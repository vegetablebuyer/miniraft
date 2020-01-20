package main

import (
	"sync"
)

// The key-value database
type DB struct {
	data map[string]string
	mutex sync.RWMutex
}

// Creates a new database
func NewDB() *DB {
	return &DB{
		data: make(map[string]string),
	}
}

// Retrieves the value for the given key
func (db *DB) Get(key string) string {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return db.data[key]
}

// Sets the value for the given key
func (db *DB) Set(key string, value string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.data[key] = value
}
