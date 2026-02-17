package main

import (
	"log"

	"github.com/wind-c/zenoh-go/pkg/zenoh"
)

func main() {
	config, err := zenoh.NewDefaultConfig()
	if err != nil {
		log.Fatal("Failed to create config: ", err)
	}
	defer config.Drop()

	config.InsertJSON5("mode", "\"peer\"")
	config.InsertJSON5("transport/shared_memory/enabled", "true")

	session, err := zenoh.Open(config)
	if err != nil {
		log.Fatal("Failed to open session: ", err)
	}
	defer session.Drop()

	shmProvider, err := session.SharedMemoryProvider()
	if err != nil {
		log.Printf("Shared memory not available: %v", err)
		log.Println("Note: Shared memory requires zenoh-c built with Z_FEATURE_SHM enabled")
		return
	}

	if !shmProvider.IsValid() {
		log.Println("Shared memory provider is not valid")
		return
	}

	log.Printf("Shared memory provider obtained: %v", shmProvider.IsValid())

	keyExpr := "demo/shm-pub/test"
	publisher, err := zenoh.DeclarePublisherWithKeyExpr(session, keyExpr)
	if err != nil {
		log.Fatal("Failed to declare publisher: ", err)
	}
	defer publisher.Undeclare()

	log.Printf("Publisher declared on: %s", keyExpr)

	value := []byte("Hello via Shared Memory!")
	err = publisher.Put(value, zenoh.TextPlain())
	if err != nil {
		log.Printf("Failed to publish: %v", err)
	} else {
		log.Printf("Published: %s", value)
	}

	log.Println("Shared Memory example completed")
}
