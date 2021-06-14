package main

import (
	"fmt"
	"log"

	"sample-project/config"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal("config parse err: ", err)
	}

	fmt.Printf("%+v\n", conf)
}
