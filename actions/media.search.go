package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/searcher"
	"raiki"
	"strings"
)

func MediaSearch(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		Query string         `json:"query"`
		Type  enum.MediaType `json:"type"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	data.Query = strings.TrimSpace(data.Query)

	var results []searcher.Result
	var err error

	switch data.Type {
	case enum.MediaTypeYoutube:
		results, err = searcher.Youtube(data.Query)
	// case enum.MediaTypeSoundcloud:
	// 	media, err = searcher.Soundcloud(data.ContentID)
	default:
		return int(enum.ResponseCodeBadRequest), nil, ErrMediaType
	}

	if err != nil {
		return int(enum.ResponseCodeError), nil, ErrMediaSearch
	}

	return int(enum.ResponseCodeOk), results, nil
}
