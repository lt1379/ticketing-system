package servicebus

import (
	"context"
	"fmt"
	"github.com/lts1379/ticketing-system/infrastructure/logger"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type ITestServiceBus interface {
	SendMessage(message []byte) error
	GetMessage(count int)
}

type TestServicebus struct {
	AzservicebusClient *azservicebus.Client
}

func NewTestServiceBus(azServiceBusClient *azservicebus.Client) ITestServiceBus {
	return &TestServicebus{AzservicebusClient: azServiceBusClient}
}

func (testServiceBus *TestServicebus) SendMessage(message []byte) error {
	sender, err := testServiceBus.AzservicebusClient.NewSender("test-queue", nil)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while making new sender service bus.")
		return err
	}
	defer sender.Close(context.Background())

	sbMessage := &azservicebus.Message{
		Body: message,
	}
	err = sender.SendMessage(context.Background(), sbMessage, nil)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while sending message.")
		return err
	}

	return nil
}

func (testServiceBus *TestServicebus) GetMessage(count int) {
	receiver, err := testServiceBus.AzservicebusClient.NewReceiverForQueue("testqueue", nil)
	if err != nil {
		panic(err)
	}
	defer receiver.Close(context.Background())

	messages, err := receiver.ReceiveMessages(context.Background(), count, nil)
	if err != nil {
		panic(err)
	}

	for _, message := range messages {
		body := message.Body
		fmt.Printf("%s\n", string(body))

		err = receiver.CompleteMessage(context.Background(), message, nil)
		if err != nil {
			panic(err)
		}
	}
}
