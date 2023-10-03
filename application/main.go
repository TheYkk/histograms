package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/fjl/memsize/memsizeui"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var memsizeH memsizeui.Handler

	var histogramRandom = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "random_numbers",
		Help:    "A histogram of normally distributed random numbers.",
		Buckets: prometheus.LinearBuckets(-3, .1, 61),
	})

	histogram := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "histogram_metric",
		Buckets: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
	})
	// Create a Prometheus summary metric for HTTP request latency.
	detailReq := metrics.NewHistogram("latency")
	requestLatencySummary := promauto.NewSummary(
		prometheus.SummaryOpts{
			Name:       "http_request_latency_seconds",
			Help:       "HTTP request latency in seconds.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001},
		},
	)
	memsizeH.Add("metr", detailReq)
	okCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "ok_counter",
	})
	koCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "ko_counter",
	})
	req := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "req_metric",
		Buckets: prometheus.DefBuckets,
	})
	fail := 0
	go func() {
		for {
			rao := rand.Float64() * 2.25
			if rand.Intn(2) == 1 {
				rao += 0.370
				if rao > 2.5 {
					fail += 1
					println("FAILOO", fail)
					koCounter.Inc()
				} else {
					okCounter.Inc()
				}
			}
			detailReq.Update(rao)
			requestLatencySummary.Observe(rao)
			histogram.Observe(rao)
			req.Observe(rao)
			histogramRandom.Observe(rand.NormFloat64())
			time.Sleep(4 * time.Millisecond)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":8282", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			metrics.WritePrometheus(w, true)
		}))
	}()
	metrics.InitPush("http://localhost:8428/api/v1/import/prometheus", time.Second, "", false)
	http.Handle("/memsize/", http.StripPrefix("/memsize", &memsizeH))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
