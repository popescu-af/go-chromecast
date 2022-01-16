package command

import (
	"encoding/json"
	"fmt"

	chromecast "github.com/popescu-af/go-chromecast"
)

type Launcher struct {
	Requester chromecast.Requester
}

func (l Launcher) statusRequest(pay chromecast.IdentifiablePayload) (st chromecast.Status, err error) {
	env := chromecast.Envelope{
		Source:      DefaultSource,
		Destination: DefaultDestination,
		Namespace:   "urn:x-cast:com.google.cast.receiver",
	}

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

func (l Launcher) Status() (st chromecast.Status, err error) {
	pay := Map{
		"type": "GET_STATUS",
	}
	return l.statusRequest(pay)
}

// Launch will launch the given app, except if it is found running in one of the optional statuses
func (l Launcher) Launch(appID string, statuses ...chromecast.Status) (st chromecast.Status, err error) {
	for _, st := range statuses {
		app := st.AppWithID(appID)
		if app != nil {
			return st, nil
		}
	}
	pay := Map{
		"type":  "LAUNCH",
		"appId": appID,
	}
	return l.statusRequest(pay)
}

func (l Launcher) Stop() (st chromecast.Status, err error) {
	pay := Map{
		"type": "STOP",
	}
	return l.statusRequest(pay)
}

func (l Launcher) SetVolume(level float64) (st chromecast.Status, err error) {
	vol := chromecast.Volume{
		Level: &level,
	}
	pay := Map{
		"type":   "SET_VOLUME",
		"volume": vol,
	}
	return l.statusRequest(pay)
}

func (l Launcher) Mute(muted bool) (st chromecast.Status, err error) {
	vol := chromecast.Volume{
		Muted: &muted,
	}
	pay := Map{
		"type":   "SET_VOLUME",
		"volume": vol,
	}
	return l.statusRequest(pay)
}

// AmpController returns an AmpController
func (l Launcher) AmpController() chromecast.AmpController {
	return ampLauncher{l}
}

type ampLauncher struct {
	Launcher
}

func (a ampLauncher) Mute(muted bool) error {
	_, err := a.Launcher.Mute(muted)
	return err
}
func (a ampLauncher) SetVolume(level float64) error {
	_, err := a.Launcher.SetVolume(level)
	return err
}

func (a ampLauncher) Quit() error {
	_, err := a.Launcher.Stop()
	return err
}
