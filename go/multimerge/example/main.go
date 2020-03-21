package main

import (
	"fmt"
	"go/go/multimerge"
	"math/rand"
	"sort"
)

type Item struct {
	ID   int
	Heat int
}

func (i Item) Equal(i2 multimerge.Noder) bool {
	return i.ID == i2.(Item).ID
}

func (i Item) LessThan(i2 multimerge.Noder) bool {
	return i.Heat < i2.(Item).Heat
}

type byItemIndex []Item

func (a byItemIndex) Len() int           { return len(a) }
func (a byItemIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byItemIndex) Less(i, j int) bool { return a[i].Heat > a[j].Heat }

func main() {
	fmt.Println("============ first example ==========")
	l1 := []Item{
		{ID: 1, Heat: 100},
		{ID: 2, Heat: 98},
		{ID: 3, Heat: 97},
		{ID: 4, Heat: 95},
		{ID: 5, Heat: 80},
	}

	l2 := []Item{
		{ID: 6, Heat: 100},
		{ID: 7, Heat: 97},
		{ID: 8, Heat: 93},
		{ID: 9, Heat: 92},
		{ID: 10, Heat: 81},
	}

	l3 := []Item{
		{ID: 11, Heat: 99},
		{ID: 12, Heat: 93},
		{ID: 13, Heat: 92},
		{ID: 14, Heat: 90},
		{ID: 15, Heat: 88},
	}

	l4 := []Item{
		{ID: 16, Heat: 94},
		{ID: 17, Heat: 92},
		{ID: 18, Heat: 92},
		{ID: 19, Heat: 87},
		{ID: 20, Heat: 85},
	}

	l5 := []Item{
		{ID: 21, Heat: 95},
		{ID: 22, Heat: 91},
		{ID: 23, Heat: 89},
		{ID: 24, Heat: 88},
		{ID: 25, Heat: 83},
	}

	lb := [][]Item{l1, l2, l3, l4, l5}

	ms := multimerge.NewSort(lb)
	nodes := ms.TopK(8)

	for _, node := range nodes {
		item := node.(Item)
		fmt.Println(item)
	}

	fmt.Println("============ second example ==========")

	L := 100000
	l := make([]Item, L)

	for i := 0; i < L; i++ {
		l[i] = Item{Heat: rand.Intn(L), ID: i}
	}

	ll := make([]Item, L)

	for i := 0; i < L; i++ {
		ll[i] = Item{Heat: rand.Intn(L), ID: i + L + 1}
	}

	lll := make([]Item, L)

	for i := 0; i < L; i++ {
		lll[i] = Item{Heat: rand.Intn(L), ID: i + 2*L + 1}
	}

	sort.Sort(byItemIndex(l))
	sort.Sort(byItemIndex(ll))
	sort.Sort(byItemIndex(lll))

	lb2 := [][]Item{l, ll, lll}
	mss := multimerge.NewSort(lb2)
	nodes = mss.TopK(L)
	fmt.Println(nodes)
}
