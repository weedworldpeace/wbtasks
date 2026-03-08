package main

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	PORT     = 6060
	GOGCPERC = 100
)

var registry = prometheus.NewRegistry()

var (
	memAlloc = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_mem_alloc_bytes",
		Help: "Current memory allocation",
	})

	memTotalAlloc = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "go_mem_total_alloc_bytes_total",
		Help: "Total memory allocation over lifetime",
	})

	memSys = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_mem_sys_bytes",
		Help: "Total memory obtained from system",
	})

	memHeapAlloc = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_mem_heap_alloc_bytes",
		Help: "Heap memory currently allocated",
	})

	memHeapSys = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_mem_heap_sys_bytes",
		Help: "Heap memory obtained from system",
	})

	memHeapIdle = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_mem_heap_idle_bytes",
		Help: "Heap memory that is idle",
	})

	memHeapInuse = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_mem_heap_inuse_bytes",
		Help: "Heap memory that is in use",
	})

	memStackInuse = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_mem_stack_inuse_bytes",
		Help: "Stack memory in use",
	})

	gcCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "go_gc_count_total",
		Help: "Number of completed GC cycles",
	})

	gcPauseTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "go_gc_pause_seconds_total",
		Help: "Total GC pause time",
	})

	gcPause = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "go_gc_pause_seconds",
		Help:    "GC pause time",
		Buckets: prometheus.ExponentialBuckets(0.0001, 2, 10),
	})

	gcLast = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_gc_last_seconds",
		Help: "Time of last GC in seconds since Unix epoch",
	})

	gcPercent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_gc_target_percent",
		Help: "GC target percentage",
	})

	nextGC = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_gc_next_bytes",
		Help: "Heap size when next GC will occur",
	})
)

var (
	prevTotalAlloc uint64
	prevNumGC      uint32
	prevPauseTotal uint64
)

func init() {
	registry.MustRegister(memAlloc)
	registry.MustRegister(memTotalAlloc)
	registry.MustRegister(memSys)
	registry.MustRegister(memHeapAlloc)
	registry.MustRegister(memHeapSys)
	registry.MustRegister(memHeapIdle)
	registry.MustRegister(memHeapInuse)
	registry.MustRegister(memStackInuse)
	registry.MustRegister(gcCount)
	registry.MustRegister(gcPauseTotal)
	registry.MustRegister(gcPause)
	registry.MustRegister(gcLast)
	registry.MustRegister(gcPercent)
	registry.MustRegister(nextGC)
}

func alloc(w http.ResponseWriter, r *http.Request) {
	data := make([][]byte, 0)
	for i := 0; i < 10; i++ {
		data = append(data, make([]byte, 10<<20))
		time.Sleep(100 * time.Millisecond)
	}
	w.Write([]byte("allocated 10mB"))
}

func updateMetrics() {
	var memStats runtime.MemStats

	for {
		runtime.ReadMemStats(&memStats)

		memAlloc.Set(float64(memStats.Alloc))
		memSys.Set(float64(memStats.Sys))
		memHeapAlloc.Set(float64(memStats.HeapAlloc))
		memHeapSys.Set(float64(memStats.HeapSys))
		memHeapIdle.Set(float64(memStats.HeapIdle))
		memHeapInuse.Set(float64(memStats.HeapInuse))
		memStackInuse.Set(float64(memStats.StackInuse))
		nextGC.Set(float64(memStats.NextGC))
		gcPercent.Set(float64(debug.SetGCPercent(GOGCPERC)))
		gcLast.Set(float64(memStats.LastGC) / 1e9)

		if memStats.TotalAlloc > prevTotalAlloc {
			memTotalAlloc.Add(float64(memStats.TotalAlloc - prevTotalAlloc))
			prevTotalAlloc = memStats.TotalAlloc
		}

		if memStats.NumGC > prevNumGC {
			gcCount.Add(float64(memStats.NumGC - prevNumGC))
			prevNumGC = memStats.NumGC
		}

		if memStats.PauseTotalNs > prevPauseTotal {
			pauseDiff := float64(memStats.PauseTotalNs-prevPauseTotal) / 1e9
			gcPauseTotal.Add(pauseDiff)
			prevPauseTotal = memStats.PauseTotalNs
		}

		for i := 0; i < 256; i++ {
			if memStats.PauseNs[i] > 0 {
				gcPause.Observe(float64(memStats.PauseNs[i]) / 1e9)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	debug.SetGCPercent(GOGCPERC)

	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/alloc", alloc)

	go updateMetrics()

	fmt.Printf("Server starting on :%d\n", PORT)
	fmt.Printf("Metrics: http://localhost:%d/metrics\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
