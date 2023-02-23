package utils

import (
	"os/exec"
	"fmt"
	"bytes"
)

func ExecCommand(name string, args ...string) ([]byte, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("exec.LookPath(%s) err: %w", name, err)
	}
	var outb, errb bytes.Buffer
	cmd := exec.Command(path, args...)
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w:\n%s", err, string(errb.Bytes()))
	}
	return outb.Bytes(), nil
}
