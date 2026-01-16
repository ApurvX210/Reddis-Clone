package storage

import (
	// "log/slog"
	"REDDIS/parsing"
	"sync"
	"time"
)

type ExpirationInfo struct{
	timestamp	time.Time
	life		time.Duration
}


type DB struct {
	mu   			sync.RWMutex
	data 			map[string][]byte
	expirationMap	map[string] ExpirationInfo
}

func NewDb() *DB {
	return &DB{
		data: map[string][]byte{},
		expirationMap: make(map[string]ExpirationInfo),
	}
}

func (db *DB) Set(key, val []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[string(key)] = []byte(val)
	
	return nil
}

func (db *DB) Get(key []byte) ([]byte, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.data[string(key)]
	return val, ok
}

func (db *DB) Del(key []byte) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, exists := db.data[string(key)]
	if exists {
		delete(db.data, string(key))
	}
	return exists
}

func (db *DB) Hello() string {
	m := map[string]any{
		"server": "redis_clone",
		"role":   "master",
	}
	response := parsing.InitialHandShake(m)
	return response
}
