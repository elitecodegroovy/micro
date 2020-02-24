package mergence

import (
	"testing"
)

func TestMinInSlice(t *testing.T) {
	tests := []struct {
		a []int64
		n int
	}{
		{[]int64{12, 22, 23, 29, 10, 2, 999}, 5},
		{[]int64{111, 22, 231, 291, 101, 210, 999}, 1},
		{[]int64{1, 22, 23, 29, 100, 202, 999}, 0},
		{[]int64{1, 22, 23, 29, 100, 202, 0}, 6},
		{[]int64{11111111111}, 0},
		{[]int64{11111111111, 1}, 1},
	}
	for _, test := range tests {
		r := minInSlice(test.a)
		if r != test.n {
			t.Errorf("expected %v, got %v", test.n, r)
		}
	}
}

func TestSortItems(t *testing.T) {
	tests := []struct {
		a   []lineItem
		min int64
	}{
		{[]lineItem{{index: 0, value: 12}}, 12},
		{[]lineItem{{index: 0, value: 12}, {index: 1, value: 10}}, 10},
		{[]lineItem{{index: 0, value: 12}, {index: 1, value: 1}, {index: 0, value: 2}, {index: 1, value: 10}}, 1},
		{[]lineItem{{index: 0, value: 12}, {index: 1, value: 1}, {index: 0, value: 2}, {index: 1, value: 10},
			{index: 0, value: 222}, {index: 1, value: 122}, {index: 0, value: 22288888888}, {index: 1, value: 10888888}},
			1},
		{[]lineItem{{index: 0, value: 121}, {index: 1, value: 11}, {index: 0, value: 1112}, {index: 1, value: 11111110},
			{index: 0, value: 2211111112}, {index: 1, value: 11111111122}, {index: 0, value: 2228822888888}, {index: 1, value: 10222222888888}},
			11},
	}
	for _, test := range tests {
		sortItems(test.a)
		if test.a[0].value != test.min {
			t.Errorf("expected %v, got %v", test.min, test.a[0].value)
		}
	}

}
