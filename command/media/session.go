package media

import "github.com/oliverpool/go-chromecast/command"

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

func (s Session) ensureDo(cmd, state string) (<-chan bool, error) {
	req, err := s.do(cmd)
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

func (s Session) Pause() (<-chan bool, error) {
	return s.ensureDo("PAUSE", "PAUSED")
}

func (s Session) Play() (<-chan bool, error) {
	return s.ensureDo("PLAY", "PLAYING")
}

func (s Session) Stop() (<-chan []byte, error) {
	return s.do("STOP")
}

func playerStateIs(sr statusResponse, state string) bool {
	for _, s := range sr.Status {
		if s.PlayerState == state {
			return true
		}
	}
	return false
}
