package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/cache_service"
	"github.com/Bitstarz-eng/event-processing-challenge/data_enrichment/currency/currency_repository"
	"github.com/Bitstarz-eng/event-processing-challenge/data_enrichment/currency/currency_service"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/pubsub"
)

func main() {
	var currencyRepository = currency_repository.NewCurrencyRepository(cache_service.CreateCache(1*time.Minute, 1*time.Minute))
	var currencyService = currency_service.NewCurrencyService(currencyRepository)

	ctx := context.Background()

	client, topics := pubsub.Setup()

	subscriptionId := "data_enrichment_service.subscription"
	subscription := pubsub.GetSubscription(client, subscriptionId, topics["CasinoEvent.create"])

	// Start receiving messages
	go func() {
		err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			fmt.Printf("Received message: %s\n", string(msg.Data))

			var event casino.Event
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("Failed to unmarshal message data: %v", err)
				msg.Nack()
				return
			}

			var eventPointer = &event

			_, err := currencyService.ConvertCurrency(eventPointer)
			if err != nil {
				log.Println(err)
			}

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
