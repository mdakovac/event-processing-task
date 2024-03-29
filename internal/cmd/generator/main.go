package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/generator"
	"github.com/Bitstarz-eng/event-processing-challenge/pubsub"
	"golang.org/x/net/context"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, topics := pubsub.Setup()

	eventCh := generator.Generate(ctx)

	for event := range eventCh {
		log.Printf("%#v\n", event)

		var inputBuffer bytes.Buffer
		gob.NewEncoder(&inputBuffer).Encode(event)
		topics["CasinoEvent.create"].Publish(ctx, &pubsub.Message{Data: inputBuffer.Bytes()})

	}

	topics["CasinoEvent.create"].Stop()
	log.Println("finished")

}
