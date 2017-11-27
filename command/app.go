package command

import (
	"encoding/json"
	"fmt"

	chromecast "github.com/oliverpool/go-chromecast"
)

type launchRequest struct {
	Envelope chromecast.Envelope
	Payload  Map
}

func (l launchRequest) App(requester chromecast.Requester, id string) (st chromecast.Status, err error) {
	l.Payload["appId"] = id
	response, err := requester.Request(l.Envelope, l.Payload)
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

func TransportForNamespace(st chromecast.Status, namespace string) (transport string, err error) {
	for _, app := range st.Applications {
		if app == nil || app.TransportId == nil {
			continue
		}
		for _, ns := range app.Namespaces {
			if ns == nil || ns.Name != namespace {
				continue
			}
			return *app.TransportId, nil
		}
	}
	return "", fmt.Errorf("no app with namespace '%s' could not be found in status", namespace)
}

var Launch = launchRequest{
	Envelope: statusEnv,
	Payload:  Map{"type": "LAUNCH"},
}
