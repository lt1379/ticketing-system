package pubsub

import (
	"context"
	"log"

	"sync"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

var (
	topic *pubsub.Topic

	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string

	// token is used to verify push requests.
)

type PubSubHandler struct {
	Topic        string
	Subscription string
	Handler      Handler
}

type Handler func(ctx context.Context, msg *pubsub.Message)

func NewPubSub(ctx context.Context, projectID string) (*pubsub.Client, error) {
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("./serviceAccount.json"))
	if err != nil {
		return nil, err
	}
	// defer client.Close()

	log.Printf("PubSub connected...")
	return client, err
}
