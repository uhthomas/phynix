package searcher

import (
	"net/http"
	"phynix/enum"
	"strings"
	"time"

	"github.com/6f7262/pipe"

	"code.google.com/p/google-api-go-client/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

var (
	youtubeService *youtube.Service
	yp             = pipe.New(10)
)

func init() {
	client := &http.Client{
		Transport: &transport.APIKey{Key: "AIzaSyBAdDIgUc_loht-bJyBtaRcD8aDeupAaeE"},
	}

	var err error
	youtubeService, err = youtube.New(client)
	if err != nil {
		panic(err)
	}
}

func Youtube(query string) ([]Result, error) {
	defer yp.One()()

	call := youtubeService.Search.List("id,snippet").
		Q(query).
		Type("video").
		MaxResults(50)

	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(res.Items) <= 0 {
		return nil, nil
	}

	videoCall := youtubeService.Videos.List("snippet,contentDetails").Id(func() string {
		ids := []string{}
		for _, item := range res.Items {
			ids = append(ids, item.Id.VideoId)
		}
		return strings.Join(ids, ",")
	}())

	videoRes, err := videoCall.Do()

	if len(videoRes.Items) <= 0 {
		return nil, nil
	}

	results := []Result{}

	for _, item := range videoRes.Items {
		var (
			artist string
			title  string = item.Snippet.Title
		)

		if s := strings.Split(title, " - "); len(s) > 1 {
			artist = s[0]
			title = strings.Join(s[1:], " - ")
		} else {
			artist = item.Snippet.ChannelTitle
		}

		dur, err := time.ParseDuration(strings.ToLower(item.ContentDetails.Duration[2:]))
		if err != nil {
			continue
		}

		thumbnails := item.Snippet.Thumbnails
		var thumbnail *youtube.Thumbnail
		switch {
		case thumbnails.Maxres != nil:
			thumbnail = thumbnails.Maxres
		case thumbnails.High != nil:
			thumbnail = thumbnails.High
		case thumbnails.Medium != nil:
			thumbnail = thumbnails.Medium
		case thumbnails.Standard != nil:
			thumbnail = thumbnails.Standard
		case thumbnails.Default != nil:
			thumbnail = thumbnails.Default
		default:
			continue
		}

		results = append(results, Result{
			Type:      int(enum.MediaTypeYoutube),
			ContentID: item.Id,
			Image:     thumbnail.Url,
			Duration:  int(dur.Seconds()),
			Artist:    artist,
			Title:     title,
		})
	}

	return results, nil
}
