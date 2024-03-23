package rendering

import (
	"fmt"
	"strconv"
)

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

func PrettySize(raw int64) string {
	if raw < 1 {
		return "? б"
	}

	var div, mod int64

	const divider = 1024

	div = raw
	step := 0

	for div/divider > 0 {
		step++
		mod = div % divider
		div = div / divider
	}

	return strconv.FormatInt(div, 10) + "." + strconv.FormatInt(mod*10/1024, 10) + " " + sizeUnitFromStep(step)
}

func sizeUnitFromStep(step int) string {
	switch step {
	case 0:
		return "б"
	case 1:
		return "Кб"
	case 2:
		return "Мб"
	case 3:
		return "Гб"
	case 4:
		return "Тб"
	default:
		return "??"
	}
}
