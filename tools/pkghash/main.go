package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/hashicorp/logutils"

	"github.com/LilithGames/moxa/utils"
)

type Package struct {
	data map[string]any
	name string
}

func NewPackage(pkg string) (*Package, error) {
	bs, err := utils.ExecCommand("go", "list", "-json", pkg)
	if err != nil {
		return nil, fmt.Errorf("utils.ExecCommand(%s) err: %w", pkg, err)
	}
	var data map[string]any
	if err := json.Unmarshal(bs, &data); err != nil {
		return nil, fmt.Errorf("json.Unmarshal err: %w", err)
	}
	return &Package{data: data, name: pkg}, nil
}

func (it *Package) Name() string {
	return it.name
}

func (it *Package) String() string {
	return it.name
}

func (it *Package) GoMod() string {
	return it.data["Module"].(map[string]any)["GoMod"].(string)
}

func (it *Package) GoSum() string {
	return strings.ReplaceAll(it.GoMod(), "go.mod", "go.sum")
}

func (it *Package) GoVersion() string {
	return it.data["Module"].(map[string]any)["GoVersion"].(string)
}

func (it *Package) GoFiles() []string {
	files, ok := lo.FromAnySlice[string](it.data["GoFiles"].([]any))
	if !ok {
		panic(fmt.Errorf("unexpected json structure: %v", it.data["GoFiles"]))
	}
	dir := it.data["Dir"].(string)
	results := lo.Map(files, func(file string, _ int) string {
		return filepath.Join(dir, file)
	})
	return results
}

func (it *Package) Deps() []string {
	result, ok := lo.FromAnySlice[string](it.data["Deps"].([]any))
	if !ok {
		panic(fmt.Errorf("unexpected Deps structure: %v", it.data["Deps"]))
	}
	return result
}

func (it *Package) InnerDeps() []string {
	base := it.data["Module"].(map[string]any)["Path"].(string)
	return lo.Filter(it.Deps(), func(dep string, _ int) bool {
		return strings.HasPrefix(dep, base)
	})
}

func getFilesHash(files []string, salt []byte) (string, error) {
	result, err := utils.ParallelMap(files, func(path string) ([]byte, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("os.Open(%s) err: %w", path, err)
		}
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return nil, fmt.Errorf("io.Copy err: %w", err)
		}
		return hasher.Sum(nil), nil
	})
	if err != nil {
		return "", fmt.Errorf("utils.ParallelMap err: %w", err)
	}
	hasher := sha256.New()
	for _, file := range files {
		if _, err := hasher.Write(result[file]); err != nil {
			return "", fmt.Errorf("hasher.Write err: %w", err)
		}
	}
	if len(salt) > 0 {
		if _, err := hasher.Write(salt); err != nil {
			return "", fmt.Errorf("hasher.Write salt err: %w", err)
		}
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}


func getPkgFiles(pkg *Package) ([]string, error) {
	files := pkg.GoFiles()

	deps := pkg.InnerDeps()
	log.Println("[DEBUG]", fmt.Sprintf("ParallelMap deps start: %v", deps))
	result, err := utils.ParallelMap(deps, func(dep string) ([]string, error) {
		pkg1, err := NewPackage(dep)
		if err != nil {
			return nil, fmt.Errorf("NewPackage(%s) err: %w", dep, err)
		}
		return pkg1.GoFiles(), nil
	})
	log.Println("[DEBUG]", fmt.Sprintf("ParallelMap deps end"))
	if err != nil {
		return nil, fmt.Errorf("utils.ParallelMap err: %w", err)
	}
	for _, dep := range deps {
		files = append(files, result[dep]...)
	}
	return files, nil
}

func getPkgHash(pkg *Package) (string, error) {
	files, err := getPkgFiles(pkg)
	if err != nil {
		return "", fmt.Errorf("getPkgFiles(%s) err: %w", pkg.Name(), err)
	}
	ver := pkg.GoVersion()
	gomod := pkg.GoMod()
	gosum := pkg.GoSum()
	log.Println("[DEBUG]", fmt.Sprintf("hash files: %v", files))
	hash, err := getFilesHash(append(files, gomod, gosum), []byte(ver))
	log.Println("[DEBUG]", fmt.Sprintf("hash done"))
	if err != nil {
		return "", fmt.Errorf("getFilesHash(%s) err: %w", pkg.Name(), err)
	}
	return hash, nil
}

func main() {
	pkgpath := flag.String("pkg", "", "golang pkg path")
	short := flag.Bool("short", false, "hash short style")
	debug := flag.Bool("debug", false, "enable debug log")
	flag.Parse()

	logfilter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	}
	if *debug {
		logfilter.MinLevel = logutils.LogLevel("DEBUG")
	}
	log.SetOutput(logfilter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Println("[DEBUG]", fmt.Sprintf("initialize done."))

	if len(*pkgpath) == 0 {
		log.Fatalln("[FATAL]", fmt.Errorf("-pkg is required"))
	}
	pkg, err := NewPackage(*pkgpath)
	if err != nil {
		log.Fatalln("[FATAL]", fmt.Errorf("NewPackage err: %w", err))
	}
	hash, err := getPkgHash(pkg)
	if err != nil {
		log.Fatalln("[FATAL]", fmt.Errorf("getPkgHash err: %w", err))
	}
	if *short {
		fmt.Printf("%+v\n", hash[:8])
	} else {
		fmt.Printf("%+v\n", hash)
	}
}
