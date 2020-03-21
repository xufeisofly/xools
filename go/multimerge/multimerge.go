package multimerge

import (
	"reflect"
)

type MSorter struct {
	ListBundle [][]Noder
	ListPtr    []int
	MaxHeap    Heap
}

type Noder interface {
	LessThan(Noder) bool
	Equal(Noder) bool
}

// // ListNode is a Noder with ListBundle's list index
// type ListNode struct {
// 	ListIndex int
// 	Node      Noder
// }

// func (ln ListNode) Equal(ln2 Noder) bool {
// 	return ln.Node.Equal(ln2.(ListNode).Node)
// }

// func (ln ListNode) LessThan(ln2 Noder) bool {
// 	return ln.Node.LessThan(ln2.(ListNode).Node)
// }

type List interface{}

func NewSort(lists interface{}) MSorter {
	s := reflect.ValueOf(lists)
	if s.Kind() != reflect.Slice {
		panic("NewSort is given a non-slice type")
	}
	lb := make([][]Noder, s.Len())
	lptr := make([]int, s.Len())

	for i := 0; i < s.Len(); i++ {
		listValue := s.Index(i)
		if listValue.Kind() != reflect.Slice {
			panic("NewSort slice element is a non-slice type")
		}

		l := make([]Noder, listValue.Len())

		for j := 0; j < listValue.Len(); j++ {
			l[j] = listValue.Index(j).Interface().(Noder)
		}
		lb[i] = l
		lptr[i] = 0
	}

	ret := MSorter{
		ListBundle: lb,
		ListPtr:    lptr,
	}
	ret.initMaxHeap()
	return ret
}

func (ms *MSorter) initMaxHeap() {
	cNodes := ms.candidateNodes()

	h := Heap(cNodes)
	maxH := h.MakeMaxHeap()

	ms.MaxHeap = maxH
}

func (ms MSorter) candidateNodes() []Noder {
	ret := make([]Noder, len(ms.ListBundle))
	for i, l := range ms.ListBundle {
		ret[i] = ms.nodeByNodeIndex(0, l)
	}
	return ret
}

func (ms MSorter) nodeByNodeIndex(idx int, l []Noder) Noder {
	if idx > len(l)-1 {
		return nil
	}

	ret := l[idx]
	return ret
}

func (ms MSorter) ShiftMaxNode() Noder {
	maxNode := ms.MaxHeap.RootNode()
	ms.MaxHeap = ms.MaxHeap.deleteRootNode()
	node := ms.nextNode(maxNode)

	if node == nil {
		panic("one of the list overflow, busy, I will enhance the code sometime")
	}
	ms.MaxHeap = ms.MaxHeap.PushNode(node)
	return maxNode
}

func (ms MSorter) nextNode(lastNode Noder) Noder {
	var listIndex = -1
	for i, l := range ms.ListBundle {
		if l[ms.ListPtr[i]].Equal(lastNode) {
			listIndex = i
			break
		}
	}

	if listIndex == -1 {
		panic("max node home list no found")
	}

	l := ms.ListBundle[listIndex]
	nodeIndex := ms.ListPtr[listIndex]

	// move list node ptr
	ms.ListPtr[listIndex] = nodeIndex + 1
	node := ms.nodeByNodeIndex(nodeIndex+1, l)
	return node
}

func (ms MSorter) TopK(k int) []Noder {
	ret := make([]Noder, k)
	for i := 0; i < k; i++ {
		ret[i] = ms.ShiftMaxNode()
	}

	return ret
}
