package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
    "flag"
    "encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

type Data struct {
    Value interface{} `json:"value"`
}


func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    key := []byte(r.URL.Path)
    key = key[1:]
    log.Println(r.Method, r.URL, r.ContentLength, string(key))

    body := Data{}
    json.NewDecoder(r.Body).Decode(&body)
    defer r.Body.Close()

    // we need to lock records if a put/delete/patch is happening, then unlock
    if r.Method == "POST" || r.Method == "DELETE" || r.Method == "PATCH" {
        if !a.lockRecord(key) {
            // Cant lock the key
            w.WriteHeader(409)
            return
        }
        defer a.unlockRecord(key)
    }

    switch r.Method {
        case "GET":
            data, code := a.Get(key)
            w.WriteHeader(code)
            w.Write(data)
        case "POST":
            _, code := a.Get(key)
            if code == 200 {
                w.WriteHeader(403)
                return
            }
            code = a.Set(key, &body)
            w.WriteHeader(code)
        case "DELETE":
            _, code := a.Get(key)
            if code == 404 {
                w.WriteHeader(409) // Conflict if it doesn't exist
                return
            }
            code = a.Delete(key)
            w.WriteHeader(code)
        case "PATCH":
            _, code := a.Get(key)
            if code != 200 {
                w.WriteHeader(code)
                return
            }
            code = a.Set(key, &body)
            w.WriteHeader(code)
    }
}

func main() {
    // Connection options
    port := flag.Int("port", 3000, "port the kv server listens on")
    leveldbPath := flag.String("leveldbPath", "lvldb", "path to leveldb")
    flag.Parse()

    if *leveldbPath == "" {
        panic("leveldbPath must have a value")
    }

    // connect to level db 
    db, err := leveldb.OpenFile(*leveldbPath, nil)
    check(err, "Error opening leveldb")
    defer db.Close()

    // Serve
    http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
    rand.Seed(time.Now().Unix())
    a:= App{db: db, lock: make(map[string]bool)}

    fmt.Println("kv listening at localhost:", *port, "\n leveldb path:", *leveldbPath)
    http.ListenAndServe(fmt.Sprintf(":%d", *port), &a)
}

func check(err error, message string) {
    if err != nil {
        panic(fmt.Sprintf("%s: %s", message, err))
    }
}
