package rendering

import "fmt"

func convertSlice[From any, To any](to []To, from []From, conv func(From) To) {
	for i, v := range from {
		if i >= len(to) {
			return
		}

		to[i] = conv(v)
	}
}

func convertSliceWithAddr[From any, To any](addr string, to []To, from []From, conv func(string, From) To) {
	for i, v := range from {
		if i >= len(to) {
			return
		}

		to[i] = conv(addr, v)
	}
}

func fileURL(addr string, bookID int, pageNumber int, ext string) string {
	return fmt.Sprintf("%s/file/%d/%d.%s", addr, bookID, pageNumber, ext)
}
