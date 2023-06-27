package models_test

import (
	"testing"

	"github.com/demisang/ozon-fintech-test/internal/models"
	"github.com/stretchr/testify/require"
)

func TestGenerateLinkCodeByURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		hash string
	}{
		{
			name: "correct_hash_1",
			url:  "http://google.com/",
			hash: "4sAtrvfFwX",
		},
		{
			name: "empty_url",
			hash: "0000000000",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.hash, models.GenerateLinkCodeByURL(test.url))
		})
	}
}
