package tvnow

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
			path:     "/vox/die-pferdeprofis/wallach-watabi-schimmelstute-naira-2018-06-02-19-10-00/player",
			expected: "https://api.tvnow.de/v3/movies/die-pferdeprofis/wallach-watabi-schimmelstute-naira-2018-06-02-19-10-00?fields=*,format,files,manifest,breakpoints,paymentPaytypes,trailers,packages,isLiveStream&station=vox",
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
			body:     strings.NewReader(`{"orderWeight":1000,"id":571634,"title":"Wallach \"Watabi\" \/ Schimmelstute \"Naira\"","free":true,"uhdEnabled":null,"broadcastStartDate":"2018-06-02 19:10:00","broadcastPreviewStartDate":"2018-05-30 20:05:00","articleShort":"Nach seiner Genesung nimmt \"Watabi\" mit voller Kraft Reißaus und schleift seine Pflegerin über den Platz. Auch Julia Klefisch hat Probleme mit ihrem neuen Pferd. \"Naira\" rempelt, bleibt nicht stehen, bestimmt selbst, wo sie hingehen will.","articleLong":"\"Watabi\" kam mit einer gebrochenen Schulter beim Pferdeschutzhof in Köln an. Nach der Genesung nimmt er mit voller Kraft Reißaus und schleift seine Pflegerin über den Platz. Ihn anzureiten, ist derzeit ein utopischer Wunsch. Derweil hat Julia Klefisch Probleme mit ihrem neuen Pferd. \"Naira\" rempelt, bleibt nicht stehen, bestimmt selbst, wo sie hingehen will. Für Bernd Hackl ist schnell klar: \"Naira\" muss ins Trainingslager nach Bischofsreut.","teaserText":"\"Watabi\" und \"Naira\" lassen sich nicht anreiten","exclusive":false,"episode":6,"seoUrl":"wallach-watabi-schimmelstute-naira-2018-06-02-19-10-00","mamiArticleId":226138,"season":4,"replaceMovieInformation":571634,"duration":"00:51:20","dontCall":false,"productPlacement":false,"productPlacementType":false,"ignoreAd":false,"ignoreChapters":false,"archiveId":"22000735518","containerId":554400,"kenburns":true,"leadTeaserTrailer":false,"cornerLogo":"vox","aspectRatio":"16:9","format":{"id":1786,"title":"Die Pferdeprofis","titleGroup":"D","station":"vox","icon":"free","hasFreeEpisodes":true,"hasPayEpisodes":true,"seoUrl":"die-pferdeprofis","tabSeason":true,"tabSeasonStartDate":"2015-08-19 12:20:16","tabSeasonEndDate":"2025-01-01 00:00:00","tabSpecialTeaserPosition":"none","formatTabs":null,"formatPaymentProductTeasers":null,"cmsLeadTeaser":null,"defaultHeaderSmallImage":"http:\/\/bilder.rtl.de\/rtlnow\/adminupload\/die_pferdeprofis_formatseiten_header_tvnow_620x169.jpg","defaultHeaderMediumImage":"http:\/\/bilder.rtl.de\/rtlnow\/adminupload\/die_pferdeprofis_formatseiten_header_tvnow_940x169.jpg","defaultHeaderLargeImage":"http:\/\/bilder.rtl.de\/rtlnow\/adminupload\/die_pferdeprofis_formatseiten_header_tvnow_1260x169.jpg","defaultDvdImage":"http:\/\/bilder.rtl.de\/rtlnow\/adminupload\/die_pferdeprofis_default_dvd_cover_teaser_300x358.jpg","defaultDvdCoverImage":null,"defaultBannerImage":"http:\/\/bilder.rtl.de\/rtlnow\/adminupload\/die_pferdeprofis_default_banner_teaser_620x169.jpg","chapterPlayer":"1","hbbtvCutLogic":"usedefaults","videoPlazaAdTag":"now_voxnow\/die_pferdeprofis","genre1":"Docutainment","genre2":"Ratgeber","genres":["Docutainment","Ratgeber","Tiere"],"genreUris":["docutainment","ratgeber","tiere"],"kenburns":true,"flashHds":"1","movieCount":0,"categoryId":"serie","taggingSubSite":"pferdeprofis","taggingPack":"tiere","szm":{"format":{"angebot":"tvnow","ivw":"\/vox_now\/die_pferdeprofis\/home","agof2_bild":"dbrsowf_ten_vnpferdeprofis","agof2_video":"dvrsowf_ten_vnpferdeprofis","name":"die_pferdeprofis"},"detail":{"angebot":"tvnow","ivw":"\/vox_now\/die_pferdeprofis\/detail","agof2_bild":"dbrsowf_ten_vnpferdeprofis","agof2_video":"dvrsowf_ten_vnpferdeprofis","name":"die_pferdeprofis"},"player":{"angebot":"tvnow","ivw":"\/vox_now\/die_pferdeprofis\/player","agof2_bild":"dbrsowf_ten_vnpferdeprofis","agof2_video":"dvrsowf_ten_vnpferdeprofis","name":"die_pferdeprofis"}},"searchAliasName":"","enableAd":true,"enablePayAd":false,"leadTeaserTrailer":false,"metaTitle":"Die Pferdeprofis online schauen als Stream bei TV NOW","metaDescription":"Pferde sind nicht nur stolze und ehrfurchtgebietende Tiere, sie haben auch ihren eigenen Kopf. Davon kann so mancher Pferdebesitzer ein Lied singen. Und auch bei der Erziehung der Vierbeiner kann einiges schiefgehen. Rat und Hilfe gibt hier die zweiteilig","metaTags":"video, videos, online sehen, internet tv, fernsehen, video on demand, alle folgen, vod, vox, die pferdeprofis, coaching-doku, pferdetrainer, bernd hackl, sandra schneider, dokumentation, pferde","replaceFormatInformation":1786,"showMovieCount":"0","hideEpisodeNumber":false,"tabInTVDate":"00\/0000","annualNavigation":null,"emptyListText":"","headerLinkText":"Die Pferdeprofis","headerLinkSrc":"http:\/\/www.vox.de\/cms\/sendungen\/die-pferdeprofis.html","headerLinkStation":"vox","airplayEnabled":true,"availableForWeb":true,"availableForMobile":true,"availableForDevice":null,"deeplinkUrl":"http:\/\/link.tvnow.de\/?f=1786","titleShort":"Die Pferdeprofis","defaultImage169Format":"http:\/\/bilder.rtl.de\/rtlnow\/adminupload\/die_pferdeprofis_default_autoimage_teaser_1280x720.jpg","defaultImage169Logo":"http:\/\/bilder.rtl.de\/rtlnow\/adminupload\/die_pferdeprofis_default_logo_768x432.jpg","defaultLogo":"http:\/\/bilder.rtl.de\/rtlnow\/adminupload\/Die_Pferdeprofis_LOGO_Shape_150x97.png","defaultLogoHd":"http:\/\/bilder.rtl.de\/rtlnow\/adminupload\/Die_Pferdeprofis_LOGO_Shape_768x432.png","defaultExternalSquareLogo":null,"ivwPlayer":"die_pferdeprofis\/player","agof":"ren_vodpferdp","agof2":"","sellPreview":1,"moveWholeWeek":"0","justmissedDuration":30,"deleteBy":"nodelete","productImportInformation":null,"rtl":false,"vox":true,"superRtl":false,"rtl2":false,"ntv":false,"nitro":false,"tvNow":false,"crime":false,"passion":false,"living":false,"rtlPlus":false,"toggoPlus":false,"infoText":"Die Pferdeprofis kümmern sich um tierische Problemfälle.","infoTextLong":"Die Pferdeprofis kümmern sich um tierische Problemfälle, die ihre Besitzer gehörig auf Trab halten.","onlineDate":"2012-02-06 12:08:33","formatTech":"DIEPFERDEPROFIS","downscaling":false,"platform":["tvnow"],"genreDetails":null,"tags":null,"ivwFormat":"die_pferdeprofis\/home","highlight":"0","productionYear":null,"productionCountry":null,"cast":null,"metadata":null,"fsk":null,"descriptionLong":null,"sexy":null,"offlineDate":"2020-02-06 12:08:33","geoAllowedCountries":null,"isGeoblocked":false,"deeplink":null,"uhdEnabled":null,"portability":{}},"files":{"items":[{"id":3368954,"bitrate":1528,"audioBitrate":128,"videoBitrate":1400,"path":"\/abr\/mediaserverfra\/v-22\/f-1786\/c-226138\/v_226138_abr1500_G5WjCF.mp4","type":"video\/x-abr","drmmetadata":0,"isDrm":false,"s3":null},{"id":3368956,"bitrate":928,"audioBitrate":128,"videoBitrate":800,"path":"\/abr\/mediaserverfra\/v-22\/f-1786\/c-226138\/v_226138_abr1000_la2aGR.mp4","type":"video\/x-abr","drmmetadata":0,"isDrm":false,"s3":null},{"id":3368958,"bitrate":286,"audioBitrate":56,"videoBitrate":230,"path":"\/abr\/mediaserverfra\/v-22\/f-1786\/c-226138\/v_226138_abr300_gr8vAn.mp4","type":"video\/x-abr","drmmetadata":0,"isDrm":false,"s3":null},{"id":3368960,"bitrate":136,"audioBitrate":56,"videoBitrate":80,"path":"\/abr\/mediaserverfra\/v-22\/f-1786\/c-226138\/v_226138_abr100_2SrrbF.mp4","type":"video\/x-abr","drmmetadata":0,"isDrm":false,"s3":null},{"id":3368964,"bitrate":556,"audioBitrate":56,"videoBitrate":500,"path":"\/abr\/mediaserverfra\/v-22\/f-1786\/c-226138\/v_226138_abr550_guP2pt.mp4","type":"video\/x-abr","drmmetadata":0,"isDrm":false,"s3":null}],"total":5},"alternativeFiles":null,"paymentPaytypes":{"items":[{"id":2708,"name":"Die Pferdeprofis PPV","price":"0.99","currency":"EUR","period":"none","type":"ppv","geoblockingAllowedCountries":"*","format":null,"activeProviders":null,"kdgProduct":"","umkbwProduct":""}],"total":0},"pictures":null,"picture":null,"breakpoints":{"items":[{"id":20379922,"extId":-1,"frames":0,"time":"00:00:00:00","type":"S","movie":null},{"id":20379924,"extId":-1,"frames":11735,"time":"00:07:49:10","type":"a","movie":null},{"id":20379926,"extId":-1,"frames":25297,"time":"00:16:51:22","type":"s","movie":null},{"id":20379928,"extId":-1,"frames":43549,"time":"00:29:01:24","type":"a","movie":null},{"id":20379930,"extId":-1,"frames":53998,"time":"00:35:59:23","type":"s","movie":null},{"id":20379932,"extId":-1,"frames":68040,"time":"00:45:21:15","type":"a","movie":null},{"id":20379934,"extId":-1,"frames":76300,"time":"00:50:52:00","type":"e","movie":null},{"id":20379936,"extId":-1,"frames":77013,"time":"00:51:20:13","type":"E","movie":null}],"total":8},"trailers":null,"blockadeText":"0","isDrm":false,"timeType":2,"refplanningId":452644926,"geoblocked":false,"geoIpPath":null,"timecode":0,"watched":0,"paytype":0,"favourite":false,"payed":false,"manifest":{"hls":"https:\/\/api.tvnow.de\/v3\/movies\/571634\/manifest-hls.m3u8?faxs=1","hlsfairplay":"https:\/\/vodnowusohls.secure.footprint.net\/proxy\/clear\/manifest\/tvnow\/571634-1-36465.ism\/fairplay.m3u8?filter=(type%3d%3d%22audio%22)%7c%7c(type%3d%3d%22video%22%26%26systemBitrate%3c1550000)","hlsclear":"https:\/\/vodnowusohls.secure.footprint.net\/proxy\/clear\/manifest\/tvnow\/571634-1-36465.ism\/fairplay.m3u8?filter=(type%3d%3d%22audio%22)%7c%7c(type%3d%3d%22video%22%26%26systemBitrate%3c1550000)","hds":"https:\/\/api.tvnow.de\/v3\/movies\/571634\/manifest-hds.f4m","dash":"https:\/\/vodnowusodash.secure.footprint.net\/proxy\/clear\/manifest\/tvnow\/571634-1-36465.ism\/.mpd?filter=(type%3d%3d%22audio%22)%7c%7c(type%3d%3d%22video%22%26%26systemBitrate%3c1550000)","windows":"https:\/\/vodnowusohss-a.akamaihd.net\/proxy\/manifest\/tvnow\/571634-1.ism\/manifest?filter=%28type%3D%3D%22audio%22%26%26systemBitrate%3E90000%29%7C%7C%28type%3D%3D%22video%22%26%26systemBitrate%3E600000%29","chromecast":"https:\/\/vodnowusohss-a.akamaihd.net\/proxy\/manifest\/tvnow\/571634-1.ism\/manifest?filter=%28type%3D%3D%22audio%22%26%26systemBitrate%3E90000%29%7C%7C%28type%3D%3D%22video%22%26%26systemBitrate%3E600000%29","fire":"https:\/\/vodnowusodash.secure.footprint.net\/proxy\/clear\/manifest\/tvnow\/571634-1-36465.ism\/.mpd","hbbtv":null,"dashclear":"https:\/\/vodnowusodash.secure.footprint.net\/proxy\/clear\/manifest\/tvnow\/571634-1-36465.ism\/.mpd?filter=(type%3d%3d%22audio%22)%7c%7c(type%3d%3d%22video%22%26%26systemBitrate%3c1550000)","smoothclear":"https:\/\/vodnowusohss-a.akamaihd.net\/proxy\/clear\/manifest\/tvnow\/571634-1-36465.ism\/manifest?filter=(type%3d%3d%22audio%22)%7c%7c(type%3d%3d%22video%22%26%26systemBitrate%3c1550000)"},"previewStartDate":"2018-05-30 20:05:00","previewEndDate":"2018-06-02 20:00:00","catchupStartDate":"2018-06-02 20:00:00","catchupEndDate":"2018-07-02 20:15:00","archiveStartDate":"2018-07-02 20:15:00","archiveEndDate":"2030-01-01 00:00:00","licenceStartDate":"2018-05-29 15:20:17","licenceEndDate":"2018-05-29 15:20:17","availableDate":"2018-07-02 20:15:00","videoPlazaTags":["j4=2","j5=2","j5=3","j5=4","j5=5","app"],"airplayEnabled":true,"availableForWeb":true,"availableForMobile":true,"availableForDevice":null,"deeplink":null,"deeplinkUrl":"http:\/\/link.tvnow.de\/?f=1786&e=571634","packages":{"items":[{"id":4,"name":"Everest - AYCS"}],"total":0},"paymentItunesContainer":null,"itunesTitle":"Wallach \"Watabi\" \/ Schimmelstute \"Naira\"","itunesPosition":6,"fsk6":false,"fsk12":true,"fsk16":false,"fsk18":false,"fsk":12,"platform":["tvnow"],"qualityCheck":"1","kdgGenre1":"0000","kdgGenre2":"0000","umkbwGenre":"","similarCategories":[],"productionYear":0,"productionCountry":"","cast":"","imdbId":"","geoAllowedCountries":null,"clipfishId":0,"portability":{},"isLiveStream":false,"alternateBroadcastDateText":""}`),
			expected: "https://vodnowusodash.secure.footprint.net/proxy/clear/manifest/tvnow/571634-1-36465.ism/.mpd?filter=(type%3d%3d%22audio%22)%7c%7c(type%3d%3d%22video%22%26%26systemBitrate%3c1550000)"},
	}
	for _, c := range cc {
		got, err := extractIDFromAPIResponse(c.body)
		if got != c.expected {
			t.Errorf("got '%s', expected '%s'", got, c.expected)
		}
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}
	}
}

func TestURLParsing(t *testing.T) {
	if os.Getenv("TEST_ONLINE") == "" {
		t.Skip("online test skipped")
	}
	cc := []struct {
		url    string
		prefix string
	}{
		{
			url:    "https://www.tvnow.de/vox/die-pferdeprofis/wallach-watabi-schimmelstute-naira-2018-06-02-19-10-00/player",
			prefix: "https://vodnowusodash-a.akamaihd.net/proxy/clear/manifest/tvnow/571634-1-36465.ism/.mpd",
		},
	}

	for _, c := range cc {
		got, err := ExtractID(c.url)
		if !strings.HasPrefix(got, c.prefix) {
			t.Errorf("got '%s', expected prefix '%s' for '%s'", got, c.prefix, c.url)
		}
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}
	}
}
