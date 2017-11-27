package command

import (
	"encoding/json"
	"fmt"

	"github.com/oliverpool/go-chromecast"
)

type statusRequest identifiableCommand

type statusResponse struct {
	Status *chromecast.Status `json:"status"`
}

func (s statusRequest) Get(requester chromecast.Requester) (st chromecast.Status, err error) {
	response, err := requester.Request(s.Envelope, s.Payload)
	if err != nil {
		return st, err
	}

	payload := <-response
	if payload == nil {
		return st, fmt.Errorf("empty status payload")
	}

	sr := statusResponse{
		Status: &st,
	}

	err = json.Unmarshal(payload, &sr)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal into status: %s", err)
	}
	return st, err
}

var statusEnv = chromecast.Envelope{
	Source:      "sender-0",
	Destination: "receiver-0",
	Namespace:   "urn:x-cast:com.google.cast.receiver",
}

var Status = statusRequest{
	Envelope: statusEnv,
	Payload:  &chromecast.PayloadWithID{Type: "GET_STATUS"},
}
