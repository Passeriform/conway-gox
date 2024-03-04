package multiplexer

import (
	"sync"

	"github.com/passeriform/conway-gox/internal/utility"
)

const (
	multiplexerIdLength = 5
)

type Multiplexer[T any] struct {
	source <-chan T
	sinks  map[string]chan T
	done   chan struct{}
	once   sync.Once
}

func Of[T any](source <-chan T) Multiplexer[T] {
	return Multiplexer[T]{source: source, sinks: make(map[string]chan T), done: make(chan struct{})}
}

func (m *Multiplexer[T]) ProvisionSink() (string, <-chan T) {
	sink := make(chan T)
	id := utility.NewRandomString(multiplexerIdLength)
	m.sinks[id] = sink
	return id, sink
}

func (m *Multiplexer[T]) CloseSink(id string) {
	close(m.sinks[id])
	delete(m.sinks, id)
}

func (m *Multiplexer[T]) Forward() {
	defer m.Close()
	for {
		select {
		case <-m.done:
			return
		case in := <-m.source:
			for _, s := range m.sinks {
				s <- in
			}
		}
	}
}

func (m *Multiplexer[T]) Close() {
	m.once.Do(func() {
		for id := range m.sinks {
			close(m.sinks[id])
		}
		m.done <- struct{}{}
	})
}
