package main

import (
	"log"
	"time"

	"github.com/wind-c/zenoh-go/pkg/zenoh"
)

func main() {
	// Create default zenoh configuration
	config, err := zenoh.NewDefaultConfig()
	if err != nil {
		log.Fatal("Failed to create config: ", err)
	}
	defer config.Drop()

	// Use peer mode with UDP only
	config.InsertJSON5("mode", "\"peer\"")
	config.InsertJSON5("transport/unicast/enabled", "true")
	config.InsertJSON5("transport/multicast/enabled", "true")

	// Open a zenoh session
	session, err := zenoh.Open(config)
	if err != nil {
		log.Fatal("Failed to open session: ", err)
	}
	defer session.Drop()

	// Declare a publisher for the key expression
	keyExpr := "demo/example/test"
	publisher, err := zenoh.DeclarePublisherWithKeyExpr(session, keyExpr)
	if err != nil {
		log.Fatal("Failed to declare publisher: ", err)
	}
	defer publisher.Undeclare()

	log.Printf("Publisher declared on: %s", keyExpr)

	// Publish messages in a loop
	for i := 0; i < 10; i++ {
		value := []byte("Hello from zenoh-go publisher!")
		err := publisher.Put(value, zenoh.TextPlain())
		if err != nil {
			log.Printf("Failed to publish: %v", err)
		} else {
			log.Printf("Published: %s", value)
		}
		time.Sleep(1 * time.Second)
	}

	log.Println("Publisher example completed")
}
