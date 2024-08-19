package aprometheus

import (
	"github.com/lightmen/nami/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var _ metrics.Observer = (*summary)(nil)

type summary struct {
	sv  *prometheus.SummaryVec
	lvs []string
}

func NewSummary(sv *prometheus.SummaryVec) metrics.Observer {
	return &summary{
		sv: sv,
	}
}

func (s *summary) With(lvs ...string) metrics.Observer {
	return &summary{
		sv:  s.sv,
		lvs: lvs,
	}
}

func (s *summary) Observe(value float64) {
	s.sv.WithLabelValues(s.lvs...).Observe(value)
}
