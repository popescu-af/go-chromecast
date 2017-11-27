package media

import "github.com/oliverpool/go-chromecast/command"

type Session struct {
	*App
	ID int `json:"mediaSessionId"`
}

func (s Session) do(cmd string, options ...Option) (<-chan []byte, error) {
	payload := command.Map{
		"type":           cmd,
		"mediaSessionId": s.ID,
	}
	for _, opt := range options {
		opt(payload)
	}
	return s.App.request(payload)
}

func (s Session) doEnsure(cmd, state string, options ...Option) (<-chan bool, error) {
	req, err := s.do(cmd, options...)
	if err != nil {
		return nil, err
	}
	ch := make(chan bool, 1)
	go func() {
		for payload := range req {
			sr, err := unmarshalStatus(payload)
			ch <- (err == nil && playerStateIs(sr, state))
		}
		close(ch)
	}()

	return ch, nil
}

func (s Session) Pause(options ...Option) (<-chan bool, error) {
	return s.doEnsure("PAUSE", "PAUSED", options...)
}

func (s Session) Seek(options ...Option) (<-chan []byte, error) {
	return s.do("SEEK", options...)
}

func (s Session) Stop(options ...Option) (<-chan bool, error) {
	return s.doEnsure("STOP", "IDLE", options...)
}

func (s Session) Play(options ...Option) (<-chan bool, error) {
	return s.doEnsure("PLAY", "PLAYING", options...)
}

func playerStateIs(sr statusResponse, state string) bool {
	for _, s := range sr.Status {
		if s.PlayerState == state {
			return true
		}
	}
	return false
}
