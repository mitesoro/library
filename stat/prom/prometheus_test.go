package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})

	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		}, []string{"device"})

	// 初始化 web_reqeust_total， counter类型指标， 表示接收http请求总次数
	WebRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "web_reqeust_total",
			Help: "Number of hello requests in total",
		},
		// 设置两个标签 请求方法和 路径 对请求总次数在两个
		[]string{"method", "endpoint"})

	// web_request_duration_seconds，Histogram类型指标，bucket代表duration的分布区间
	WebRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "web_request_duration_seconds",
			Help:    "web request duration distribution",
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
		},
		[]string{"method", "endpoint"})
)

func init() {

	// 注册监控指标
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
	prometheus.MustRegister(WebRequestTotal)
	prometheus.MustRegister(WebRequestDuration)
}

// 包装 handler function,不侵入业务逻辑

func Monitor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		duration := time.Since(start)
		// counter类型 metric的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		// Histogram类型 metric的记录方式
		WebRequestDuration.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Observe(duration.Seconds())
	}
}

func Query(w http.ResponseWriter, r *http.Request) {
	//模拟业务查询耗时0~1s
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	_, _ = io.WriteString(w, "some results")
}

func TestPrometheus(t *testing.T) {
	cpuTemp.Set(65.3)
	hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/query", Monitor(Query))
	log.Fatal(http.ListenAndServe(":8888", nil))
	// pro采集数据是通过定期请求该服务http端口来实现的。
}
