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
	config.InsertJSON5("transport/unicast/quic/enabled", "true")
	config.InsertJSON5("transport/unicast/quic/listen", "7448")

	session, err := zenoh.Open(config)
	if err != nil {
		log.Fatal("Failed to open session: ", err)
	}
	defer session.Drop()

	keyExpr := "demo/quic/test"
	subscriber, err := zenoh.DeclareSubscriber(session, keyExpr, func(sample zenoh.Sample) {
		log.Printf("Received - Key: %s, Value: %s", sample.KeyExpr, string(sample.Payload))
	})
	if err != nil {
		log.Fatal("Failed to declare subscriber: ", err)
	}
	defer subscriber.Undeclare()

	log.Printf("QUIC Subscriber declared on: %s", keyExpr)
	log.Println("Waiting for messages... Press Ctrl+C to exit")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	log.Println("QUIC Subscriber example completed")
}
