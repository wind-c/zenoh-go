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

	session, err := zenoh.Open(config)
	if err != nil {
		log.Fatal("Failed to open session: ", err)
	}
	defer session.Drop()

	selector := "demo/**"
	iterator, err := zenoh.GetWithIterator(session, selector)
	if err != nil {
		log.Fatal("Failed to get: ", err)
	}

	log.Printf("Querying: %s", selector)

	count := 0
	for iterator.Next() {
		reply := iterator.Reply()
		if reply.IsOk() {
			log.Printf("Reply[%d] - Key: %s, Value: %s", count, reply.KeyExpr(), string(reply.Value()))
			count++
		} else {
			log.Printf("Reply[%d] - Error: %s", count, reply.Error())
		}
	}

	log.Printf("Query example completed. Got %d replies", count)
}
