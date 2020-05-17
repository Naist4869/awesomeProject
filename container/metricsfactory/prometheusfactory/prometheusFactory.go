package prometheusfactory

import (
	"net/http"

	"github.com/Naist4869/awesomeProject/container/metricsfactory"

	"github.com/Naist4869/awesomeProject/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct {
}

func (p PromMetrics) NewHistogramVec(cfg *metricsfactory.HistogramVecOpts) metricsfactory.HistogramVec {
	return NewHistogramVec(cfg)
}

func (p PromMetrics) NewGaugeVec(cfg *metricsfactory.GaugeVecOpts) metricsfactory.GaugeVec {
	return NewGaugeVec(cfg)
}

func (p PromMetrics) NewCounterVec(cfg *metricsfactory.CounterVecOpts) metricsfactory.CounterVec {
	return NewCounterVec(cfg)
}

func (p PromMetrics) Handler() http.Handler {
	return promhttp.Handler()
}
func PrometheusFactory(mc config.MetricsConfig) (metricsfactory.Metrics, error) {
	return PromMetrics{}, nil
}
