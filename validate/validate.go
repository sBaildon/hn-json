package validate

import (
	"net/url"
)

func Title(title string) bool {
	return (len(title) > 0) && (len(title) <= 256)
}

func Author(author string) bool {
	return (len(author) > 0) && (len(author) <= 256)
}

func URI(uri string) bool {
	if _, err := url.Parse(uri); err != nil {
		return false
	}

	return true
}

func Points(points int) bool {
	return (points >= 0)
}

func Comments(comments int) bool {
	return (comments >= 0)
}

func Rank(rank int) bool {
	return (rank >= 0)
}
