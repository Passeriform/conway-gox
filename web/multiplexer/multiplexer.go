package multiplexer

type Multiplexer[T any] struct {
	source <-chan T
	sinks  []chan T
}

func Of[T any](source <-chan T) Multiplexer[T] {
	return Multiplexer[T]{source, make([]chan T, 0)}
}

func (m *Multiplexer[T]) ProvisionSink() <-chan T {
	sink := make(chan T)
	m.sinks = append(m.sinks, sink)
	return sink
}

func (m *Multiplexer[T]) Forward() {
	for {
		select {
		case in := <-m.source:
			for _, s := range m.sinks {
				s <- in
			}
		}
	}
}
