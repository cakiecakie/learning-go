package main

import "fmt"

// min heap
type Heap []int

func NewHeap(length int) Heap {
	return make([]int, length+1)
}

func (h Heap) Len() int {
	return h[0]
}

func (h Heap) Peak() int {
	if h.Len() == 0 {
		return -1
	}

	return h[1]
}

func (h Heap) Push(num int) {

	if h.Len() == len(h) {
		h.Pop()
	}

	idx := h[0] + 1
	h[idx] = num

	for idx != 1 {
		parentIdx := idx / 2
		if h[parentIdx] > h[idx] {
			h[parentIdx], h[idx] = h[idx], h[parentIdx]
			idx = parentIdx
		}
	}

	h[0] += 1
}

func (h Heap) Pop() int {

	if h.Len() == 0 {
		return -1
	}

	idx := 1
	res := h[idx]
	for 2*idx+1 <= h.Len() {
		if h[2*idx] < h[2*idx+1] {
			h[idx] = h[2*idx]
			idx = 2 * idx
		} else {
			h[idx] = h[2*idx+1]
			idx = 2*idx + 1
		}
	}

	if 2*idx <= h.Len() {
		h[idx] = h[2*idx]
	}

	h[0] -= 1
	return res
}

func main() {

	h := NewHeap(4)

	for i := 4; i > 0; i-- {
		h.Push(i)
		fmt.Printf("push: %+v\n", h)
		if h.Peak() != i {
			fmt.Printf("error push exp: %d, got %d", i, h.Peak())
		}
	}

	for i := 1; i <= 4; i++ {
		num := h.Pop()
		fmt.Printf("pop: %+v\n", h)
		if num != i {
			fmt.Printf("error pop exp: %d, got %d", i, num)
		}
	}

	for i := 4; i > 0; i-- {
		h.Push(i)
		fmt.Printf("push: %+v\n", h)
		if h.Peak() != i {
			fmt.Printf("error push exp: %d, got %d", i, h.Peak())
		}
	}
}
