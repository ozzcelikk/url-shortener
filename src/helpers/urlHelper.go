package helpers

import (
	"math/rand"
	"net/url"
)

func isEmpty(str string) bool {
	if len(str) == 0 {
		return true
	}

	return false
}

func IsValidUrl(str string) bool {
	isStrEmpty := isEmpty(str)

	if isStrEmpty == true {
		return false
	}

	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

var chars = []rune("0123456789abcdefghijklmnoprstuvwxyzABCDEFGHIJKLMNOPRSTUVWXYZ")

func GetRandomShortUrl() string {
	shortUrl := ""

	for len(shortUrl) < 7 {
		rnd := rand.Intn(56)

		println(rnd)

		shortUrl += string(chars[rnd])
	}

	return shortUrl
}
