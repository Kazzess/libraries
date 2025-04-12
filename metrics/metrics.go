package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Registerer = prometheus.DefaultRegisterer
)

type (
	Opts          = prometheus.Opts
	CounterOpts   = prometheus.CounterOpts
	GaugeOpts     = prometheus.GaugeOpts
	SummaryOpts   = prometheus.SummaryOpts
	HistogramOpts = prometheus.HistogramOpts

	Counter     = prometheus.Counter
	CounterVec  = prometheus.CounterVec
	CounterFunc = prometheus.CounterFunc

	Gauge     = prometheus.Gauge
	GaugeVec  = prometheus.GaugeVec
	GaugeFunc = prometheus.GaugeFunc

	Summary    = prometheus.Summary
	SummaryVec = prometheus.SummaryVec

	Histogram    = prometheus.Histogram
	HistogramVec = prometheus.HistogramVec
)

func NewCounter(opts CounterOpts) Counter {
	return promauto.With(Registerer).NewCounter(opts)
}

func NewCounterVec(opts CounterOpts, labels []string) *CounterVec {
	return promauto.With(Registerer).NewCounterVec(opts, labels)
}

func NewCounterFunc(opts CounterOpts, fn func() float64) CounterFunc {
	return promauto.With(Registerer).NewCounterFunc(opts, fn)
}

func NewGauge(opts GaugeOpts) Gauge {
	return promauto.With(Registerer).NewGauge(opts)
}

func NewGaugeVec(opts GaugeOpts, labels []string) *GaugeVec {
	return promauto.With(Registerer).NewGaugeVec(opts, labels)
}

func NewGaugeFunc(opts GaugeOpts, fn func() float64) GaugeFunc {
	return promauto.With(Registerer).NewGaugeFunc(opts, fn)
}

func NewSummary(opts SummaryOpts) Summary {
	return promauto.With(Registerer).NewSummary(opts)
}

func NewSummaryVec(opts SummaryOpts, labels []string) *SummaryVec {
	return promauto.With(Registerer).NewSummaryVec(opts, labels)
}

func NewHistogram(opts HistogramOpts) Histogram {
	return promauto.With(Registerer).NewHistogram(opts)
}

func NewHistogramVec(opts HistogramOpts, labels []string) *HistogramVec {
	return promauto.With(Registerer).NewHistogramVec(opts, labels)
}
