package media

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command"
)

// https://developers.google.com/cast/docs/reference/messages

const Namespace = "urn:x-cast:com.google.cast.media"

type App struct {
	Envelope chromecast.Envelope
	Client   chromecast.Client

	mu           sync.Mutex
	latestStatus []Status
}

func FromStatus(client chromecast.Client, st chromecast.Status) (a *App, err error) {
	transport, err := command.TransportForNamespace(st, Namespace)
	if err != nil {
		return a, err
	}
	a = &App{
		Envelope: chromecast.Envelope{
			Source:      command.DefaultSource,
			Destination: transport,
			Namespace:   Namespace,
		},
		Client: client,
	}

	return a, command.Connect.SendTo(client, a.Envelope.Destination)
}

type Item struct {
	ContentID   string `json:"contentId"`
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

func (a *App) LatestStatus() []Status {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.latestStatus
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

// Option to customize the loading
type Option func(command.Map)

func PreventAutoplay(c command.Map) {
	c["autoplay"] = false
}

func Seek(t time.Duration) func(command.Map) {
	return func(c command.Map) {
		c["currentTime"] = t.Seconds()
	}
}

func CustomData(data interface{}) func(command.Map) {
	return func(c command.Map) {
		c["customData"] = data
	}
}

func (a *App) Load(item Item, options ...Option) (*Session, error) {
	payload := command.Map{
		"type":  "LOAD",
		"media": item,
	}
	for _, opt := range options {
		opt(payload)
	}
	response, err := a.Client.Request(a.Envelope, payload)
	if err != nil {
		return nil, err
	}
	body := <-response
	s, err := unmarshalStatus(body)
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

	s, err := unmarshalStatus(body)
	if err == nil {
		a.setStatus(s.Status)
	}
	return s.Status, err
}

func (a *App) UpdateStatus() {
	ch := make(chan []byte, 1)
	env := chromecast.Envelope{
		Source:      a.Envelope.Destination,
		Destination: a.Envelope.Source,
		Namespace:   a.Envelope.Namespace,
	}
	a.Client.Listen(env, "MEDIA_STATUS", ch)

	for payload := range ch {
		s, err := unmarshalStatus(payload)
		if err != nil {
			continue
		}
		a.setStatus(s.Status)
	}
}

func unmarshalStatus(payload []byte) (s statusResponse, err error) {
	err = json.Unmarshal(payload, &s)
	return s, err
}
