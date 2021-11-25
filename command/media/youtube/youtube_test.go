package youtube_test

import (
	"testing"

	"github.com/popescu-af/go-chromecast/command/media/youtube"
)

func TestURLParsing(t *testing.T) {
	cc := []struct {
		url      string
		expected string
	}{
		{"https://www.youtube.com/watch?v=b-GIBLX3nAk", "b-GIBLX3nAk"},
		{"https://youtu.be/b-GIBLX3nAk", "b-GIBLX3nAk"},
		{"https://youtu.be/b-GIBLX3nAk?t=1s", "b-GIBLX3nAk"},
		{"https://www.youtube.com/embed/b-GIBLX3nAk", "b-GIBLX3nAk"},
		{"https://www.youtube-nocookie.com/embed/b-GIBLX3nAk?start=10", "b-GIBLX3nAk"},
	}

	for _, c := range cc {
		got, err := youtube.ExtractID(c.url)
		if got != c.expected {
			t.Errorf("got '%s', expected '%s' for '%s'", got, c.expected, c.url)
		}
		if err != nil {
			t.Errorf("got unexpected error: %w", err)
		}
	}
}
