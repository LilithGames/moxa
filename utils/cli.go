package utils

import (
	"io"
	"os/exec"
	"fmt"
	"bytes"
)

func ExecCommand(name string, args ...string) ([]byte, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("exec.LookPath(%s) err: %w", name, err)
	}
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(path, args...)
	outpipe, _ := cmd.StdoutPipe()
	errpipe, _ := cmd.StderrPipe()
	var err1, err2 error
	go func() {
		_, err1 = io.Copy(&outbuf, outpipe)
	}()
	go func() {
		_, err2 = io.Copy(&errbuf, errpipe)
	}()
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w:\n%s", err, string(errbuf.Bytes()))
	}
	if err1 != nil {
		return nil, fmt.Errorf("io.Copy stdout err: %w", err1)
	}
	if err2 != nil {
		return nil, fmt.Errorf("io.Copy stderr err: %w", err2)
	}
	return outbuf.Bytes(), nil
}
