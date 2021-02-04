package downloader

import (
	"errors"
	"net/http"
	"phynix/enum"
	"phynix/models"
	"pipeline"
	"strings"
	"time"

	"code.google.com/p/google-api-go-client/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

var (
	youtubeService  *youtube.Service
	youtubePipeline = pipeline.New(10)
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

func Youtube(id string) (models.Media, error) {
	youtubePipeline.Populate()
	defer youtubePipeline.Free()

	call := youtubeService.Videos.List("snippet,contentDetails").
		Id(id)

	res, err := call.Do()
	if err != nil {
		return models.Media{}, err
	}

	if len(res.Items) <= 0 {
		return models.Media{}, errors.New("returned no videos")
	}

	item := res.Items[0]

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
		return models.Media{}, err
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
		return models.Media{}, errors.New("video does not contain a thumbnail")
	}

	m := models.Media{
		Type:      int(enum.MediaTypeYoutube),
		ContentID: id,
		Image:     thumbnail.Url,
		Duration:  int(dur.Seconds()),
		Artist:    artist,
		Title:     title,
		Blurb:     item.Snippet.Description,
	}

	if err := models.DB.Where("content_id = ? and type = ?", m.ContentID, m.Type).
		FirstOrCreate(&m).Error; err != nil {
		return models.Media{}, err
	}

	return m, nil
}
