package enum

import "net/http"

type (
	ResponseCode  int
	CommunityRole int
	GlobalRole    int
	MediaType     int
	VoteType      int
)

const (
	ResponseCodeError ResponseCode = iota + 1
	ResponseCodeOk
	ResponseCodeBadRequest
	ResponseCodeForbidden
)

const (
	CommunityRoleVip CommunityRole = iota + 1
	CommunityRoleModerator
	CommunityRoleManager
	CommunityRoleBot
	CommunityRoleCoHost
	CommunityRoleHost
)

const (
	GlobalRoleUser      GlobalRole = 0
	GlobalRoleModerator            = iota + 100
	GlobalRoleAdmin
	GlobalRoleServer
)

const (
	MediaTypeYoutube MediaType = iota + 1
	MediaTypeSoundcloud
)

const (
	VoteTypeWoot VoteType = iota
	VoteTypeSave
	VoteTypeMeh
)

func (r ResponseCode) String() string {
	return map[ResponseCode]string{
		ResponseCodeError:      "error",
		ResponseCodeOk:         "ok",
		ResponseCodeBadRequest: "badRequest",
		ResponseCodeForbidden:  "forbidden",
	}[r]
}

func (r ResponseCode) Http() int {
	return map[ResponseCode]int{
		ResponseCodeError:      http.StatusInternalServerError,
		ResponseCodeOk:         http.StatusOK,
		ResponseCodeBadRequest: http.StatusBadRequest,
		ResponseCodeForbidden:  http.StatusForbidden,
	}[r]
}
