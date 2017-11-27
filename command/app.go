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

func (l launchRequest) App(requester chromecast.Requester, id string) (env chromecast.Envelope, err error) {
	l.Payload["appId"] = id
	response, err := requester.Request(l.Envelope, l.Payload)
	if err != nil {
		return env, err
	}

	payload := <-response
	if payload == nil {
		return env, fmt.Errorf("empty status payload")
	}

	var st chromecast.Status
	sr := statusResponse{
		Status: &st,
	}

	err = json.Unmarshal(payload, &sr)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal into status: %s", err)
	}

	return AppEnvFromStatus(st, id, l.Envelope.Source)
}

func AppEnvFromStatus(st chromecast.Status, appID, source string) (env chromecast.Envelope, err error) {
	for _, app := range st.Applications {
		if app != nil && app.AppID != nil && *app.AppID == appID {
			if app.TransportId != nil && len(app.Namespaces) > 0 {
				env = chromecast.Envelope{
					Source:      source,
					Destination: *app.TransportId,
					Namespace:   app.Namespaces[len(app.Namespaces)-1].Name,
				}
				return env, nil
			}
			return env, fmt.Errorf("transportId or namespaces are empty")
		}
	}
	return env, fmt.Errorf("appId could not be found in status")
}

var Launch = launchRequest{
	Envelope: statusEnv,
	Payload:  Map{"type": "LAUNCH"},
}
