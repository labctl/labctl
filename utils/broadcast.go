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
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for c := range hub.ch {
		c <- m
	}
}

func (hub *BroadcastChan[T]) Add(ch chan<- T) {
	hub.mu.Lock()
	hub.ch[ch] = true
	hub.mu.Unlock()
}

func (hub *BroadcastChan[T]) Remove(ch chan<- T) {
	hub.mu.Lock()
	delete(hub.ch, ch)
	hub.mu.Unlock()
}
