package pubsub

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

type PubsubClientType = pubsub.Client
type PubsubTopicType = pubsub.Topic
type PubsubMessageType = pubsub.Message

func Setup(projectId string) (*PubsubClientType, map[string]*PubsubTopicType) {
	client := CreateClient(projectId)
	topics := CreateTopics(client)

	return client, topics
}

func CreateClient(projectId string) *PubsubClientType {
	ctx := context.Background()

	// Create a Pub/Sub client
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func CreateTopics(client *PubsubClientType) map[string]*PubsubTopicType {
	topicNames := []string{"CasinoEvent.create"}

	topics := make(map[string]*PubsubTopicType)
	for _, v := range topicNames {
		topic := createTopicIfNotExists(client, v)
		topics[v] = topic
	}

	return topics
}

func createTopicIfNotExists(client *PubsubClientType, t string) *PubsubTopicType{
	ctx := context.Background()

	topic := client.Topic(t)
	ok, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("Failed to check if topic exists: %v", err)
	}

	if ok{
		return topic
	} else {
		topic, err := client.CreateTopic(ctx, t)
		if err != nil {
			log.Fatalf("Failed to create topic: %v", err)
		}
		log.Printf("Topic %s created.\n", t)
		return topic
	}
}