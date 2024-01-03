package hgraber

import "testing"

func TestGeneratePagination(t *testing.T) {
	tests := []struct {
		Name        string
		PageCount   int
		CurrentPage int
		Output      []int
	}{
		{
			Name:        "Full-1",
			PageCount:   12,
			CurrentPage: 1,
			Output:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		{
			Name:        "Full-2",
			PageCount:   5,
			CurrentPage: 1,
			Output:      []int{1, 2, 3, 4, 5},
		},
		{
			Name:        "Big-1",
			PageCount:   100,
			CurrentPage: 50,
			Output:      []int{1, 2, 3, -1, 48, 49, 50, 51, 52, -1, 98, 99, 100},
		},
		{
			Name:        "Big-Left-1",
			PageCount:   100,
			CurrentPage: 4,
			Output:      []int{1, 2, 3, 4, 5, 6, 7, 8, -1, 98, 99, 100},
		},
		{
			Name:        "Big-Left-2",
			PageCount:   100,
			CurrentPage: 8,
			Output:      []int{1, 2, 3, -1, 6, 7, 8, 9, 10, -1, 98, 99, 100},
		},
		{
			Name:        "Big-Left-3",
			PageCount:   100,
			CurrentPage: 7,
			Output:      []int{1, 2, 3, 4, 5, 6, 7, 8, -1, 98, 99, 100},
		},
		{
			Name:        "Big-Right-1",
			PageCount:   100,
			CurrentPage: 96,
			Output:      []int{1, 2, 3, -1, 93, 94, 95, 96, 97, 98, 99, 100},
		},
		{
			Name:        "Big-Right-2",
			PageCount:   100,
			CurrentPage: 93,
			Output:      []int{1, 2, 3, -1, 91, 92, 93, 94, 95, -1, 98, 99, 100},
		},
		{
			Name:        "Big-Right-3",
			PageCount:   100,
			CurrentPage: 94,
			Output:      []int{1, 2, 3, -1, 93, 94, 95, 96, 97, 98, 99, 100},
		},
	}

	for index := range tests {
		test := tests[index]

		t.Run(test.Name, func(t *testing.T) {
			result := generatePagination(test.CurrentPage, test.PageCount)
			if !compareSlice(result, test.Output) {
				t.Log(result)
				t.Fail()
			}
		})
	}
}

func compareSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for index := range a {
		if a[index] != b[index] {
			return false
		}
	}

	return true
}
