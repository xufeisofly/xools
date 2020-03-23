package multimerge

import (
	"reflect"
	"sync"
)

type MSorter struct {
	ListIn     []List
	ListBundle [][]Noder
	ListPtr    []int
	MaxHeap    Heap
}

type Noder interface {
	LessThan(Noder) bool
	Equal(Noder) bool
}

type List interface{}

func NewSort(lists interface{}) MSorter {
	s := reflect.ValueOf(lists)
	if s.Kind() != reflect.Slice {
		panic("NewSort is given a non-slice type")
	}
	// lb := make([][]Noder, s.Len())
	lptr := make([]int, s.Len())
	lin := make([]List, s.Len())

	wg := sync.WaitGroup{}
	for i := 0; i < s.Len(); i++ {
		wg.Add(1)
		go func(i int, s reflect.Value) {
			listValue := s.Index(i)
			if listValue.Kind() != reflect.Slice {
				panic("NewSort slice element is a non-slice type")
			}

			lin[i] = s.Index(i).Interface().(List)
			lptr[i] = 0
			wg.Done()
		}(i, s)
	}

	wg.Wait()
	ret := MSorter{
		ListIn: lin,
		// ListBundle: lb,
		ListPtr: lptr,
	}
	// ret.initMaxHeap()
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
		panic("Only K < L is supported, K is TopK arg, L is the minimum length of list")
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
		panic("Only K < L is supported, K is TopK arg, L is the minimum length of list")
	}

	l := ms.ListBundle[listIndex]
	nodeIndex := ms.ListPtr[listIndex]

	// move list node ptr
	ms.ListPtr[listIndex] = nodeIndex + 1
	node := ms.nodeByNodeIndex(nodeIndex+1, l)
	return node
}

func (ms MSorter) TopK(k int) []Noder {
	// 我只能说这里是最最最拖慢速度的
	lb := make([][]Noder, len(ms.ListIn))
	wg := sync.WaitGroup{}
	for j, list := range ms.ListIn {
		wg.Add(1)
		go func(j int, list List) {
			lv := reflect.ValueOf(list)
			if lv.Kind() != reflect.Slice {
				panic("NewSort slice element is a non-slice type")
			}

			var listLength int
			// 如果 L > N*K，那么只需要在 N*K 的 长度里取值就行了
			if lv.Len() > k*len(ms.ListIn) {
				listLength = k * len(ms.ListIn)
			} else {
				listLength = lv.Len()
			}
			l := make([]Noder, listLength)

			for ii := 0; ii < listLength; ii++ {
				l[ii] = lv.Index(ii).Interface().(Noder)
			}

			lb[j] = l
			wg.Done()
		}(j, list)
	}
	wg.Wait()

	ms.ListBundle = lb
	ms.initMaxHeap()

	ret := make([]Noder, k)
	for i := 0; i < k; i++ {
		ret[i] = ms.ShiftMaxNode()
	}

	return ret
}
