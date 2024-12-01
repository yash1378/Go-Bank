package main

import (
	"fmt"
	"log"
)

func main() {

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", store)
	// Call Init to create the account table
	if err := store.Init(); err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	fmt.Println("yeah buddy")
	server := NewAPIServer(":3000", store)
	server.Run()

}
