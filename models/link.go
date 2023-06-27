package models

import (
	"errors"
)

type Link struct {
	Code string `json:"code"`
	Url  string `json:"url"`
}

type CreateLinkDto struct {
	Url string `json:"url"`
}

var ErrLinkNotFound = errors.New("link not found")

const LinkSymbols = "0123456789_AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
const LinkSymbolsLen = len(LinkSymbols)

func GenerateLinkCodeByUrl(url string) string {
	var hash [int8(10)]byte
	hashPosNumber := int8(0)
	for urlSymbolPos, urlSymbol := range url {
		currentHashValue := hash[hashPosNumber]
		newHashValueSymbolPos := int(int32(currentHashValue)+int32(urlSymbolPos)+urlSymbol) % LinkSymbolsLen
		newHashValue := LinkSymbols[newHashValueSymbolPos]

		hash[hashPosNumber] = newHashValue
		if hashPosNumber++; hashPosNumber >= 10 {
			hashPosNumber = 0
		}
	}

	if len(url) < 10 {
		for i := len(url); i < 10; i++ {
			hash[i] = LinkSymbols[0]
		}
	}

	return string(hash[:])
}
