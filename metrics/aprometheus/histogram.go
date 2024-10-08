package aprometheus

import (
	"github.com/lightmen/nami/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var _ metrics.Observer = (*histogram)(nil)

type histogram struct {
	hv  *prometheus.HistogramVec
	lvs []string
}

func NewHistogram(hv *prometheus.HistogramVec) metrics.Observer {
	return &histogram{
		hv: hv,
	}
}

func (h *histogram) With(lvs ...string) metrics.Observer {
	return &histogram{
		hv:  h.hv,
		lvs: lvs,
	}
}

func (h *histogram) Observe(value float64) {
	h.hv.WithLabelValues(h.lvs...).Observe(value)
}
