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

func (a *App) set(w http.ResponseWriter, body *Request) {
    err := a.db.Put([]byte(body.Key), []byte(fmt.Sprint(body.Value)), nil)

    // TODO: Better error handling
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("500 - Example error message"))
    } else {
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte("201 - OK"))
    }
}

func (a *App) get(w http.ResponseWriter, body *Request) {
    data, err := a.db.Get([]byte(body.Key), nil)
    
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("500 - Example error message"))
    } else {
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(data))
    }
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
    }
}

func main() {
    // Handle multiple transport
    http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
    rand.Seed(time.Now().Unix())

    // connect to level db 
    db, err := leveldb.OpenFile("test", nil)
    check(err, "Error opening leveldb")
    defer db.Close()

    // Serve
    a:= App{db: db}


    fmt.Println("kv running at localhost:3000")
    http.ListenAndServe(":3000", &a)
}

func check(err error, message string) {
    if err != nil {
        panic(fmt.Sprintf("%s: %s", message, err))
    }
}
