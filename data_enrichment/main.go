package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bitstarz-eng/event-processing-challenge/pubsub"
)

func main() {
	ctx := context.Background()

	client, topics := pubsub.Setup()

	subscriptionId := "data_enrichment_service.subscription"
	subscription := pubsub.GetSubscription(client, subscriptionId, topics["CasinoEvent.create"])

	// Start receiving messages
	go func() {
		err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			fmt.Printf("Received message: %s\n", string(msg.Data))
			msg.Ack()
		})
		if err != nil {
			log.Fatalf("Error receiving message: %v", err)
		}
	}()
	fmt.Println("Pub/Sub listener started")

	// Create a channel to handle shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down Pub/Sub listener")
}
