package cluster

import (
	"fmt"
	"time"
	"strings"
	"log"
	"errors"
	"strconv"
	"os"

	"github.com/lni/goutils/vfs"
	"github.com/samber/lo"
)

var ErrSnapshotNotFound error = errors.New("ErrSnapshotNotFound")

type Snapshot struct {
	ShardID    uint64
	Version    string
	NodeHostID string
	Index      uint64
	BaseDir    string
}

func (it *Snapshot) Dir() string {
	return fmt.Sprintf("%s/%d/%s", it.BaseDir, it.ShardID, it.Version)
}

func (it *Snapshot) Path() string {
	return fmt.Sprintf("%s/%s/snapshot-%016X", it.Dir(), it.NodeHostID, it.Index)
}

func (it *Snapshot) CommitPath() string {
	return fmt.Sprintf("%s.commit", it.Path())
}

type ISnapshotManager interface {
	Prepare(shardID uint64, version string) string
	New(shardID uint64, version string, index uint64) *Snapshot
	Get(shardID uint64, version string, nhid string, index uint64) (*Snapshot, error)
	GetLatest(shardID uint64, version string, expire time.Duration) (*Snapshot, error)
	Commit(s *Snapshot) error
	Committed(s *Snapshot) bool
}

type SnapshotManager struct {
	baseDir string
	nhid string
	fs vfs.FS
}

func NewSnapshotManager(baseDir string, nhid string, fs vfs.FS) ISnapshotManager {
	return &SnapshotManager{baseDir: baseDir, nhid: nhid, fs: fs}
}

func (it *SnapshotManager) Prepare(shardID uint64, version string) string {
	path := fmt.Sprintf("%s/%d/%s/%s", it.baseDir, shardID, version, it.nhid)
	if err := it.fs.MkdirAll(path, os.ModePerm); err != nil {
		log.Println("[WARN]", fmt.Errorf("fs.MkdirAll(%s) err: %w", path, err))
	}
	return path
}

func (it *SnapshotManager) New(shardID uint64, version string, index uint64) *Snapshot {
	return &Snapshot{ShardID: shardID, Version: version, NodeHostID: it.nhid, Index: index, BaseDir: it.baseDir}
}

func (it *SnapshotManager) Get(shardID uint64, version string, nhid string, index uint64) (*Snapshot, error) {
	snapshot := &Snapshot{ShardID: shardID, Version: version, NodeHostID: nhid, Index: index, BaseDir: it.baseDir}
	if !it.Committed(snapshot) {
		return nil, ErrSnapshotNotFound
	}
	return snapshot, nil
}

func (it *SnapshotManager) GetLatest(shardID uint64, version string, expire time.Duration) (*Snapshot, error) {
	dir := fmt.Sprintf("%s/%d/%s", it.baseDir, shardID, version)
	now := time.Now()
	snapshots := make([]*Snapshot, 0)
	nhids, err := it.fs.List(dir)
	if err != nil {
		return nil, fmt.Errorf("fs.List(%s) err: %w", dir, err)
	}
	for _, nhid := range nhids {
		ndir := fmt.Sprintf("%s/%s", dir, nhid)
		indexes, err := it.fs.List(ndir)
		if err != nil {
			return nil, fmt.Errorf("fs.List(%s) err: %w", ndir, err)
		}
		for _, index := range indexes {
			if strings.HasSuffix(index, ".commit") {
				continue
			}
			i, err := getIndex(index)
			if err != nil {
				log.Println("[WARN]", fmt.Errorf("skip invalid snapshot index %s err: %w", index, err))
				continue
			}
			idir := fmt.Sprintf("%s/%s", ndir, index)
			fi, err := it.fs.Stat(idir)
			if err != nil {
				log.Println("[WARN]", fmt.Errorf("skip stat snapshot index %s err: %w", idir, err))
				continue
			}
			if !fi.IsDir() {
				log.Println("[WARN]", fmt.Errorf("skip snapshot index %s is not a directory", idir))
				continue
			}
			if now.Sub(fi.ModTime()) > expire {
				continue
			}
			snapshot := &Snapshot{ShardID: shardID, Version: version, NodeHostID: nhid, Index: i, BaseDir: it.baseDir}
			if !it.Committed(snapshot) {
				continue
			}
			snapshots = append(snapshots, snapshot)
		}
	}
	snapshot := lo.MaxBy(snapshots, func(a *Snapshot, b *Snapshot) bool {
		return a.Index > b.Index
	})
	if snapshot == nil {
		return nil, ErrSnapshotNotFound
	}
	return snapshot, nil
}

func (it *SnapshotManager) Commit(s *Snapshot) error {
	path := s.CommitPath()
	f, err := it.fs.Create(path)
	if err != nil {
		return fmt.Errorf("fs.Create(%s) err: %w", path, err)
	}
	f.Close()
	return nil
}

func (it *SnapshotManager) Committed(s *Snapshot) bool {
	path := s.CommitPath()
	_, err := it.fs.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func getIndex(filename string) (uint64, error) {
	items := strings.Split(filename, "-")
	if len(items) != 2 {
		return 0, fmt.Errorf("Invalid snapshot filename: items(%d) != 2", len(items))
	}
	if items[0] != "snapshot" {
		return 0, fmt.Errorf("Invalid snapshot filename: first item(%s) != snapshot", items[0])
	}
	index, err := strconv.ParseUint(items[1], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid snaphost filename: ParseUint(%s) err: %w", items[1], err)
	}
	return index, nil
}
