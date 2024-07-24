package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	issuedOrdersCounter prometheus.Counter
}

func NewPrometheusMetrics() *PrometheusMetrics {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "issued_orders_total",
		Help: "Total number of orders issued",
	})
	prometheus.MustRegister(counter)

	return &PrometheusMetrics{
		issuedOrdersCounter: counter,
	}
}

func (m *PrometheusMetrics) IncIssuedOrders() {
	m.issuedOrdersCounter.Inc()
}
