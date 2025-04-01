package gokick_test

import (
	"fmt"
	"testing"

	"github.com/scorfly/gokick"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLivestreamSortError(t *testing.T) {
	testCases := map[string]string{
		"empty":         "",
		"not supported": "not supported",
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := gokick.NewLivestreamSort(value)
			assert.EqualError(t, err, fmt.Sprintf("unknown livestream sort: %s", value))
		})
	}
}

func TestNewLivestreamSortSuccess(t *testing.T) {
	testCases := map[string]gokick.LivestreamSort{
		"viewer_count": gokick.LivestreamSortViewerCount,
		"started_at":   gokick.LivestreamSortStartedAt,
	}

	for name, value := range testCases {
		t.Run(name, func(t *testing.T) {
			LivestreamSort, err := gokick.NewLivestreamSort(value.String())
			require.NoError(t, err)
			assert.Equal(t, LivestreamSort, value)
		})
	}
}
