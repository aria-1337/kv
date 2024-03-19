package main

import (
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

type App struct {
    db *leveldb.DB 
    mlock sync.Mutex
    lock map[string]bool
}

func (a *App) lockRecord(key []byte) bool {
    a.mlock.Lock()
    defer a.mlock.Unlock()
    if _, exists := a.lock[string(key)]; exists {
        return false
    }
    a.lock[string(key)] = true
    return true
}

func (a *App) unlockRecord(key []byte) {
    a.mlock.Lock()
    delete(a.lock, string(key))
    a.mlock.Unlock()
}

func (a *App) Get(key []byte) ([]byte, int) {
    data, err := a.db.Get(key, nil)
    if err == leveldb.ErrNotFound {
        return make([]byte, 0), 404
    }
    return data, 200
}

func (a *App) Set(key []byte, body *Data) int {
    err := a.db.Put(key, []byte(fmt.Sprint(body.Value)), nil)
    if err != nil {
        return 500
    }
    return 200
}

func (a *App) Delete(key []byte) int {
    err := a.db.Delete(key, nil)
    if err != nil {
        return 500
    }
    return 200
}
