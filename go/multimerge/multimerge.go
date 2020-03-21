package multimerge

type MSorter struct {
	ListBundle []List
}

type Noder interface {
	LessThan(Noder) bool
	Equal(Noder) bool
}

type List interface{}

func New(lists []List) MSorter {
	return MSorter{
		ListBundle: lists,
	}
}

func (ms MSorter) TopK(k int) ([]Noder, error) {
	return nil, nil
}
