package main

import (
	"flag"
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

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName)
	store, err := storage.New(dsn)
	if err != nil {
		log.Fatal("create storage layer err: ", err)
	}

	defer store.Pool.Close()

	v := flag.Int("db-version", 0, "sets DB migration specific version")
	flag.Parse()

	dbVer := *v
	switch dbVer {
	case 0:
		err = store.UpDBVersion(dsn)
		if err != nil {
			log.Fatal("db migrate error: ", err)
		}
	default:
		err = store.SetDBVersion(dsn, uint(dbVer))
		if err != nil {
			log.Fatal("db migrate to specific version error: ", err)
		}
	}
}
