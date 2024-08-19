package aprometheus

import (
	"github.com/lightmen/nami/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var _ metrics.Gauge = (*gauge)(nil)

type gauge struct {
	gv  *prometheus.GaugeVec
	lvs []string
}

func NewGauge(gv *prometheus.GaugeVec) metrics.Gauge {
	return &gauge{
		gv: gv,
	}
}

func (g *gauge) With(lvs ...string) metrics.Gauge {
	return &gauge{
		gv:  g.gv,
		lvs: lvs,
	}
}

func (g *gauge) Add(delta float64) {
	g.gv.WithLabelValues(g.lvs...).Add(delta)
}

func (g *gauge) Set(value float64) {
	g.gv.WithLabelValues(g.lvs...).Set(value)
}

func (g *gauge) Sub(delta float64) {
	g.gv.WithLabelValues(g.lvs...).Sub(delta)
}
