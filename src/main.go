package main

import (
    "fmt"
	"net/http"
    "math/rand"
    "time"

	"github.com/syndtr/goleveldb/leveldb"
)

type App struct {
    db *leveldb.DB
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hey\n")
}


func main() {
    // Handle multiple transport
    http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
    rand.Seed(time.Now().Unix())

    // connect to level db 
    db, err := leveldb.OpenFile("test", nil)
    if err != nil {
        panic(fmt.Sprintf("Could not establish connection to LevelDB: %s", err))
    }
    defer db.Close()

    // Serve
    a:= App{db: db}


    fmt.Println("kv running at localhost:3000")
    http.ListenAndServe(":3000", &a)
}
