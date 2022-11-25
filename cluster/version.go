package cluster

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

const DataVersionKey = "DataVersion"

type IVersionManager interface {
	Migrating() bool
	CodeVersion() string
	DataVersion() string
	Commit()
}

type VersionManager struct {
	code  string
	store IStorage
}

func NewVersionManager(store IStorage) (IVersionManager, error) {
	code, err := getExecutableSHA256()
	if err != nil {
		return nil, fmt.Errorf("getExecutableSHA256 err: %w", err)
	}
	return &VersionManager{store: store, code: code}, nil
}

func (it *VersionManager) Migrating() bool {
	if it.DataVersion() == "" {
		return false
	}
	return it.CodeVersion() != it.DataVersion()
}
func (it *VersionManager) CodeVersion() string {
	return it.code
}
func (it *VersionManager) DataVersion() string {
	return it.store.Get(DataVersionKey)
}
func (it *VersionManager) Commit() {
	it.store.Set(DataVersionKey, it.code)
}

func getExecutableSHA256() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("os.Executable err: %w", err)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("os.Open(%s) err: %w", path, err)
	}
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("io.Copy err: %w", err)
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
