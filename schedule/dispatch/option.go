package dispatch

type Option struct {
	maxWorker int
}

type OptionFunc func(o *Option)

func WithMaxWorker(maxWorker int) OptionFunc {
	return func(o *Option) {
		if maxWorker <= 0 {
			o.maxWorker = 128
			return
		}
		o.maxWorker = maxWorker
	}
}
