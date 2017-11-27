package media

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command"
)

// ID chosen from https://gist.github.com/jloutsenhizer/8855258
const ID = "CC1AD845"

type App struct {
	Envelope chromecast.Envelope
	Client   chromecast.Client

	mu           sync.Mutex
	latestStatus []Status
}

func New(client chromecast.Client) (a *App, err error) {
	env, err := command.Launch.App(client, ID)
	if err != nil {
		return a, err
	}
	a = &App{
		Envelope: env,
		Client:   client,
	}

	return a, command.Connect.SendTo(client, env.Destination)
}

func FromStatus(client chromecast.Client, st chromecast.Status) (a *App, err error) {
	env, err := command.AppEnvFromStatus(st, ID, command.Status.Envelope.Source)
	if err != nil {
		return a, err
	}
	a = &App{
		Envelope: env,
		Client:   client,
	}

	return a, command.Connect.SendTo(client, env.Destination)
}

type Item struct {
	ContentId   string `json:"contentId"`
	StreamType  string `json:"streamType"`
	ContentType string `json:"contentType"`
}

type Status struct {
	SessionID              int                    `json:"mediaSessionId"`
	PlaybackRate           float64                `json:"playbackRate"`
	PlayerState            string                 `json:"playerState"`
	CurrentTime            float64                `json:"currentTime"`
	SupportedMediaCommands int                    `json:"supportedMediaCommands"`
	Volume                 *chromecast.Volume     `json:"volume,omitempty"`
	Item                   *ItemStatus            `json:"media"`
	CustomData             map[string]interface{} `json:"customData"`
	RepeatMode             string                 `json:"repeatMode"`
	IdleReason             string                 `json:"idleReason"`
}

type statusResponse struct {
	Status []Status `json:"status"`
}

type ItemStatus struct {
	ContentId   string  `json:"contentId"`
	StreamType  string  `json:"streamType"`
	ContentType string  `json:"contentType"`
	Duration    float64 `json:"duration"`
}

// FOR DEBUG ONLY!
func (a *App) syncedRequestDEBUG(payload chromecast.IdentifiablePayload) error {
	response, err := a.Client.Request(a.Envelope, payload)
	if err != nil {
		return err
	}
	<-response // FIXME: do something with the response?
	return nil
}

func (a *App) request(payload chromecast.IdentifiablePayload) (<-chan []byte, error) {
	return a.Client.Request(a.Envelope, payload)
}

func (a *App) setStatus(st []Status) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.latestStatus = st
}

func (a *App) CurrentSession() (*Session, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.firstSession(a.latestStatus)
}

func (a *App) firstSession(st []Status) (*Session, error) {
	for _, status := range st {
		if status.SessionID > 0 {
			return &Session{
				App: a,
				ID:  status.SessionID,
			}, nil
		}
	}
	return nil, fmt.Errorf("no valid SessionId has been found in the status")
}

func (a *App) Load(item Item) (*Session, error) {
	payload := command.Map{
		"type":     "LOAD",
		"media":    item,
		"autoplay": true,
		// "currentTime": 0,
		// "customData":  struct{}{},
	}
	response, err := a.Client.Request(a.Envelope, payload)
	if err != nil {
		return nil, err
	}
	body := <-response
	var s statusResponse
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	a.setStatus(s.Status)
	return a.firstSession(s.Status)
}

func (a *App) GetStatus() ([]Status, error) {
	payload := command.Map{"type": "GET_STATUS"}
	response, err := a.Client.Request(a.Envelope, payload)
	if err != nil {
		return nil, err
	}
	body := <-response
	var s statusResponse
	err = json.Unmarshal(body, &s)
	if err == nil {
		a.setStatus(s.Status)
	}
	return s.Status, err
}

func (a *App) UpdateStatus() {
	ch := make(chan []byte, 1)
	a.Client.Listen(a.Envelope, "MEDIA_STATUS", ch)

	for payload := range ch {
		var s statusResponse
		if err := json.Unmarshal(payload, &s); err != nil {
			continue
		}
		a.setStatus(s.Status)
	}
}

type Session struct {
	*App
	ID int `json:"mediaSessionId"`
}

func (s Session) do(cmd string) (<-chan []byte, error) {
	payload := command.Map{
		"type":           cmd,
		"mediaSessionId": s.ID,
	}
	return s.App.request(payload)
}

func (s Session) Pause() (<-chan []byte, error) {
	return s.do("PAUSE")
}

func (s Session) Play() (<-chan []byte, error) {
	return s.do("PLAY")
}

func (s Session) Stop() (<-chan []byte, error) {
	return s.do("STOP")
}
