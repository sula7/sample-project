package main

import (
	"fmt"
	"log"

	"sample-project/config"
	"sample-project/storage"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal("config parse err: ", err)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName)
	store, err := storage.New(connStr)
	if err != nil {
		log.Fatal("create storage layer err: ", err)
	}

	defer store.Pool.Close()
}
