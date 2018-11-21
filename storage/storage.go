package storage

import (
    "log"

    "github.com/boltdb/bolt"
)

const (
    dbLocation       = "./bolt.db"
    TraveltextBucket = "TraveltextUsers"
)

var (
    DB *bolt.DB
)

func init() {
    var err error
    DB, err = bolt.Open(dbLocation, 0644, nil)
    if err != nil {
        log.Fatal(err)
    }
}
