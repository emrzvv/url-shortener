package db

import "fmt"

type Storage interface {
	Init()
	Get(key string) (string, bool)
	Contains(key string) bool
	Set(key string, value string)
	String() string
	Clear()
}

type InMemoryDB struct {
	db map[string]string
}

func (storage *InMemoryDB) Init() {
	storage.db = make(map[string]string)
}

func (storage *InMemoryDB) Get(key string) (string, bool) {
	value, ok := storage.db[key]
	return value, ok
}

func (storage *InMemoryDB) Contains(key string) bool {
	_, ok := storage.db[key]
	return ok
}

func (storage *InMemoryDB) Set(key string, value string) {
	storage.db[key] = value
}

func (storage *InMemoryDB) String() string {
	return fmt.Sprintf("%v", storage.db)
}

func (storage *InMemoryDB) Clear() {
	clear(storage.db)
}
