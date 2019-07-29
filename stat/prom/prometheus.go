package prom

import "github.com/prometheus/client_golang/prometheus"

// https://godoc.org/github.com/prometheus/client_golang/prometheus

var (
	BusinessInfoCount = New().WithCounter("go_business_info_count", []string{"name"}).WithState("go_business_info_state", []string{"name"})
)

// Prom 结构
// CounterVec是用来管理相同metric下不同label的一组Counter
type Prom struct {
	timer   *prometheus.HistogramVec //Histogram,Summary提供跟多的统计信息;典型的应用如：请求持续时间，响应大小。可以对观察结果采样，分组及统计
	counter *prometheus.CounterVec   //Counter对数据只增不减;典型的应用如：请求的个数，结束的任务数，出现的错误数等。随着客户端不断请求，数值越来越大。
	state   *prometheus.GaugeVec     // Gauge可增可减;典型的应用如：温度，运行的 goroutines 的个数。返回的数值会上下波动。
}

// New 新建
func New() *Prom {
	return &Prom{}
}

// WithTimer with summary timer
func (p *Prom) WithTimer(name string, labels []string) *Prom {
	if p == nil || p.timer != nil {
		return p
	}
	p.timer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.timer)
	return p
}

// WithCounter sets counter.
func (p *Prom) WithCounter(name string, labels []string) *Prom {
	if p == nil || p.counter != nil {
		return p
	}
	p.counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.counter)
	return p
}

// WithState sets state.
func (p *Prom) WithState(name string, labels []string) *Prom {
	if p == nil || p.state != nil {
		return p
	}
	p.state = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: name,
		}, labels)
	prometheus.MustRegister(p.state)
	return p
}

// Timing log timing information (in milliseconds) without sampling
func (p *Prom) Timing(name string, time int64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.timer != nil {
		p.timer.WithLabelValues(label...).Observe(float64(time))
	}
}

// Incr increments one stat counter without sampling
func (p *Prom) Incr(name string, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counter != nil {
		p.counter.WithLabelValues(label...).Inc()
	}
	if p.state != nil {
		p.state.WithLabelValues(label...).Inc()
	}
}

// Decr decrements one stat counter without sampling
func (p *Prom) Decr(name string, extra ...string) {
	if p.state != nil {
		label := append([]string{name}, extra...)
		p.state.WithLabelValues(label...).Dec()
	}
}

// State set state
func (p *Prom) State(name string, v int64, extra ...string) {
	if p.state != nil {
		label := append([]string{name}, extra...)
		p.state.WithLabelValues(label...).Set(float64(v))
	}
}

// Add add count    v must > 0
func (p *Prom) Add(name string, v int64, extra ...string) {
	label := append([]string{name}, extra...)
	if p.counter != nil {
		p.counter.WithLabelValues(label...).Add(float64(v))
	}
	if p.state != nil {
		p.state.WithLabelValues(label...).Add(float64(v))
	}
}
