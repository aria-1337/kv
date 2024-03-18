package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
    "sync"
    "flag"

	"github.com/syndtr/goleveldb/leveldb"
)

type Request struct {
    Key string `json:"key"`
    Value interface{} `json:"value"`
}

type App struct {
    db *leveldb.DB
    mlock sync.Mutex
    lock map[string]bool
}

func (a *App) lockRecord(key string) bool {
    a.mlock.Lock()
    defer a.mlock.Unlock()
    if _, exists := a.lock[key]; exists {
        return false
    }
    a.lock[key] = true
    return true
}

func (a *App) unlockRecord(key string) {
    a.mlock.Lock()
    delete(a.lock, key)
    a.mlock.Unlock()
}

func (a *App) getRecord(key string) ([]byte, int) {
    data, err := a.db.Get([]byte(key), nil)
    code := 200
    if err != leveldb.ErrNotFound  && err != nil {
        code = 500
    }
    if err == leveldb.ErrNotFound {
        code = 404
    }
    return data, code
}

// POST - 201=success 403=forbidden overwrite 
func (a *App) set(w http.ResponseWriter, body *Request) {
    _, code := a.getRecord(body.Key)
    if code == 200 {
        w.WriteHeader(403)
        return
    }
    a.db.Put([]byte(body.Key), []byte(fmt.Sprint(body.Value)), nil)
    w.WriteHeader(201)
}

// GET - 200=success 404=record doesnt exist
func (a *App) get(w http.ResponseWriter, body *Request) {
    data, code := a.getRecord(body.Key)
    w.WriteHeader(code)
    w.Write([]byte(data))
}

// DELETE - 204=success 409=not deleted
func (a *App) delete(w http.ResponseWriter, body *Request) {
    _, code := a.getRecord(body.Key)
    if code == 404 {
        w.WriteHeader(409)
        return
    }
    a.db.Delete([]byte(body.Key), nil)
    w.WriteHeader(204)
}

// PATCH - 200=success 404=record doesnt exist
func (a *App) update(w http.ResponseWriter, body *Request) {
    _, code := a.getRecord(body.Key)
    if code == 404 {
        w.WriteHeader(404)
        return
    }
    a.db.Put([]byte(body.Key), []byte(fmt.Sprint(body.Value)), nil)
    w.WriteHeader(200)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Println(r.Method, r.URL, r.ContentLength)

    // Read request
    defer r.Body.Close()
    body := Request{}
    json.NewDecoder(r.Body).Decode(&body)

    // we need to lock records if a put/delete/patch is happening, then unlock
    if r.Method == "POST" || r.Method == "DELETE" || r.Method == "PATCH" {
        if !a.lockRecord(body.Key) {
            // Cant lock the key
            w.WriteHeader(409)
            return
        }
        defer a.unlockRecord(body.Key)
    }


    // Route request
    switch r.Method {
        case "POST":
            a.set(w, &body)
        case "GET":
            a.get(w, &body)
        case "DELETE":
            a.delete(w, &body)
        case "PATCH":
            a.update(w, &body)
    }
}

func main() {
    // Connection options
    port := flag.Int("port", 3000, "port the kv server listens on")

    // connect to level db 
    db, err := leveldb.OpenFile("test", nil)
    check(err, "Error opening leveldb")
    defer db.Close()

    // Serve
    http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
    rand.Seed(time.Now().Unix())
    a:= App{db: db, lock: make(map[string]bool)}

    fmt.Println("kv listening at localhost:", *port)
    http.ListenAndServe(fmt.Sprintf(":%d", *port), &a)
}

func check(err error, message string) {
    if err != nil {
        panic(fmt.Sprintf("%s: %s", message, err))
    }
}
