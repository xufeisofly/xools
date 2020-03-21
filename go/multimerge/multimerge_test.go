package multimerge_test

import (
	"go/go/multimerge"
	"math/rand"
	"os"
	"sort"
	"testing"
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

func TestTopK(t *testing.T) {
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

	expected := []int{100, 100, 99, 98, 97}

	ms := multimerge.NewSort([][]Item{l1, l2, l3, l4, l5})
	var values = make([]int, 5)
	rets := ms.TopK(5)

	for i, ret := range rets {
		values[i] = ret.(Item).Heat
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

type byItemIndex []Item

func (a byItemIndex) Len() int           { return len(a) }
func (a byItemIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byItemIndex) Less(i, j int) bool { return a[i].Heat > a[j].Heat }

const L = 100000
const numList = 50
const K = 100

func getLists(L int) [][]Item {
	var listBundle = make([][]Item, numList)

	for i := 0; i < numList; i++ {
		l := make([]Item, L)

		for j := 0; j < L; j++ {
			l[j] = Item{Heat: rand.Intn(L), ID: j + i*(L+1)}
		}

		sort.Sort(byItemIndex(l))
		listBundle[i] = l
	}

	return listBundle
}

// BenchmarkTopK-4      	       1	1522500513 ns/op
func BenchmarkTopK(b *testing.B) {
	lb := getLists(L)

	for i := 0; i < b.N; i++ {
		ms := multimerge.NewSort(lb)
		_ = ms.TopK(K)
	}
}

// BenchmarkOldTopK-4   	       1	1981787002 ns/op
func BenchmarkOldTopK(b *testing.B) {
	lb := getLists(L)

	for i := 0; i < b.N; i++ {
		var lall = make([]Item, len(lb)*L)
		for _, l := range lb {
			lall = append(lall, l...)
		}

		sort.Sort(byItemIndex(lall))
		_ = lall[:K]
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
