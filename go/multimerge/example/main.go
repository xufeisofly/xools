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
	l1 := []Item{
		{Value: 18},
		{Value: 17},
		{Value: 14},
		{Value: 12},
		{Value: 11},
	}

	l2 := []Item{
		{Value: 19},
		{Value: 16},
		{Value: 15},
		{Value: 13},
		{Value: 10},
	}

	l3 := []Item{
		{Value: 17},
		{Value: 15},
		{Value: 12},
		{Value: 11},
		{Value: 9},
	}

	lb := [][]Item{l1, l2, l3}

	ms := multimerge.NewSort(lb)
	nodes := ms.TopK(5)
	fmt.Println(nodes)
}
