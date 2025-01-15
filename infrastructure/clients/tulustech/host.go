package tulustech

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lts1379/ticketing-system/infrastructure/clients"
	"github.com/lts1379/ticketing-system/infrastructure/clients/tulustech/models"
)

type ITulusHost interface {
	GetRandomTyping(ctx context.Context, reqHeader models.ReqHeader) (models.ResTypingRandom, error)
}

type TulusHost struct {
	id   string
	host string
}

func NewTulusHost(host string) ITulusHost {
	return &TulusHost{host: host}
}

func (TulusHost *TulusHost) GetRandomTyping(ctx context.Context, reqHeader models.ReqHeader) (models.ResTypingRandom, error) {
	var res models.ResTypingRandom

	endpoint := "/api/typings/random"
	method := "POST"

	reqMapHeader := map[string]string{
		"Accept":       reqHeader.Accept,
		"Content-Type": reqHeader.ContentType,
		"Cookie":       reqHeader.Cookie,
	}
	hostClient := clients.NewHost(TulusHost.host, endpoint, method, nil, reqMapHeader, nil)
	byteData, statusCode, err := hostClient.HTTPPost()
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(byteData, &res); err != nil {
		return res, err
	}

	if statusCode < 200 && statusCode > 299 {
		return res, fmt.Errorf("something occurred with server")
	}

	return res, nil
}
