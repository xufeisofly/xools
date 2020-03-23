package multimerge

import (
	"reflect"
)

type Heap []Noder

func (h Heap) RootNode() Noder {
	return h[0]
}

func (h Heap) rootNodeIndex() int {
	return 0
}

func (h Heap) lastParentNodeIndex() int {
	length := len(h)
	if isEven(len(h)) {
		return length/2 - 1
	}
	return (length-1)/2 - 1
}

func (h Heap) lastNodeIndex() int {
	return len(h) - 1
}

func (h Heap) shouldSwapWithChild(i int) bool {
	pNode := h[i]
	biggerChildNode, _ := h.biggerChildNodeWithIndex(i)
	return pNode.LessThan(biggerChildNode)
}

func (h Heap) biggerChildNodeWithIndex(pIdx int) (Noder, int) {
	node1, idx1 := h.leftChildNodeWithIndex(pIdx)
	node2, idx2 := h.rightChildNodeWithIndex(pIdx)

	if node1 == nil {
		return node2, idx2
	}
	if node2 == nil {
		return node1, idx1
	}

	if node1.LessThan(node2) {
		return node2, idx2
	}
	return node1, idx1
}

func (h Heap) leftChildNodeWithIndex(i int) (Noder, int) {
	index := 2*i + 1
	return h[index], index
}

func (h Heap) rightChildNodeWithIndex(i int) (Noder, int) {
	if isEven(len(h)) && i == h.lastParentNodeIndex() {
		return nil, 0
	}

	index := 2*i + 2
	return h[index], index
}

func NewHeap(l List) Heap {
	s := reflect.ValueOf(l)
	if s.Kind() != reflect.Slice {
		panic("NewHeap is given a non-slice type")
	}
	ret := make([]Noder, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface().(Noder)
	}
	return ret
}

func (h Heap) flowUpByIndex(index int) {
	pNode, pIdx := h.parentNodeWithIndex(index)
	cNode, cIdx := h.biggerChildNodeWithIndex(pIdx)

	if pNode.LessThan(cNode) {
		h.swapByIndex(pIdx, cIdx)
		h.flowDownByIndex(cIdx)
		h.flowUpByIndex(pIdx)
	}
}

func (h Heap) flowDownByIndex(index int) {
	node := h[index]
	if h.hasChild(index) {
		cNode, cIdx := h.biggerChildNodeWithIndex(index)
		if node.LessThan(cNode) {
			h.swapByIndex(index, cIdx)
			h.flowDownByIndex(cIdx)
		}
	}
}

func (h Heap) hasChild(i int) bool {
	return i <= h.lastParentNodeIndex()
}

// TODO 这个 for 循环太傻逼了
func (h Heap) index(n Noder) int {
	for i, node := range h {
		if node.Equal(n) {
			return i
		}
	}
	panic("cannot find index")
}

func (h Heap) swapByIndex(index1, index2 int) {
	tmp := h[index1]
	h[index1] = h[index2]
	h[index2] = tmp
}

func (h Heap) parentNodeWithIndex(idx int) (Noder, int) {
	node := h[idx]
	if idx == 0 {
		return node, 0
	}
	if isEven(idx) {
		return h[idx/2-1], idx/2 - 1
	}
	return h[(idx+1)/2-1], (idx+1)/2 - 1
}

func (h Heap) MakeMaxHeap() Heap {
	for i := h.lastParentNodeIndex(); i >= 0; i-- {
		if h.shouldSwapWithChild(i) {
			_, upIdx := h.biggerChildNodeWithIndex(i)
			h.flowUpByIndex(upIdx)
		}
	}
	return h
}

func (h Heap) LastNode() Noder {
	return h[len(h)-1]
}

func (h Heap) deleteRootNode() Heap {
	h.swapByIndex(h.rootNodeIndex(), h.lastNodeIndex())
	h = h[0 : len(h)-1]
	h.flowDownByIndex(h.rootNodeIndex())
	return h
}

func (h Heap) PushNode(node Noder) Heap {
	h = append(h, node)
	h.MakeMaxHeap()
	return h
}
