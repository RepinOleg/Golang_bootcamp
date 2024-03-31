package ex02

import (
	"container/heap"
	"fmt"
)

type Present struct {
	Value int
	Size  int
}

type PresentHeap []Present

func (p PresentHeap) Len() int { return len(p) }

func (p PresentHeap) Less(i, j int) bool {
	if p[i].Value == p[j].Value {
		return p[i].Size < p[j].Size
	}
	return p[i].Value > p[j].Value
}

func (p PresentHeap) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p *PresentHeap) Push(x interface{}) {
	*p = append(*p, x.(Present))
}

func (p *PresentHeap) Pop() interface{} {
	lenOldHeap := len(*p)
	popElem := (*p)[lenOldHeap-1] // сохраняем последний элемент
	*p = (*p)[0 : lenOldHeap-1]
	return popElem
}

func GetNCoolestPresents(in []Present, n int) ([]Present, error) {
	if n < 0 || n > len(in) {
		return nil, fmt.Errorf("n = %d, index out of range", n)
	}
	tmp := make(PresentHeap, len(in))
	for i, present := range in {
		tmp[i].Value = present.Value
		tmp[i].Size = present.Size
	}
	heap.Init(&tmp)
	res := make([]Present, n)
	for i := range res {
		res[i] = heap.Pop(&tmp).(Present)
	}
	return res, nil
}
