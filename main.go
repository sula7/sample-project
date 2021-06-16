package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4/source/httpfs"

	"sample-project/config"
	"sample-project/storage"
)

//go:embed migrations
var migrations embed.FS

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

	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		log.Fatal("migrate source create error: ", err)
	}

	dbVer := *v
	switch dbVer {
	case 0:
		err = store.UpDBVersion(dsn, source)
		if err != nil {
			log.Fatal("db migrate error: ", err)
		}
	default:
		err = store.SetDBVersion(dsn, uint(dbVer), source)
		if err != nil {
			log.Fatal("db migrate to specific version error: ", err)
		}
		fmt.Println("WARNING! DB migrated into version ", dbVer)
	}
}
