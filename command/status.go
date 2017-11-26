package command

import (
	"encoding/json"
	"fmt"

	cast "github.com/oliverpool/go-chromecast"
)

type statusRequest identifiableCommand

type statusResponse struct {
	Status *cast.Status `json:"status"`
}

func (s statusRequest) Get(requester requestFunc) (st cast.Status, err error) {
	response, err := requester(s.Envelope, s.Payload)
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

var statusEnv = cast.Envelope{
	Source:      "sender-0",
	Destination: "receiver-0",
	Namespace:   "urn:x-cast:com.google.cast.receiver",
}

var Status = statusRequest{
	Envelope: statusEnv,
	Payload:  &cast.PayloadWithID{Type: "GET_STATUS"},
}
