package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/cache_service"
	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/aggregation/aggregation_controller"
	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/aggregation/aggregation_service"
	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/currency/currency_repository"
	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/currency/currency_service"
	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/description/description_service"
	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/player/player_repository"
	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/player/player_service"
	"github.com/Bitstarz-eng/event-processing-challenge/db"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/pubsub"
	"github.com/Bitstarz-eng/event-processing-challenge/util/env_vars"
	"github.com/gin-gonic/gin"
)

func main() {
	env_vars.SetEnvVars()

	db := db.Connect()

	var currencyRepository = currency_repository.NewCurrencyRepository(cache_service.CreateCache(1*time.Minute, 1*time.Minute))
	var currencyService = currency_service.NewCurrencyService(currencyRepository)

	var playerRepository = player_repository.NewPlayerRepository(db)
	var playerService = player_service.NewPlayerService(playerRepository)

	var descriptionService = description_service.NewDescriptionService(currencyService)

	var aggregationService = aggregation_service.NewAggregationService()

	ctx := context.Background()

	client, topics := pubsub.Setup()

	subscriptionId := "data_processing_service.subscription"
	subscription := pubsub.GetSubscription(client, subscriptionId, topics["CasinoEvent.create"])

	go func() {
		err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			//fmt.Printf("Received message: %s\n", string(msg.Data))

			var event casino.Event
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("Failed to unmarshal message data: %v", err)
				msg.Nack()
				return
			}

			currencyService.ConvertCurrency(&event)
			playerService.AssignPlayerData(&event)
			descriptionService.AssignDescription(&event)

			aggregationService.AddEventToAggregation(&event)

			forPrint1, _ := json.MarshalIndent(aggregationService.GetAggregation(), "", "    ")
			log.Println("Aggregation event", string(forPrint1))

			//forPrint, _ := json.MarshalIndent(event, "", "    ")
			//log.Println("Enriched event", string(forPrint))

			msg.Ack()
		})
		if err != nil {
			log.Fatalf("Error receiving message: %v", err)
		}
	}()

	log.Println("Pub/Sub listener started")


	// Setup Gin router
	router := gin.Default()
	aggregation_controller.SetupAggregationController(router, aggregationService)
	router.Run("localhost:8080")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down Pub/Sub listener")
}
