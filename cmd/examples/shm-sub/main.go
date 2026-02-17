package main

import (
	"log"
	"os"
	"os/signal"

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
	subscriber, err := zenoh.DeclareSubscriber(session, keyExpr, func(sample zenoh.Sample) {
		log.Printf("Received - Key: %s, Value: %s", sample.KeyExpr, string(sample.Payload))
	})
	if err != nil {
		log.Fatal("Failed to declare subscriber: ", err)
	}
	defer subscriber.Undeclare()

	log.Printf("Shared Memory Subscriber declared on: %s", keyExpr)
	log.Println("Waiting for messages... Press Ctrl+C to exit")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	log.Println("Shared Memory Subscriber example completed")
}
