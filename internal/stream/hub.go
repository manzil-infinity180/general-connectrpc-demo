package stream

import "sync"

type Hub[T any] struct {
	mu   sync.Mutex
	subs map[string][]chan T
}

func NewHub[T any]() *Hub[T] {
	return &Hub[T]{subs: make(map[string][]chan T)}
}

func (h *Hub[T]) Subscribe(key string) chan T {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan T, 10)
	h.subs[key] = append(h.subs[key], ch)
	return ch
}

func (h *Hub[T]) Publish(key string, msg T) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, ch := range h.subs[key] {
		ch <- msg
	}
}
