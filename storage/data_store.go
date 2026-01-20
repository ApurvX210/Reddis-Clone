package storage

import (
	// "log/slog"
	"REDDIS/parsing"
	"sync"
	"time"
)


type DB struct {
	mu   			sync.RWMutex
	data 			map[string][]byte
	expiryMu		sync.RWMutex
	expirationMap	map[string] time.Time
	
	keyList			[]string
	indexMap		map[string] int

	wg				sync.WaitGroup
}

func NewDb() *DB {
	db := &DB{
		data: map[string][]byte{},
		expirationMap: make(map[string]time.Time),
		keyList: make([]string, 0),
		indexMap: make(map[string]int),
	}

	db.cleanup()
	return db
}

func (db *DB) cleanup(){
	db.wg.Add(1)
	go db.activeCleanup()
}

func (db *DB) activeCleanup(){
	
}

func (db *DB) Set(key, val []byte, exp time.Time) error {
	db.mu.Lock()
	db.expiryMu.Lock()
	// Storing data in map
	db.data[string(key)] = []byte(val)
	// Storing expiry info in map for future cleanup
	db.expirationMap[string(key)] = exp
	db.mu.Unlock()
	db.expiryMu.Unlock()
	return nil
}

func (db *DB) Get(key []byte) ([]byte, bool) {
	db.mu.RLock()
	db.expiryMu.RLock()
	Key := string(key)
	exp,_ := db.expirationMap[Key]
	var value []byte;
	var exist bool;
	if exp.Before(time.Now()){
		value, exist = db.data[string(key)]
	}else{
		value = nil
		exist = false
	}
	db.mu.RUnlock()
	db.expiryMu.RUnlock()
	return value, exist
}

func (db *DB) Del(key []byte) bool {
	db.mu.Lock()
	db.expiryMu.Lock()

	_, exists := db.data[string(key)]
	if exists {
		delete(db.data, string(key))
		delete(db.expirationMap, string(key))
	}
	db.mu.Unlock()
	db.expiryMu.Unlock()
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
