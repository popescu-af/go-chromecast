package vimeo

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestIframeExtraction(t *testing.T) {
	cc := []struct {
		body     io.Reader
		expected string
	}{
		{
			body: strings.NewReader(`
				<component is="lesson-view" inline-template>
				<div>
				<div class="video-player-wrap">
				<div class="video-player" v-cloak>
				<video-player lesson="1148" vimeo-id="231780045" inline-template>
				<div id="laracasts-video" class="container"></div>
				</video-player>

				<div class="next-lesson-arrow previous" v-cloak>`),
			expected: "https://player.vimeo.com/video/231780045"},
	}
	for _, c := range cc {
		got, err := extractIframeFromPage(c.body)
		if got != c.expected {
			t.Errorf("got '%s', expected '%s'", got, c.expected)
		}
		if err != nil {
			t.Errorf("got unexpected error: %w", err)
		}
	}
}
func TestMp4Extraction(t *testing.T) {
	cc := []struct {
		body     io.Reader
		expected string
	}{
		{
			body:     strings.NewReader(`,"default_cdn":"akfire_interconnect_quic","cdns":{"akfire_interconnect_quic":{"url":"https://46skyfiregce-vimeo.akamaized.net/exp=1529961460~acl=%2F231780045%2F%2A~hmac=3e9c6fb6936f69d51b891d6a4213bec94f4efa6e2640dd41acad57344861b3af/231780045/video/820169965,820170024,820170021,820169954/master.m3u8","origin":"gcs"},"fastly_skyfire":{"url":"https://skyfire.vimeocdn.com/1529961460-0x1e04e1b6f2d57441c4e46b0eb86067758218cb97/231780045/video/820169965,820170024,820170021,820169954/master.m3u8","origin":"gcs"}}},"progressive":[{"profile":174,"width":1280,"mime":"video/mp4","fps":30,"url":"https://gcs-vimeo.akamaized.net/exp=1529961460~acl=%2A%2F820170024.mp4%2A~hmac=1b0c809bc92d1924a50ae061b2cc633f08087ccfafdedfbf53e5200e3250cce9/vimeo-prod-skyfire-std-us/01/1356/9/231780045/820170024.mp4","cdn":"akamai_interconnect","quality":"720p","id":820170024,"origin":"gcs","height":720},{"profile":175,"width":1920,"mime":"video/mp4","fps":30,"url":"https://gcs-vimeo.akamaized.net/exp=1529961460~acl=%2A%2F820170021.mp4%2A~hmac=8cc4c2cb65f4269a693c4de059bd74d5f2797057ed73bc8cd7d9b0c4dc0582df/vimeo-prod-skyfire-std-us/01/1356/9/231780045/820170021.mp4","cdn":"akamai_interconnect","quality":"1080p","id":820170021,"origin":"gcs","height":1080},{"profile":164,"width":640,"mime":"video/mp4","fps":30,"url":"https://gcs-vimeo.akamaized.net/exp=1529961460~acl=%2A%2F820169965.mp4%2A~hmac=1486b9fbe03f66d405e8e2c3b651518185d06d282477c7d2c2c715d9fa1a6e52/vimeo-prod-skyfire-std-us/01/1356/9/231780045/820169965.mp4","cdn":"akamai_interconnect","quality":"360p","id":820169965,"origin":"gcs","height":360},{"profile":165,"width":960,"mime":"video/mp4","fps":30,"url":"https://gcs-vimeo.akamaized.net/exp=1529961460~acl=%2A%2F820169954.mp4%2A~hmac=df91aed8aa7cc15851e44ecbdbf3baa4906d6c9213539d51f6e7d5e6e73bf88c/vimeo-prod-skyfire-std-us/01/1356/9/231780045/820169954.mp4","cdn":"akamai_interconnect","quality":"540p","id":820169954,"origin":"gcs","height":540}]},"lang":"en","sentry":{"url":"https://9e9ab33f1870463393a4a1e85a1280c2@sentry.cloud.vimeo.com/2","enabled":false,"debug_enab`),
			expected: "https://46skyfiregce-vimeo.akamaized.net/exp=1529961460~acl=%2F231780045%2F%2A~hmac=3e9c6fb6936f69d51b891d6a4213bec94f4efa6e2640dd41acad57344861b3af/231780045/video/820169965,820170024,820170021,820169954/master.m3u8",
		},
		{
			body:     strings.NewReader(`8","origin":"gcs"},"fastly_skyfire":{"url":"https://skyfire.vimeocdn.com/1529964151-0x626e27b83ed503bb96417a8a85643ad5106742a3/238821524/video/853435837,853435913,853435901,853435825/master.m3u8?token","origin":"gcs"}}},"progressive":[{"profile":175,"width":1920,"mime":"video/mp4","fps":30,"url":"https://fpdl.vimeocdn.com/vimeo-prod-skyfire-std-us/01/2764/9/238821524/853435913.mp4?token=1529964151-0x901607a1e3b3243c38d31be12c27d71548916b7b","cdn":"fastly","quality":"1080p","id":853435913,"origin":"gcs","height":1080},{"profile":165,"width":960,"mime":"video/mp4","fps":30,`),
			expected: "https://skyfire.vimeocdn.com/1529964151-0x626e27b83ed503bb96417a8a85643ad5106742a3/238821524/video/853435837,853435913,853435901,853435825/master.m3u8?token",
		},
	}
	for _, c := range cc {
		got, err := extractM3u8FromIframe(c.body)
		if got != c.expected {
			t.Errorf("got '%s', expected '%s'", got, c.expected)
		}
		if err != nil {
			t.Errorf("got unexpected error: %w", err)
		}
	}
}

func TestPageRequest(t *testing.T) {
	if os.Getenv("TEST_ONLINE") == "" {
		t.Skip("online test skipped")
	}
	cc := []struct {
		url    string
		iframe string
	}{
		{
			url:    "https://laracasts.com/series/whats-new-in-laravel-5-5/episodes/20",
			iframe: "https://player.vimeo.com/video/231780045",
		},
	}

	for _, c := range cc {
		got, err := ExtractIframe(c.url)
		if got != c.iframe {
			t.Errorf("got '%s', expected '%s' for '%s'", got, c.iframe, c.url)
		}
		if err != nil {
			t.Errorf("got unexpected error: %w", err)
		}
	}
}

func TestIframeRequest(t *testing.T) {
	if os.Getenv("TEST_ONLINE") == "" {
		t.Skip("online test skipped")
	}
	cc := []struct {
		url       string
		iframe    string
		mp4prefix string
	}{
		{
			url:       "https://laracasts.com/series/whats-new-in-laravel-5-5/episodes/20",
			iframe:    "https://player.vimeo.com/video/231780045",
			mp4prefix: "https://gcs-vimeo.akamaized.net/",
		},
	}

	for _, c := range cc {
		got, err := ExtractM3u8(c.url, c.iframe)
		if !strings.HasPrefix(got, c.mp4prefix) {
			t.Errorf("got '%s', expected '%s' prefix for '%s'", got, c.mp4prefix, c.url)
		}
		if err != nil {
			t.Errorf("got unexpected error: %w", err)
		}
	}
}
