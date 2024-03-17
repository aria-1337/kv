package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type Request struct {
    Key string `json:"key"`
    Value interface{} `json:"value"`
}

type App struct {
    db *leveldb.DB
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

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Println(r.Method, r.URL, r.ContentLength)

    // Read request
    defer r.Body.Close()
    body := Request{}
    json.NewDecoder(r.Body).Decode(&body)

    // Route request
    switch r.Method {
        case "POST":
            a.set(w, &body)
        case "GET":
            a.get(w, &body)
        case "DELETE":
            a.delete(w, &body)
    }
}

func main() {
    // connect to level db 
    db, err := leveldb.OpenFile("test", nil)
    check(err, "Error opening leveldb")
    defer db.Close()


    // Serve
    http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
    rand.Seed(time.Now().Unix())
    a:= App{db: db}

    fmt.Println("kv running at localhost:3000")
    http.ListenAndServe(":3000", &a)
}

func check(err error, message string) {
    if err != nil {
        panic(fmt.Sprintf("%s: %s", message, err))
    }
}
