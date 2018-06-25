package vimeo_test

import (
	"testing"

	"github.com/oliverpool/go-chromecast/command/media/vimeo"
)

func TestURLParsing(t *testing.T) {
	cc := []struct {
		url      string
		expected string
	}{
		{"https://vimeo.com/channels/staffpicks/276738707", "/videos/276738707"},
		{"https://vimeo.com/276405604", "/videos/276405604"},
	}

	for _, c := range cc {
		got, err := vimeo.ExtractID(c.url)
		if got != c.expected {
			t.Errorf("got '%s', expected '%s' for '%s'", got, c.expected, c.url)
		}
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}
	}
}
