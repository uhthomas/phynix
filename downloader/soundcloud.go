package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"phynix/enum"
	"phynix/models"
	"pipeline"
	"strings"
)

var soundcloudPipeline = pipeline.New(10)

func Soundcloud(id string) (models.Media, error) {
	soundcloudPipeline.Populate()
	defer soundcloudPipeline.Free()

	var media models.Media
	if !models.DB.First(&media, "content_id = ?", id).RecordNotFound() {
		return media, nil
	}

	var data struct {
		Image       string `json:"artwork_url"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Length      int    `json:"duration"`
		User        struct {
			Avatar   string `json:"avatar_url"`
			Username string `json:"username"`
		} `json:"user"`
	}

	res, err := http.Get(fmt.Sprintf("https://api.soundcloud.com/tracks/%s?client_id=2d5a9eda4866bc2d96d385db39509df5",
		url.QueryEscape(id)))
	if err != nil {
		return models.Media{}, err
	}

	if res.StatusCode != http.StatusOK {
		return models.Media{}, errors.New("unexpected response code")
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return models.Media{}, err
	}

	image := data.Image
	if image == "" {
		image = data.User.Avatar
	}

	var (
		artist string
		title  string = data.Title
	)

	if s := strings.Split(title, " - "); len(s) > 1 {
		artist = s[0]
		title = strings.Join(s[1:], " - ")
	} else {
		artist = data.User.Username
	}

	m := models.Media{
		Type:      int(enum.MediaTypeSoundcloud),
		ContentID: id,
		Image:     image,
		Duration:  data.Length / 1000,
		Artist:    artist,
		Title:     title,
		Blurb:     data.Description,
	}

	if err := models.DB.Where("content_id = ? and type = ?", m.ContentID, m.Type).
		FirstOrCreate(&m).Error; err != nil {
		return models.Media{}, err
	}

	return m, nil
}
