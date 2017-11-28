package receiver

import (
	"encoding/json"
	"fmt"

	chromecast "github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command"
)

const Namespace = "urn:x-cast:com.google.cast.receiver"

var env = chromecast.Envelope{
	Source:      command.DefaultSource,
	Destination: command.DefaultDestination,
	Namespace:   Namespace,
}

type Launcher struct {
	Requester chromecast.Requester
}

func (l Launcher) Launch(appID string) (st chromecast.Status, err error) {
	pay := command.Map{
		"type":  "LAUNCH",
		"appId": appID,
	}
	return l.statusRequest(pay)
}

func (l Launcher) Stop() (st chromecast.Status, err error) {
	pay := command.Map{
		"type": "STOP",
	}
	return l.statusRequest(pay)
}

func (l Launcher) statusRequest(pay chromecast.IdentifiablePayload) (st chromecast.Status, err error) {
	response, err := l.Requester.Request(env, pay)
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
