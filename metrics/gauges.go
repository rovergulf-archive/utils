package metrics

import "github.com/prometheus/client_golang/prometheus"

func StatusGauge() prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "app",
		Name:      "status",
		Help:      "Describes application status",
	})
}
