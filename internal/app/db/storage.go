package db

import (
	"fmt"
	"sync"
)

type Storage interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	String() string
	Clear()
}

type InMemoryDB struct {
	db map[string]string
	m  sync.RWMutex
}

func NewInMemoryDBStorage(db map[string]string) Storage {
	return &InMemoryDB{db: db, m: sync.RWMutex{}}
}

func (storage *InMemoryDB) Get(key string) (string, bool) {
	storage.m.RLock()
	value, ok := storage.db[key]
	storage.m.RUnlock()
	return value, ok
}

func (storage *InMemoryDB) Set(key string, value string) {
	storage.m.Lock()
	storage.db[key] = value
	storage.m.Unlock()
}

func (storage *InMemoryDB) String() string {
	return fmt.Sprintf("%v", storage.db)
}

func (storage *InMemoryDB) Clear() {
	clear(storage.db)
}
