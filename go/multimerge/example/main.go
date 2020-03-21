package main

import (
	"fmt"
	"go/go/multimerge"
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

func main() {
	l := []Item{
		{Value: 9},
		{Value: 7},
		{Value: 13},
		{Value: 15},
		{Value: 11},
	}

	h := multimerge.NewHeap(l)
	fmt.Println(h)
	hmax := h.MakeMaxHeap()
	fmt.Println(hmax)
}
