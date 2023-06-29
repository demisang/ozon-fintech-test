package models

import (
	"errors"
	"regexp"
)

const shortLinkInvalidPattern = `[^a-zA-Z\d_]+`

var CompiledTemplate = regexp.MustCompile(shortLinkInvalidPattern)

type Link struct {
	Code string `json:"code"`
	URL  string `json:"url"`
}

type CreateLinkDTO struct {
	URL string `json:"url"`
}

var ErrLinkNotFound = errors.New("link not found")

const (
	linkSymbols    = "0123456789_AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	linkSymbolsLen = len(linkSymbols)
)

func GenerateLinkCodeByURL(url string) string {
	var hash [10]byte

	hashPosNumber := int8(0)
	for urlSymbolPos, urlSymbol := range url {
		currentHashValue := hash[hashPosNumber]
		newHashValueSymbolPos := int(int32(currentHashValue)+int32(urlSymbolPos)+urlSymbol) % linkSymbolsLen
		newHashValue := linkSymbols[newHashValueSymbolPos]

		hash[hashPosNumber] = newHashValue

		if hashPosNumber++; hashPosNumber >= 10 {
			hashPosNumber = 0
		}
	}

	if len(url) < 10 {
		for i := len(url); i < 10; i++ {
			hash[i] = linkSymbols[len(url)]
		}
	}

	return string(hash[:])
}
