package aprometheus

import (
	"github.com/lightmen/nami/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var _ metrics.Counter = (*counter)(nil)

type counter struct {
	cv  *prometheus.CounterVec
	lvs []string
}

func NewCounter(cv *prometheus.CounterVec) metrics.Counter {
	return &counter{
		cv: cv,
	}
}

func (c *counter) With(lvs ...string) metrics.Counter {
	return &counter{
		cv:  c.cv,
		lvs: lvs,
	}
}

func (c *counter) Add(delta float64) {
	c.cv.WithLabelValues(c.lvs...).Add(delta)
}

func (c *counter) Inc() {
	c.cv.WithLabelValues(c.lvs...).Inc()
}
