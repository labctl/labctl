package utils

import "sync"

// Broadcast to many channels
type BroadcastChan[T any] struct {
	mu sync.Mutex
	ch map[chan<- T]bool
}

func NewBroadcastChan[T any]() *BroadcastChan[T] {
	return &BroadcastChan[T]{ch: make(map[chan<- T]bool)}
}

func (hub *BroadcastChan[T]) Broadcast(m T) {
	for c := range hub.ch {
		c <- m
	}
}

func (hub *BroadcastChan[T]) Add(ch chan<- T) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.ch[ch] = true
}

func (hub *BroadcastChan[T]) Remove(ch chan<- T) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	delete(hub.ch, ch)
}
