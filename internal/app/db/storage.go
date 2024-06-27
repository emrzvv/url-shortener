package db

import "fmt"

type Storage struct {
	db map[string]string
}

func (storage *Storage) Init() {
	storage.db = make(map[string]string)
}

func NewStorage() Storage {
	return Storage{}
}

func (storage *Storage) Get(key string) (string, bool) {
	value, ok := storage.db[key]
	return value, ok
}

func (storage *Storage) Contains(key string) bool {
	_, ok := storage.db[key]
	return ok
}

func (storage *Storage) Set(key string, value string) {
	storage.db[key] = value
}

func (storage *Storage) String() string {
	return fmt.Sprintf("%v", storage.db)
}
