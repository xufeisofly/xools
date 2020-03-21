package multimerge

import (
	"reflect"
)

type Heap []Noder

func (h Heap) RootNode() Noder {
	return h[0]
}

func (h Heap) LastParentNode() Noder {
	length := len(h)
	if isEven(len(h)) {
		return h[length/2-1]
	}
	return h[(length-1)/2-1]
}

func (h Heap) shouldSwapWithChild(pNode Noder) bool {
	biggerChildNode := h.biggerChildNode(pNode)
	return pNode.LessThan(biggerChildNode)
}

func (h Heap) biggerChildNode(pNode Noder) Noder {
	node1 := h.leftChildNode(pNode)
	node2 := h.rightChildNode(pNode)
	if node1 == nil {
		return node2
	}
	if node2 == nil {
		return node1
	}

	if node1.LessThan(node2) {
		return node2
	}
	return node1
}

func (h Heap) leftChildNode(pNode Noder) Noder {
	return h[2*h.index(pNode)+1]
}

func (h Heap) rightChildNode(pNode Noder) Noder {
	idx := h.index(pNode)
	if isEven(len(h)) && idx == h.index(h.LastParentNode()) {
		return nil
	}
	return h[2*idx+2]
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

func (h Heap) flowUp(node Noder) {
	pNode := h.parentNode(node)
	cNode := h.biggerChildNode(pNode)
	pIdx := h.index(pNode)
	cIdx := h.index(cNode)
	if pNode.LessThan(cNode) {
		h.swap(pNode, cNode)
		h.flowDown(h[cIdx])
		h.flowUp(h[pIdx])
	}
}

func (h Heap) flowDown(node Noder) {
	if h.hasChild(node) {
		cNode := h.biggerChildNode(node)
		cIdx := h.index(cNode)
		if node.LessThan(cNode) {
			h.swap(node, cNode)
			h.flowDown(h[cIdx])
		}
	}
}

func (h Heap) hasChild(node Noder) bool {
	return h.index(node) <= h.index(h.LastParentNode())
}

func (h Heap) index(n Noder) int {
	for i, node := range h {
		if node.Equal(n) {
			return i
		}
	}
	panic("cannot find index")
}

func (h Heap) swap(node1, node2 Noder) {
	idx1 := h.index(node1)
	idx2 := h.index(node2)
	h[idx1] = node2
	h[idx2] = node1
}

func (h Heap) parentNode(node Noder) Noder {
	idx := h.index(node)
	if idx == 0 {
		return node
	}
	if isEven(idx) {
		return h[idx/2-1]
	}
	return h[(idx+1)/2-1]
}

func (h Heap) MakeMaxHeap() Heap {
	for i := h.index(h.LastParentNode()); i >= 0; i-- {
		curPNode := h[i]
		if h.shouldSwapWithChild(curPNode) {
			upNode := h.biggerChildNode(curPNode)
			h.flowUp(upNode)
		}
	}
	return h
}

func (h Heap) LastNode() Noder {
	return h[len(h)-1]
}

func (h Heap) deleteRootNode() Heap {
	h.swap(h.RootNode(), h.LastNode())
	h = h[0 : len(h)-1]
	h.flowDown(h.RootNode())
	return h
}

func (h Heap) PushNode(node Noder) Heap {
	h = append(h, node)
	h.MakeMaxHeap()
	return h
}
