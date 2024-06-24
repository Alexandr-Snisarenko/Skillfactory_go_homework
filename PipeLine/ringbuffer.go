package main

import (
	"errors"
	"sync"
)

// кольцевой буфер
type RingBufferInt struct {
	array []int
	pos   int
	size  int
	mu    sync.Mutex
}

func (b *RingBufferInt) Push(v int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.pos == b.size-1 {
		for i := 1; i < b.size; i++ {
			b.array[i-1] = b.array[i]
		}
		b.array[b.pos] = v
	} else {
		b.pos++
		b.array[b.pos] = v
	}
}

func (b *RingBufferInt) GetAll() []int {
	if b.pos == -1 {
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	arr := b.array[:b.pos+1]
	b.pos = -1
	return arr
}

func (b *RingBufferInt) Get() (int, error) {
	if b.pos == -1 {
		return 0, errors.New("buffer is empty")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.pos--
	return b.pos, nil
}

func NewRingBufferInt(size int) *RingBufferInt {
	return &RingBufferInt{make([]int, size), -1, size, sync.Mutex{}}
}
