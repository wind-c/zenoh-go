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
	config.InsertJSON5("transport/unicast/enabled", "true")
	config.InsertJSON5("transport/multicast/enabled", "true")

	session, err := zenoh.Open(config)
	if err != nil {
		log.Fatal("Failed to open session: ", err)
	}
	defer session.Drop()

	keyExpr := "demo/**"

	queryable, err := zenoh.DeclareQueryable(session, keyExpr, func(query zenoh.Query) {
		log.Printf("Received query - Key: %s, Parameters: %s", query.KeyExpr(), query.Parameters())

		value := "Hello from Queryable! Key: " + query.KeyExpr()
		err := query.Reply(query.KeyExpr(), []byte(value), zenoh.EncodingTextPlain)
		if err != nil {
			log.Printf("Failed to reply: %v", err)
		} else {
			log.Printf("Replied with: %s", value)
		}
	})
	if err != nil {
		log.Fatal("Failed to declare queryable: ", err)
	}
	defer queryable.Undeclare()

	log.Printf("Queryable declared on: %s", keyExpr)
	log.Println("Waiting for queries... Press Ctrl+C to exit")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	log.Println("Queryable example completed")
}
