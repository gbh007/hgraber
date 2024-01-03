package rendering

import (
	"app/internal/domain/hgraber"
)

type FirstHandleMultipleResult struct {
	TotalCount     int64    `json:"total_count"`
	LoadedCount    int64    `json:"loaded_count"`
	DuplicateCount int64    `json:"duplicate_count"`
	ErrorCount     int64    `json:"error_count"`
	NotHandled     []string `json:"not_handled,omitempty"`
}

func HandleMultipleResultFromDomain(raw *hgraber.FirstHandleMultipleResult) FirstHandleMultipleResult {
	out := FirstHandleMultipleResult{
		TotalCount:     raw.TotalCount,
		LoadedCount:    raw.LoadedCount,
		DuplicateCount: raw.DuplicateCount,
		ErrorCount:     raw.ErrorCount,
		NotHandled:     make([]string, len(raw.NotHandled)),
	}

	copy(out.NotHandled, raw.NotHandled)

	return out
}
