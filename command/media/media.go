package media

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/popescu-af/go-chromecast"
	"github.com/popescu-af/go-chromecast/command"
)

// https://developers.google.com/cast/docs/reference/messages

const Namespace = "urn:x-cast:com.google.cast.media"

type App struct {
	*command.App

	mu           sync.Mutex
	latestStatus []Status
}

func LaunchAndConnect(client chromecast.Client, id string, statuses ...chromecast.Status) (*App, error) {
	a, err := command.LaunchAndConnect(client, id, statuses...)
	if err != nil {
		return nil, err
	}
	a.Envelope.Namespace = Namespace
	return &App{
		App: a,
	}, nil
}

func ConnectFromStatus(client chromecast.Client, st chromecast.Status) (*App, error) {
	a, err := command.ConnectFromStatus(client, st, Namespace)
	if err != nil {
		return nil, err
	}
	return &App{
		App: a,
	}, nil
}

//// EDIT by me !!!!
type Track struct {
	Name        string   `json:"name"`             // Human-readable name
	TrackID     int      `json:"trackId"`          // Unique ID in the context of the parent MediaInformation object
	Type        string   `json:"type"`             // e.g. TEXT
	ContentID   string   `json:"trackContentId"`   // URL
	ContentType string   `json:"trackContentType"` // e.g. VTT or text/vtt
	Subtype     string   `json:"subtype"`          // e.g. SUBTITLES
	Roles       []string `json:"roles"`            // e.g. TEXT: main, alternate, subtitle, supplementary, commentary, dub, description, forced_subtitle
	Language    string   `json:"language"`         // e.g. en-US
	IsInband    bool     `json:"isInband"`         // embedded in the media or not
}

type TextTrackStyle struct {
	BackgroundColor   string  `json:"backgroundColor"`
	EdgeColor         string  `json:"edgeColor"`         // RGBA, e.g. ff0000ff for opaque red
	EdgeType          string  `json:"edgeType"`          // e.g. OUTLINE
	FontFamily        string  `json:"fontFamily"`        // e.g. ARIAL
	FontGenericFamily string  `json:"fontGenericFamily"` // e.g. SANS_SERIF
	FontScale         float64 `json:"fontScale"`         // e.g. 1, 0.5
	FontStyle         string  `json:"fontStyle"`         // e.g. NORMAL
	ForegroundColor   string  `json:"foregroundColor"`
}

////

type Item struct {
	ContentID   string `json:"contentId"`
	StreamType  string `json:"streamType"`
	ContentType string `json:"contentType"`

	//// EDIT by me !!!!
	Tracks             []Track        `json:"tracks"`
	SubtitleTrackStyle TextTrackStyle `json:"textTrackStyle"`
	////
}

type Status struct {
	SessionID              int                    `json:"mediaSessionId"`
	PlaybackRate           float64                `json:"playbackRate"`
	PlayerState            string                 `json:"playerState"`
	CurrentTime            Seconds                `json:"currentTime"`
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
	ContentId   string                 `json:"contentId"`
	StreamType  string                 `json:"streamType"`
	ContentType string                 `json:"contentType"`
	Duration    Seconds                `json:"duration"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type Seconds struct {
	time.Duration
}

func (s *Seconds) UnmarshalJSON(b []byte) (err error) {
	var seconds float64
	err = json.Unmarshal(b, &seconds)
	s.Duration = time.Duration(seconds * float64(time.Second))
	return err
}

func (s Seconds) MarshalJSON() (b []byte, err error) {
	return json.Marshal(s.Duration.Seconds())
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

func (a *App) Load(item Item, options ...Option) (<-chan []byte, error) {
	payload := command.Map{
		"type":  "LOAD",
		"media": item,
		//// EDIT by me !!!!
		"activeTrackIds": []int{0},
		////
	}
	for _, opt := range options {
		opt(payload)
	}
	return a.Client.Request(a.Envelope, payload)
}

func (a *App) LoadAndGetSession(item Item, options ...Option) (*Session, error) {
	response, err := a.Load(item, options...)
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

func (a *App) Status() ([]Status, error) {
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

type URLLoader func(rawurl string, options ...Option) (func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error), error)
