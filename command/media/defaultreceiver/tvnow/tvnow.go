package tvnow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oliverpool/go-chromecast"
	"github.com/oliverpool/go-chromecast/command/media"
	"github.com/oliverpool/go-chromecast/command/media/defaultreceiver"
)

type App struct {
	*media.App
}

func LaunchAndConnect(client chromecast.Client, statuses ...chromecast.Status) (App, error) {
	app, err := defaultreceiver.LaunchAndConnect(client, statuses...)
	return App{app}, err
}

func (a App) Load(id string, options ...media.Option) (<-chan []byte, error) {
	item := media.Item{
		ContentID:   id,
		ContentType: "application/dash+xml",
		StreamType:  "BUFFERED",
	}
	return a.App.Load(item, options...)
}

func (a App) LoadURL(rawurl string, options ...media.Option) (<-chan []byte, error) {
	id, err := ExtractID(rawurl)
	if err != nil {
		return nil, err
	}
	return a.Load(id, options...)
}

func URLLoader(rawurl string, options ...media.Option) (func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error), error) {
	id, err := ExtractID(rawurl)
	if err != nil {
		return nil, err
	}
	return func(client chromecast.Client, statuses ...chromecast.Status) (<-chan []byte, error) {
		app, err := LaunchAndConnect(client, statuses...)
		if err != nil {
			return nil, err
		}
		return app.Load(id, options...)
	}, nil
}

func ExtractID(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("could not parse url '%s': %v", rawurl, err)
	}

	hosts := map[string]struct{}{
		"www.tvnow.de": struct{}{},
		"tvnow.de":     struct{}{},
	}
	if _, ok := hosts[u.Host]; !ok {
		return "", fmt.Errorf("unsupported host: %s", u.Host)
	}
	apiURL := getAPIURL(u.Path)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("could not fetch api url '%s': %v", apiURL, err)
	}

	defer resp.Body.Close()
	id, err := extractIDFromAPIResponse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could extract ID from api response '%s': %v", apiURL, err)
	}
	if id != "" {
		return id, nil
	}
	return "", fmt.Errorf("could not find id inside URL: %s", rawurl)
}

func getAPIURL(path string) string {
	parts := strings.Split(path, "/")
	if parts[len(parts)-1] == "player" {
		parts = parts[:len(parts)-1]
	}
	if parts[0] == "" {
		parts = parts[1:]
	}
	return "https://api.tvnow.de/v3/movies/" + strings.Join(parts[1:], "/") + "?fields=*,format,files,manifest,breakpoints,paymentPaytypes,trailers,packages,isLiveStream&station=" + parts[0]
}

func extractIDFromAPIResponse(body io.Reader) (string, error) {
	var response struct {
		Manifest struct {
			Dash string
		}
	}
	err := json.NewDecoder(body).Decode(&response)
	if err != nil {
		return "", err
	}
	return response.Manifest.Dash, nil
}
