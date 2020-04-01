package tatort

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestAPIURLGeneration(t *testing.T) {
	cc := []struct {
		path     string
		expected string
	}{
		{
			path:     "http://www.daserste.de/unterhaltung/krimi/tatort/videos/freies-land-video-100.html",
			expected: "http://www.daserste.de/unterhaltung/krimi/tatort/videos/freies-land-video-100~playerXml.xml",
		},
	}
	for _, c := range cc {
		got := getAPIURL(c.path)
		if got != c.expected {
			t.Errorf("got '%s', expected '%s' for '%s'", got, c.expected, c.path)
		}
	}
}

func TestAPIParsingGeneration(t *testing.T) {
	cc := []struct {
		body     io.Reader
		expected string
	}{
		{
			body: strings.NewReader(`<?xml version="1.0" encoding="UTF-8" ?>
				<playlist startIndex="0" title="">
				<video commentable="false">
				<videoId>freies-land-video-100</videoId>
				<title>Tatort: Freies Land</title>
				<shareTitle>Freies Land</shareTitle>
				<duration>01:28:59</duration>
				<broadcast>Tatort</broadcast><broadcastDate>2018-06-03T20:15:00.000+02:00</broadcastDate>
				<author></author>
				<permalink>http://www.daserste.de/unterhaltung/krimi/tatort/videos/freies-land-video-100.html</permalink>
				<recommendUrl></recommendUrl>
				<embeddable>false</embeddable>
				<ivwPixelUrl></ivwPixelUrl>
				<trackingPixelUrl></trackingPixelUrl>
				<agf-tracking>
				<c10>DasErste</c10>
				<c12>Content</c12>
				<c15>X004502998</c15>
				<c16>DasErste_Fiction|Film|Sonntagskrimi</c16>
				<c18>N</c18>
				<c19>br</c19>
				<c7>crid://daserste.de/tatort/afbc8935-50d8-4851-abcf-b7d168248ef3</c7>
				<c8>5339</c8>
				<c9>Tatort_Tatort: Freies Land_03.06.2018 20:15</c9>
				</agf-tracking>
				<fskRating></fskRating>
				<geoblocking></geoblocking>
				<desc>Batic und Leitmayr Batic ermitteln im Todesfall eines jungen Mannes und tauchen in die Welt der sogenannten &#034;Freiländer&#034; ein. Der polizeiliche Status der Ermittler aus München gilt nichts, weshalb dieser Fall ganz neue Grenzen aufweist.</desc>
				<previewImagesVtt url="https://previewthumbnails.daserste.de/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/preview.vtt"></previewImagesVtt>
				<dataTimedText url="//www.daserste.de/unterhaltung/krimi/tatort/videos/freies-land-video-ut102.xml">
				</dataTimedText>
				<dataTimedTextNoOffset url="//www.daserste.de/unterhaltung/krimi/tatort/videos/freies-land-video-ut102~_withoutOffset-true.xml">
				</dataTimedTextNoOffset>
				<dataTimedTextVtt url="//www.daserste.de/unterhaltung/krimi/tatort/videos/freies-land-video-ut102~_type-webvtt.vtt">
				</dataTimedTextVtt>
				<streamingUrlIPhone>
				</streamingUrlIPhone>
				<streamingUrlIPad>
				</streamingUrlIPad>
				<similarVideos>//www.daserste.de/api/unterhaltung/krimi/tatort/videos/freies-land-video-100~similarVideos.json</similarVideos>
				<relatedUrl>//www.daserste.de/api/unterhaltung/krimi/tatort/videos/freies-land-video-100~moreVideos.json</relatedUrl>
				<assets type="audiodesc">
				<asset type="1.2.3.6.1 Prog 320x180">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/320-3.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>320</frameWidth>
				<frameHeight>180</frameHeight>
				<bitrateVideo>128</bitrateVideo>
				<bitrateAudio>56</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>184</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.11.1 Web L">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/640-3.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>640</frameWidth>
				<frameHeight>360</frameHeight>
				<bitrateVideo>1024</bitrateVideo>
				<bitrateAudio>192</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>1216</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.9.1 Web M">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/512-3.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>512</frameWidth>
				<frameHeight>288</frameHeight>
				<bitrateVideo>512</bitrateVideo>
				<bitrateAudio>96</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>608</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.12.2 Web L">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/960-3.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>960</frameWidth>
				<frameHeight>540</frameHeight>
				<bitrateVideo>1800</bitrateVideo>
				<bitrateAudio>192</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>1992</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.7.1 Prog 480x270">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/480-3.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>480</frameWidth>
				<frameHeight>270</frameHeight>
				<bitrateVideo>256</bitrateVideo>
				<bitrateAudio>64</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>320</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="5.2.13.12.1 Web L">
				<serverPrefix></serverPrefix>
				<fileName>https://dasersteuni-vh.akamaihd.net/i/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/,320-3,640-3,512-3,960-3,480-3,.mp4.csmil/master.m3u8</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth></frameWidth>
				<frameHeight></frameHeight>
				<bitrateVideo></bitrateVideo>
				<bitrateAudio></bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate></totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="5.2.7.12.1 HDS VoD Web-L">
				<serverPrefix></serverPrefix>
				<fileName>https://dasersteuni-vh.akamaihd.net/z/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/,640-3,512-3,960-3,.mp4.csmil/manifest.f4m</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth></frameWidth>
				<frameHeight></frameHeight>
				<bitrateVideo></bitrateVideo>
				<bitrateAudio></bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate></totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.12.3 HbbTV 960x540">
				<serverPrefix></serverPrefix>
				<fileName>http://ctv-videos.daserste.de/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/960-3.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>960</frameWidth>
				<frameHeight>540</frameHeight>
				<bitrateVideo>1800</bitrateVideo>
				<bitrateAudio>192</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>1992</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				</assets>
				<assets>
				<asset type="1.2.3.11.1 Web L">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/640-1.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>640</frameWidth>
				<frameHeight>360</frameHeight>
				<bitrateVideo>1024</bitrateVideo>
				<bitrateAudio>192</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>1216</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.7.1 Prog 480x270">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/480-1.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>480</frameWidth>
				<frameHeight>270</frameHeight>
				<bitrateVideo>256</bitrateVideo>
				<bitrateAudio>64</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>320</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.12.2 Web L">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/960-1.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>960</frameWidth>
				<frameHeight>540</frameHeight>
				<bitrateVideo>1800</bitrateVideo>
				<bitrateAudio>192</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>1992</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.6.1 Prog 320x180">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/320-1.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>320</frameWidth>
				<frameHeight>180</frameHeight>
				<bitrateVideo>128</bitrateVideo>
				<bitrateAudio>56</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>184</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.9.1 Web M">
				<serverPrefix></serverPrefix>
				<fileName>https://pdvideosdaserste-a.akamaihd.net/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/512-1.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>512</frameWidth>
				<frameHeight>288</frameHeight>
				<bitrateVideo>512</bitrateVideo>
				<bitrateAudio>96</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>608</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="5.2.13.12.1 Web L">
				<serverPrefix></serverPrefix>
				<fileName>https://dasersteuni-vh.akamaihd.net/i/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/,640-1,480-1,960-1,320-1,512-1,.mp4.csmil/master.m3u8</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth></frameWidth>
				<frameHeight></frameHeight>
				<bitrateVideo></bitrateVideo>
				<bitrateAudio></bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate></totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="5.2.7.12.1 HDS VoD Web-L">
				<serverPrefix></serverPrefix>
				<fileName>https://dasersteuni-vh.akamaihd.net/z/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/,640-1,960-1,512-1,.mp4.csmil/manifest.f4m</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth></frameWidth>
				<frameHeight></frameHeight>
				<bitrateVideo></bitrateVideo>
				<bitrateAudio></bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate></totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				<asset type="1.2.3.12.3 HbbTV 960x540">
				<serverPrefix></serverPrefix>
				<fileName>http://ctv-videos.daserste.de/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/960-1.mp4</fileName>
				<mediaType>mp4</mediaType>
				<frameWidth>960</frameWidth>
				<frameHeight>540</frameHeight>
				<bitrateVideo>1800</bitrateVideo>
				<bitrateAudio>192</bitrateAudio>
				<codecVideo>H.264</codecVideo>
				<totalBitrate>1992</totalBitrate>
				<recommendedBandwidth></recommendedBandwidth>
				<readableSize>0 bytes</readableSize>
				<dimensions></dimensions>
				<downloadUrl></downloadUrl>
				<size>0</size>
				</asset>
				</assets><teaserImage>
				<variants>
				<variant name="image1280">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/sendung/freies-land-logo-100~_v-varxl.jpg</url>
				<width>1280</width>
				<height>720</height>
				</variant>
				<variant name="image960">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/sendung/freies-land-logo-100~_v-varvideol.jpg</url>
				<width>960</width>
				<height>540</height>
				</variant>
				<variant name="image512">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/sendung/freies-land-logo-100~_v-varm.jpg</url>
				<width>512</width>
				<height>288</height>
				</variant>
				<variant name="image256">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/sendung/freies-land-logo-100~_v-vars.jpg</url>
				<width>256</width>
				<height>144</height>
				</variant>
				</variants>
				<alttext>Batic und Leitmayr</alttext>
				</teaserImage>
				<teaserImage>
				<variants>
				<variant name="image1280">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/br-freies-land-bild-126~_v-varxl.jpg</url>
				<width>1280</width>
				<height>720</height>
				</variant>
				<variant name="image960">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/br-freies-land-bild-126~_v-varvideol.jpg</url>
				<width>960</width>
				<height>540</height>
				</variant>
				<variant name="image512">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/br-freies-land-bild-126~_v-varm.jpg</url>
				<width>512</width>
				<height>288</height>
				</variant>
				<variant name="image256">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/br-freies-land-bild-126~_v-vars.jpg</url>
				<width>256</width>
				<height>144</height>
				</variant>
				</variants>
				<alttext>Batic (Miroslav Nemec) und Leitmayr (Udo Wachtveitl) </alttext>
				</teaserImage>
				<teaserImage>
				<variants>
				<variant name="image1280">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/sendung/kriminalhauptkommissar-franz-leitmayr-und-kriminalhauptkommissar-ivo-batic-stehen-vor-dem-ver-100~_v-varxl.jpg</url>
				<width>1280</width>
				<height>720</height>
				</variant>
				<variant name="image960">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/sendung/kriminalhauptkommissar-franz-leitmayr-und-kriminalhauptkommissar-ivo-batic-stehen-vor-dem-ver-100~_v-varvideol.jpg</url>
				<width>960</width>
				<height>540</height>
				</variant>
				<variant name="image512">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/sendung/kriminalhauptkommissar-franz-leitmayr-und-kriminalhauptkommissar-ivo-batic-stehen-vor-dem-ver-100~_v-varm.jpg</url>
				<width>512</width>
				<height>288</height>
				</variant>
				<variant name="image256">
				<url>//www.daserste.de/unterhaltung/krimi/tatort/sendung/kriminalhauptkommissar-franz-leitmayr-und-kriminalhauptkommissar-ivo-batic-stehen-vor-dem-ver-100~_v-vars.jpg</url>
				<width>256</width>
				<height>144</height>
				</variant>
				</variants>
				<alttext>Kriminalhauptkommissar Franz Leitmayr (Udo Wachtveitl, links) und Kriminalhauptkommissar Ivo Batic (Miroslav Nemec) stehen vor dem verschlossenen Eingangstor zum &#034;Freiland&#034;-Gelände. Weiteres Bildmaterial finden Sie unter www.br-foto.de.</alttext>
				</teaserImage></video>
				</playlist>`),
			expected: "https://dasersteuni-vh.akamaihd.net/i/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/,640-1,480-1,960-1,320-1,512-1,.mp4.csmil/master.m3u8"},
	}
	for _, c := range cc {
		got, err := extractIDFromAPIResponse(c.body)
		if got != c.expected {
			t.Errorf("got '%s', expected '%s'", got, c.expected)
		}
		if err != nil {
			t.Errorf("got unexpected error: %w", err)
		}
	}
}

func TestURLParsing(t *testing.T) {
	if os.Getenv("TEST_ONLINE") == "" {
		t.Skip("online test skipped")
	}
	cc := []struct {
		url      string
		expected string
	}{
		{
			url:      "http://www.daserste.de/unterhaltung/krimi/tatort/videos/freies-land-video-100.html",
			expected: "https://dasersteuni-vh.akamaihd.net/i/int/2018/05/24/afbc8935-50d8-4851-abcf-b7d168248ef3/,640-1,480-1,960-1,320-1,512-1,.mp4.csmil/master.m3u8",
		},
	}

	for _, c := range cc {
		got, err := ExtractID(c.url)
		if got != c.expected {
			t.Errorf("got '%s', expected '%s' for '%s'", got, c.expected, c.url)
		}
		if err != nil {
			t.Errorf("got unexpected error: %w", err)
		}
	}
}
