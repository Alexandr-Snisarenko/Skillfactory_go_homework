package main

import (
	"errors"
	"log/slog"
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
		slog.Debug("RingBufer", "operation", "Push", "action", "shift")
		for i := 1; i < b.size; i++ {
			b.array[i-1] = b.array[i]
		}
		slog.Debug("RingBufer", "operation", "Push", "action", "write_value", "value", v, "position", b.pos)
		b.array[b.pos] = v
	} else {
		b.pos++
		slog.Debug("RingBufer", "operation", "Push", "action", "write_value", "value", v, "position", b.pos)
		b.array[b.pos] = v
	}
}

func (b *RingBufferInt) GetAll() []int {
	if b.pos == -1 {
		slog.Debug("RingBufer", "operation", "GetAll", "action", "return_by_empty")
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	arr := b.array[:b.pos+1]
	slog.Debug("RingBufer", "operation", "GetAll", "action", "return_data_slice", "value", arr)
	b.pos = -1
	return arr
}

func (b *RingBufferInt) Get() (int, error) {
	if b.pos == -1 {
		slog.Debug("RingBufer", "operation", "Get", "action", "return_by_empty")
		return 0, errors.New("buffer is empty")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	getVal := b.array[b.pos]
	slog.Debug("RingBufer", "operation", "Get", "action", "return_value", "value", getVal, "position", b.pos)
	b.pos--
	return getVal, nil
}

func NewRingBufferInt(size int) *RingBufferInt {
	return &RingBufferInt{make([]int, size), -1, size, sync.Mutex{}}
}
