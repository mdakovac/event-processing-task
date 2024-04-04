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

	client, topics := pubsub.Setup()

	subscriptionId := "data_processing_service.subscription"
	subscription := pubsub.GetSubscription(client, subscriptionId, topics[pubsub.TopicCasinoEventCreate])

	go func() {
		err := subscription.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			//fmt.Printf("Received message: %s\n", string(msg.Data))

			err := handleMessageReceived(msg.Data, currencyService, playerService, descriptionService, aggregationService)
			if err != nil {
				msg.Nack()
			} else {
				msg.Ack()
			}
		})
		if err != nil {
			log.Printf("Error receiving message: %v", err)
		}
	}()
	log.Println("Pub/Sub listener started")

	go func() {
		router := gin.Default()
		aggregation_controller.SetupAggregationController(router, aggregationService)
		router.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func handleMessageReceived(
	data []byte,
	currencyService currency_service.CurrencyServiceType,
	playerService player_service.PlayerServiceType,
	descriptionService description_service.DescriptionServiceType,
	aggregationService aggregation_service.AggregationServiceType,
) error {
	var event casino.Event
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("Failed to unmarshal message data: %v", err)
		return err
	}

	currencyService.ConvertCurrency(&event)
	playerService.AssignPlayerData(&event)
	descriptionService.AssignDescription(&event)

	aggregationService.AddEventToAggregation(&event)

	forPrint, _ := json.MarshalIndent(event, "", "    ")
	log.Println("Enriched event", string(forPrint))
	return nil
}
