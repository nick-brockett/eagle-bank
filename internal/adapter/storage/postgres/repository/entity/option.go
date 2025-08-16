package entity

func ApplyOptions[E any](opts []Option[E], value E) {
	for _, applyOpt := range opts {
		applyOpt(value)
	}
}

type Option[E any] func(E)

type Options[E any] []Option[E]

func (o Options[E]) Merge(opts ...Option[E]) Options[E] {
	return append(o, opts...)
}
