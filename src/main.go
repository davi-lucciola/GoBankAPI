package main

import (
	"log"
)

func main() {
	storage, err := NewPostgresStorage()

	if err != nil {
		log.Fatal(err)
	}

	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":8084", storage)
	server.Run()
}
