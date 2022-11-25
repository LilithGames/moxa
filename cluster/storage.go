package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type IStorage interface {
	Get(key string) string
	Set(key string, value string)
}

type LocalStorage struct {
	mu     sync.RWMutex
	state  map[string]string
	dbpath string
}

func NewLocalStorage(dir string) IStorage {
	os.MkdirAll(dir, os.ModePerm)
	it := &LocalStorage{
		dbpath: fmt.Sprintf("%s/storage.db", dir),
	}
	it.state = it.load()
	it.save(it.state)
	return it
}

func (it *LocalStorage) load() map[string]string {
	var result map[string]string
	bs, err := ioutil.ReadFile(it.dbpath)
	if err != nil {
		log.Println("[WARN]", fmt.Errorf("LocalStorage.load ioutil.ReadFile(%s) err: %w", it.dbpath, err))
		return result
	}
	if err := json.Unmarshal(bs, &result); err != nil {
		log.Println("[WARN]", fmt.Errorf("LocalStorage.load json.Unmarshal err: %w", err))
	}
	return result
}

func (it *LocalStorage) save(s map[string]string) {
	bs, err := json.Marshal(s)
	if err != nil {
		log.Println("[WARN]", fmt.Errorf("LocalStorage.save json.Marshal err: %w", err))
		return
	}
	if err := ioutil.WriteFile(it.dbpath, bs, os.ModePerm); err != nil {
		log.Println("[WARN]", fmt.Errorf("LocalStorage.save ioutil.WriteFile err: %w", err))
	}
	return
}

func (it *LocalStorage) Get(key string) string {
	it.mu.RLock()
	defer it.mu.RUnlock()
	if v, ok := it.state[key]; !ok {
		return ""
	} else {
		return v
	}
}

func (it *LocalStorage) Set(key string, value string) {
	it.mu.Lock()
	defer it.mu.Unlock()
	if it.state == nil {
		it.state = make(map[string]string)
	}
	it.state[key] = value
	it.save(it.state)
}
