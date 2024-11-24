package counter

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusCounter struct {
	counterVec  *prometheus.CounterVec
	labelValues []string
}

func NewPrometheusCounter(name, help string, labels []string) *PrometheusCounter {
	counterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name, // metric name
			Help: help,
		}, labels,
	)

	prometheus.MustRegister(counterVec)

	return &PrometheusCounter{
		counterVec:  counterVec,
		labelValues: labels,
	}
}

func (pc *PrometheusCounter) Increment(ctx context.Context, value float64) error {
	if value != 0 {
		pc.counterVec.WithLabelValues(pc.labelValues...).Add(value)
	} else {
		pc.counterVec.WithLabelValues(pc.labelValues...).Inc()
	}

	return nil
}

func (pc *PrometheusCounter) Set(ctx context.Context, value float64) error {
	return nil
}

func (pc *PrometheusCounter) Get(ctx context.Context) (float64, error) {
	return 0, nil
}

func (pc *PrometheusCounter) Reset(ctx context.Context) error {
	pc.counterVec.Reset()
	return nil
}

