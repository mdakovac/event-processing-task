package pubsub

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/Bitstarz-eng/event-processing-challenge/util/env_vars"
)

type Client = pubsub.Client
type Topic = pubsub.Topic
type Message = pubsub.Message

type Subscription = pubsub.Subscription
type SubscriptionConfig = pubsub.SubscriptionConfig

const (
	TopicCasinoEventCreate = "CasinoEvent.Create"
)

var topics = []string{
	TopicCasinoEventCreate,
}

func Setup() (*Client, map[string]*Topic) {
	projectId := env_vars.EnvVariables.PUBSUB_PROJECT_ID

	c := CreateClient(projectId)
	t := CreateTopics(c)

	return c, t
}

func CreateClient(projectId string) *Client {
	ctx := context.Background()

	// Create a Pub/Sub client
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func CreateTopics(client *Client) map[string]*Topic {
	t := make(map[string]*Topic)
	for _, v := range topics {
		topic := createTopicIfNotExists(client, v)
		t[v] = topic
	}

	return t
}

func GetSubscription(client *Client, subscriptionId string, topic *Topic) *Subscription {
	ctx := context.Background()

	subscription := client.Subscription("data_processing_service.subscription")
	exists, err := subscription.Exists(ctx)
	if err != nil {
		log.Fatalf("Failed to check if subscription exists: %v", err)
	}

	if exists {
		return subscription
	}

	log.Println("Creating subscription")
	subscription, err = client.CreateSubscription(ctx, subscriptionId, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		log.Fatalf("Failed to create subscription: %v", err)
	}

	return subscription
}

func createTopicIfNotExists(client *Client, t string) *Topic {
	ctx := context.Background()

	topic := client.Topic(t)
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("Failed to check if topic exists: %v", err)
	}

	if exists {
		return topic
	}

	topic, err = client.CreateTopic(ctx, t)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}
	log.Printf("Topic %s created.\n", t)
	return topic
}
