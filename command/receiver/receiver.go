package receiver

import (
	"encoding/json"
	"fmt"

	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command"
)

const Namespace = "urn:x-cast:com.google.cast.receiver"

func Launch(requester chromecast.Requester, id string) (st chromecast.Status, err error) {
	env := chromecast.Envelope{
		Source:      command.DefaultSource,
		Destination: command.DefaultDestination,
		Namespace:   Namespace,
	}
	pay := command.Map{
		"type":  "LAUNCH",
		"appId": id,
	}
	response, err := requester.Request(env, pay)
	if err != nil {
		return st, err
	}

	payload := <-response
	if payload == nil {
		return st, fmt.Errorf("empty status payload")
	}

	sr := chromecast.StatusResponse{
		Status: &st,
	}

	err = json.Unmarshal(payload, &sr)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal into status: %s", err)
	}

	return st, err
}
