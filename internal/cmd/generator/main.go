package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/generator"
	"github.com/Bitstarz-eng/event-processing-challenge/pubsub"
	"golang.org/x/net/context"
)

func main() {
	_, topics := pubsub.Setup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventCh := generator.Generate(ctx)

	for event := range eventCh {
		log.Printf("%#v\n", event)

		msgJson, err := json.Marshal(event)
		if err != nil {
			log.Fatalf("Failed to marshal message: %v", err)
		}

		topics[pubsub.TopicCasinoEventCreate].Publish(ctx, &pubsub.Message{Data: msgJson})
	}

	log.Println("finished")
}
