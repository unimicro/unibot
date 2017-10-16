package storage

import (
	"log"

	"github.com/boltdb/bolt"
)

const (
	dbLocation    = "/tmp/bolt.db"
	JenkinsBucket = "jenkinsJobs"
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
