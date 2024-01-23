package rendering

import (
	"app/internal/domain/hgraber"
)

type FirstHandleMultipleResult struct {
	TotalCount     int64              `json:"total_count"`
	LoadedCount    int64              `json:"loaded_count"`
	DuplicateCount int64              `json:"duplicate_count"`
	ErrorCount     int64              `json:"error_count"`
	NotHandled     []string           `json:"not_handled,omitempty"`
	Details        []BookHandleResult `json:"details"`
}

type BookHandleResult struct {
	URL         string `json:"url"`
	IsDuplicate bool   `json:"is_duplicate"`
	IsHandled   bool   `json:"is_handled"`
	ErrorReason string `json:"error_reason"`
}

func HandleMultipleResultFromDomain(raw *hgraber.FirstHandleMultipleResult) FirstHandleMultipleResult {
	out := FirstHandleMultipleResult{
		TotalCount:     raw.TotalCount,
		LoadedCount:    raw.LoadedCount,
		DuplicateCount: raw.DuplicateCount,
		ErrorCount:     raw.ErrorCount,
		NotHandled:     make([]string, len(raw.NotHandled)),
		Details:        make([]BookHandleResult, len(raw.Details)),
	}

	copy(out.NotHandled, raw.NotHandled)
	convertSlice(out.Details, raw.Details, func(data hgraber.BookHandleResult) BookHandleResult {
		return BookHandleResult{
			URL:         data.URL,
			IsDuplicate: data.IsDuplicate,
			IsHandled:   data.IsHandled,
			ErrorReason: data.ErrorReason,
		}
	})

	return out
}
