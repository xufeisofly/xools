package multimerge_test

import (
	"go/go/multimerge"
	"testing"
)

type Item struct {
	Value int
}

func (i Item) Equal(i2 multimerge.Noder) bool {
	return i.Value == i2.(Item).Value
}

func (i Item) LessThan(i2 multimerge.Noder) bool {
	return i.Value < i2.(Item).Value
}

func TestTopK(t *testing.T) {
	l1 := []Item{
		{Value: 30},
		{Value: 25},
		{Value: 20},
		{Value: 18},
		{Value: 17},
		{Value: 14},
		{Value: 12},
		{Value: 11},
	}

	l2 := []Item{
		{Value: 33},
		{Value: 24},
		{Value: 22},
		{Value: 19},
		{Value: 16},
		{Value: 15},
		{Value: 13},
		{Value: 10},
	}

	l3 := []Item{
		{Value: 50},
		{Value: 36},
		{Value: 31},
		{Value: 17},
		{Value: 15},
		{Value: 12},
		{Value: 11},
		{Value: 9},
	}

	expected := []int{50, 36, 33, 31, 30}

	ms := multimerge.NewSort([][]Item{l1, l2, l3})
	var values = make([]int, 5)
	rets := ms.TopK(5)

	for i, ret := range rets {
		values[i] = ret.(Item).Value
	}

	for i, value := range values {
		if value != expected[i] {
			t.Fatalf("TopK is not right, expected: %v, actual: %v", expected, values)
		}
	}

	if len(values) != len(expected) {
		t.Fatalf("TopK is not right, expected: %v, actual: %v", expected, values)
	}
}
