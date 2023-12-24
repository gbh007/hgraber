package hgraber

type FirstHandleMultipleResult struct {
	TotalCount     int64
	LoadedCount    int64
	DuplicateCount int64
	ErrorCount     int64
	NotHandled     []string
}
