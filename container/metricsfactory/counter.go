package metricsfactory

import (
	"fmt"
	"sync/atomic"
)

var _ Metric = &counter{}

// Counter stores a numerical value that only ever goes up.
type Counter interface {
	Metric
}

// CounterOpts is an alias of Opts.
type CounterOpts Opts

type counter struct {
	val int64
}

// NewCounter creates a new Counter based on the CounterOpts.
func NewCounter(opts CounterOpts) Counter {
	return &counter{}
}

func (c *counter) Add(val int64) {
	if val < 0 {
		panic(fmt.Errorf("stat/metric: cannot decrease in negative value. val: %d", val))
	}
	atomic.AddInt64(&c.val, val)
}

func (c *counter) Value() int64 {
	return atomic.LoadInt64(&c.val)
}

// CounterVecOpts is an alias of VectorOpts.
type CounterVecOpts VectorOpts

// CounterVec counter vec.
type CounterVec interface {
	// Inc increments the counter by 1. Use Add to increment it by arbitrary
	// non-negative values.
	Inc(labels ...string)
	// Add adds the given value to the counter. It panics if the value is <
	// 0.
	Add(v float64, labels ...string)
}