package metricsfactory

import (
	"net/http"
)

type Metrics interface {
	NewHistogramVec(cfg *HistogramVecOpts) HistogramVec
	NewGaugeVec(cfg *GaugeVecOpts) GaugeVec
	NewCounterVec(cfg *CounterVecOpts) CounterVec
	Handler() http.Handler
}
