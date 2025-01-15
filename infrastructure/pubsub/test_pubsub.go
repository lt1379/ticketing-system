package pubsub

import (
	"context"
	"github.com/lts1379/ticketing-system/infrastructure/logger"
	"log"

	"cloud.google.com/go/pubsub"
)

type ITestPubSub interface {
	Publish(ctx context.Context, topic string, payload []byte) (string, error)
	GetSubscription(ctx context.Context, subID string) (*pubsub.Subscription, error)
}

type TestPubSub struct {
	PubSubClient *pubsub.Client
}

func NewTestPubSub(pubSubClient *pubsub.Client) ITestPubSub {
	return &TestPubSub{
		PubSubClient: pubSubClient,
	}
}

func (testPubSub *TestPubSub) Publish(ctx context.Context, topicName string, payload []byte) (string, error) {
	msg := &pubsub.Message{
		Data: payload,
	}

	topic = testPubSub.PubSubClient.Topic(topicName)

	// Create the topic if it doesn't exist.
	exists, err := topic.Exists(ctx)
	if err != nil {
		return "", err
	}
	if !exists {
		log.Printf("Topic %v doesn't exist - creating it", topicName)
		_, err = testPubSub.PubSubClient.CreateTopic(ctx, topicName)
		if err != nil {
			return "", err
		}
	}

	serverId, err := topic.Publish(ctx, msg).Get(ctx)
	if err != nil {
		return "", err
	}

	logger.GetLogger().WithField("server ID", serverId).Info("Message published")
	return serverId, nil
}

func (testPubSub *TestPubSub) GetSubscription(ctx context.Context, subID string) (*pubsub.Subscription, error) {
	logger.GetLogger().WithField("subID", subID).Info("PubSub starting...")

	return testPubSub.PubSubClient.Subscription(subID), nil
}
