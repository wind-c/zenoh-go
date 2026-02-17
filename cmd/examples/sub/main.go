package main

import (
	"log"
	"os"
	"os/signal"

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

	// Declare a subscriber for the key expression
	// Using wildcard pattern to match multiple keys
	keyExpr := "demo/example/test"
	subscriber, err := zenoh.DeclareSubscriber(session, keyExpr, func(sample zenoh.Sample) {
		log.Printf("Received sample - Key: %s, Value: %s", sample.KeyExpr, string(sample.Payload))
	})
	if err != nil {
		log.Fatal("Failed to declare subscriber: ", err)
	}
	defer subscriber.Undeclare()

	log.Printf("Subscriber declared on: %s", keyExpr)
	log.Println("Waiting for messages... Press Ctrl+C to exit")

	// Wait for interrupt signal to gracefully shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	log.Println("Subscriber example completed")
}
