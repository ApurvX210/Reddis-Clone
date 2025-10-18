package main

import (
	// "log/slog"
	"sync"
)

type DB struct {
	mu sync.RWMutex
	data map[string][]byte
}

func NewDb() *DB {
	return &DB{
		data: map[string][]byte{},
	}
}

func (db *DB) Set(key,val []byte) error{
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[string(key)] = []byte(val)
	// slog.Info("Set Instruction commited successfully","key",key,"value",val)
	return nil
}

func (db *DB) Get(key []byte) ([]byte,bool){
	db.mu.RLock()
	defer db.mu.RUnlock()
	val,ok := db.data[string(key)]
	return val,ok
}

func (db *DB) hello() string{
	m:= map[string]any{
		"server":"redis_clone",
		"role": "master",
	}
	response := initialHandShake(m)
	return response
}
