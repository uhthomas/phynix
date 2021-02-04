package validation

import "regexp"

func Displayname(str string) (valid bool) {
	valid = regexp.MustCompile(`^.{2,20}$`).MatchString(str)
	return
}

func Username(str string) (valid bool) {
	valid = regexp.MustCompile(`^[a-z0-9._-]{2,20}$`).MatchString(str)
	return
}

func Email(str string) (valid bool) {
	valid = len(str) >= 2 && len(str) <= 100 && regexp.MustCompile(`^.+@.+\..+`).MatchString(str)
	return
}

func Password(str string) (valid bool) {
	valid = regexp.MustCompile(`^.{2,72}$`).MatchString(str)
	return
}

func CommunitySlug(str string) (valid bool) {
	valid = regexp.MustCompile(`^[a-z0-9-]{2,20}$`).MatchString(str)
	return
}

func CommunityName(str string) (valid bool) {
	valid = regexp.MustCompile(`^.{2,30}$`).MatchString(str)
	return
}

func CommunityDescription(str string) (valid bool) {
	valid = regexp.MustCompile(`^.{0,500}$`).MatchString(str)
	return
}

func CommunityWelcomeMessage(str string) (valid bool) {
	valid = regexp.MustCompile(`^.{0,300}$`).MatchString(str)
	return
}

func BanReason(str string) (valid bool) {
	valid = regexp.MustCompile(`^.{0,500}$`).MatchString(str)
	return
}

func PlaylistName(str string) (valid bool) {
	valid = regexp.MustCompile(`^.{1,30}$`).MatchString(str)
	return
}

func PlaylistItemArtist(str string) (valid bool) {
	valid = regexp.MustCompile(`^.{1,100}$`).MatchString(str)
	return
}

func PlaylistItemTitle(str string) (valid bool) {
	valid = regexp.MustCompile(`^.{1,100}$`).MatchString(str)
	return
}
