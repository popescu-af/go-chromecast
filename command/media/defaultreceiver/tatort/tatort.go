package tatort

import (
	"encoding/xml"
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
	app, err := media.LaunchAndConnect(client, defaultreceiver.ID, statuses...)
	return App{app}, err
}

func (a App) Load(id string, options ...media.Option) (<-chan []byte, error) {
	item := media.Item{
		ContentID:   id,
		ContentType: "application/x-mpegurl",
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
		"www.daserste.de": struct{}{},
		"daserste.de":     struct{}{},
	}
	if _, ok := hosts[u.Host]; !ok {
		return "", fmt.Errorf("unsupported host: %s", u.Host)
	}
	u.RawQuery = ""
	u.Fragment = ""
	apiURL := getAPIURL(u.String())

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
	path = strings.TrimSuffix(path, ".html")
	return path + "~playerXml.xml"
}

func extractIDFromAPIResponse(body io.Reader) (string, error) {
	// playlist video assets(not audiodesc) asset[type="5.2.13.12.1 Web L"] fileName
	var response struct {
		XMLName xml.Name `xml:"playlist"`
		Video   struct {
			Title  string `xml:"title"`
			Assets []struct {
				Type  string `xml:"type,attr"`
				Asset []struct {
					Type     string `xml:"type,attr"`
					FileName string `xml:"fileName"`
				} `xml:"asset"`
			} `xml:"assets"`
		} `xml:"video"`
	}
	err := xml.NewDecoder(body).Decode(&response)
	if err != nil {
		return "", err
	}
	for _, assets := range response.Video.Assets {
		if assets.Type != "" {
			continue
		}
		for _, asset := range assets.Asset {
			if strings.HasSuffix(asset.FileName, ".m3u8") {
				return asset.FileName, nil
			}
		}
	}
	return "", fmt.Errorf("could not find .m3u8 file")
}
