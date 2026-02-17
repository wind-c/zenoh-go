package main

import (
	"log"
	"time"

	"github.com/wind-c/zenoh-go/pkg/zenoh"
)

func main() {
	config, err := zenoh.NewDefaultConfig()
	if err != nil {
		log.Fatal("Failed to create config: ", err)
	}
	defer config.Drop()

	config.InsertJSON5("mode", "\"peer\"")
	config.InsertJSON5("transport/unicast/quic/enabled", "true")
	config.InsertJSON5("transport/unicast/quic/listen", "7447")

	session, err := zenoh.Open(config)
	if err != nil {
		log.Fatal("Failed to open session: ", err)
	}
	defer session.Drop()

	keyExpr := "demo/quic/test"
	publisher, err := zenoh.DeclarePublisherWithKeyExpr(session, keyExpr)
	if err != nil {
		log.Fatal("Failed to declare publisher: ", err)
	}
	defer publisher.Undeclare()

	log.Printf("QUIC Publisher declared on: %s", keyExpr)
	log.Println("Publishing via QUIC transport...")

	for i := 0; i < 10; i++ {
		value := []byte("Hello via QUIC!")
		err := publisher.Put(value, zenoh.TextPlain())
		if err != nil {
			log.Printf("Failed to publish: %v", err)
		} else {
			log.Printf("Published: %s", value)
		}
		time.Sleep(1 * time.Second)
	}

	log.Println("QUIC Publisher example completed")
}
