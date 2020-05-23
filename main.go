package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/KouT127/pub-sub-practice/model"
	"google.golang.org/api/option"
	"log"
	"os"
)

type PubSubClient struct {
	*pubsub.Client
}
type Topic struct {
	*pubsub.Topic
}

func main() {
	ctx := context.Background()
	client := NewPubSubClient()
	topic := client.configureTopics()
	ID, err := topic.PublishMessage(ctx, "hello")
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Print(ID)
}

const resizeTopic = "resizeTopic"

func NewPubSubClient() *PubSubClient {
	ctx := context.Background()
	opt := option.WithCredentialsJSON([]byte(mustGetEnv("PUBSUB_SERVICE")))
	client, err := pubsub.NewClient(ctx, mustGetEnv("PROJECT_ID"), opt)
	if err != nil {
		log.Fatal(err)
	}
	return &PubSubClient{client}
}

func (c *PubSubClient) configureTopics() *Topic {
	ctx := context.Background()
	topic := c.Topic(resizeTopic)
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		log.Fatal("Not exists topics")
	}
	return &Topic{topic}
}

func (t *Topic) PublishMessage(ctx context.Context, msg string) (string, error) {
	message, err := json.Marshal(model.StorageInformation{
		FileName:      "c2e1fdf7a30a38dca150351659fdea8e",
		FileExtension: "png",
		Directory:     "photos/users/4d2b645a-e674-4805-b31a-d6806e7ecb08",
		ImagePath:     "photos/users/4d2b645a-e674-4805-b31a-d6806e7ecb08/c2e1fdf7a30a38dca150351659fdea8e.png",
	})
	if err != nil {
		return "", err
	}
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
	})
	serverID, err := result.Get(ctx)
	if err != nil {
		return "", err
	}
	return serverID, err
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Env missing key %s", key)
	}
	return value
}
