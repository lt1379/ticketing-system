package usecase

import (
	"context"
	"encoding/json"
	"my-project/domain/dto"
	tulushost "my-project/infrastructure/clients/tulustech"
	"my-project/infrastructure/clients/tulustech/models"
	"my-project/infrastructure/logger"
	"my-project/infrastructure/pubsub"
	"my-project/infrastructure/servicebus"
)

type ITestUsecase interface {
	Test(ctx context.Context) dto.TestDto
}

type TestUsecase struct {
	TulusTechHost  tulushost.ITulusHost
	TestPubSub     pubsub.ITestPubSub
	TestServiceBus servicebus.ITestServiceBus
}

func NewTestUsecase(tulusTechHost tulushost.ITulusHost, testPubSub pubsub.ITestPubSub, testServiceBus servicebus.ITestServiceBus) ITestUsecase {
	return &TestUsecase{TulusTechHost: tulusTechHost, TestPubSub: testPubSub, TestServiceBus: testServiceBus}
}

func (testUsecase *TestUsecase) Test(ctx context.Context) dto.TestDto {
	var res dto.TestDto

	res.PubSub = "Not OK"
	res.ServiceBus = "Not OK"

	msg := "Hello"
	byteMsg, err := json.Marshal(msg)
	if err != nil {
		logger.GetLogger().Error("Error while marshalling")
		return res
	}
	publishResponse, err := testUsecase.TestPubSub.Publish(ctx, "topic", byteMsg)
	if err != nil {
		logger.GetLogger().Error("Error while publishing message")
		res.PubSub = err.Error()
		return res
	}
	logger.GetLogger().WithField("publishResponse", publishResponse).Info("Successfully published")
	res.PubSub = "OK"

	err = testUsecase.TestServiceBus.SendMessage(byteMsg)
	if err != nil {
		logger.GetLogger().Error("Error while publishing message with service bus")
		res.ServiceBus = err.Error()
		return res
	}
	res.ServiceBus = "OK"

	reqHeader := models.ReqHeader{}
	randomTypingRes, err := testUsecase.TulusTechHost.GetRandomTyping(ctx, reqHeader)
	if err != nil {
		logger.GetLogger().Error("Error while get random typing")
		res.ServiceBus = err.Error()
		return res
	}
	logger.GetLogger().WithField("randomTypingResponse", randomTypingRes).Info("Successfully get random typing")
	res.TulusTech = "OK"

	return res
}
