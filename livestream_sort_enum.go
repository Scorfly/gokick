package gokick

import "fmt"

type LivestreamSort int

const (
	LivestreamSortViewerCount LivestreamSort = iota // viewer_count
	LivestreamSortStartedAt                         // started_at
)

func NewLivestreamSort(sort string) (LivestreamSort, error) {
	switch sort {
	case "viewer_count":
		return LivestreamSortViewerCount, nil
	case "started_at":
		return LivestreamSortStartedAt, nil
	default:
		return 0, fmt.Errorf("unknown livestream sort: %s", sort)
	}
}

func (s LivestreamSort) String() string {
	switch s {
	case LivestreamSortViewerCount:
		return "viewer_count"
	case LivestreamSortStartedAt:
		return "started_at"
	default:
		return "unknown"
	}
}
