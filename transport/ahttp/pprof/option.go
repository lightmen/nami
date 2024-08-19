package pprof

type Option func(*Transporter)

func WithParam(key, val string) Option {
	return func(t *Transporter) {
		if t.params == nil {
			t.params = make(map[string]string)
		}
		t.params[key] = val
	}
}

func WithRouter(r IRouter) Option {
	return func(t *Transporter) {
		t.router = r
	}
}
