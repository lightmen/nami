package dqueue

type option struct {
	maxWorker int
}

type OptionFunc func(o *option)

func WithMaxWorker(num int) OptionFunc {
	return func(o *option) {
		if num <= maxIdleWorkers {
			o.maxWorker = maxIdleWorkers
			return
		}
		o.maxWorker = num
	}
}
