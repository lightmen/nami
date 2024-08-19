package random

type option struct {
	randVal int32
}

type Option func(o *option)

func WithRandVal(val int32) Option {
	return func(o *option) {
		o.randVal = val
	}
}
