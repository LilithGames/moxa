package utils

import (
	"testing"
	assert "github.com/stretchr/testify/require"
)

func Test_ExecCommand(t *testing.T) {
	bs, err := ExecCommand("go", "version")
	assert.Nil(t, err)
	assert.True(t, len(bs) > 0)
	_, err = ExecCommand("go", "version1")
	assert.NotNil(t, err)
	println(err.Error())
	bs, err = ExecCommand("go", "help")
	assert.Nil(t, err)
	println(string(bs))
}
