package restserver

import (
	"github.com/Naist4869/awesomeProject/container/metricsfactory"
	"github.com/Naist4869/awesomeProject/container/metricsfactory/prometheusfactory"
)

const (
	clientNamespace = "http_client"
	serverNamespace = "http_server"
)

var metrics prometheusfactory.PromMetrics

var (
	_metricServerReqDur = metrics.NewHistogramVec(&metricsfactory.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path", "caller", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
	_metricServerReqCodeTotal = metrics.NewCounterVec(&metricsfactory.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "caller", "method", "code"},
	})
	_metricServerBBR = metrics.NewCounterVec(&metricsfactory.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "",
		Name:      "bbr_total",
		Help:      "http server bbr total.",
		Labels:    []string{"url", "method"},
	})
	_metricClientReqDur = metrics.NewHistogramVec(&metricsfactory.HistogramVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http client requests duration(ms).",
		Labels:    []string{"path", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
	_metricClientReqCodeTotal = metrics.NewCounterVec(&metricsfactory.CounterVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http client requests code count.",
		Labels:    []string{"path", "method", "code"},
	})
)
