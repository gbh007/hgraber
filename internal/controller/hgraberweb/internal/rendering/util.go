package rendering

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
