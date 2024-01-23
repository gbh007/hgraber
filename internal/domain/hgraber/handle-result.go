package hgraber

type FirstHandleMultipleResult struct {
	TotalCount     int64
	LoadedCount    int64
	DuplicateCount int64
	ErrorCount     int64
	NotHandled     []string
	Details        []BookHandleResult
}

type BookHandleResult struct {
	URL         string
	IsDuplicate bool
	IsHandled   bool
	ErrorReason string
}
