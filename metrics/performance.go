package metrics

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
)

var defaultRegistry = FunctionTimeRegistry{
	aggregates: make(map[string]FunctionTimeAggregate, 6),
}

// A FunctionTimeRegistry stores aggregates of duration and call data for specific functions
// Each unique function registered has its own aggregate
// The registry can then be used at any time (typically after some execution) to print
// performance stats for each item in the registry
type FunctionTimeRegistry struct {
	aggregates map[string]FunctionTimeAggregate
}

// A FunctionTimeAggregate represents the collected data for a function
// It is designed to be quick; expensive statistics will be calculated afterwards from it
type FunctionTimeAggregate struct {
	Sum   int64
	Max   int64
	Min   int64
	Count int64
}

// Func returns a tuple that represents a timeable event
func (r *FunctionTimeRegistry) Func(name string) (string, time.Time) {
	return name, time.Now()
}

// Func uses the default registry, see func (r *FunctionTimeRegistry) Func
func Func(name string) (string, time.Time) {
	return defaultRegistry.Func(name)
}

// Timed updates an aggregate for the given event, keeping track of the
// number of calls and some stats about the duration of the calls
// It MUST be deferred or called at the end of a function
// It's meant to take the output of a timeable event generator, like:
// defer r.Timed(r.Func("functionUnderTest"))
func (r *FunctionTimeRegistry) Timed(name string, start time.Time) {
	delta := time.Since(start).Nanoseconds()
	aggregate, ok := r.aggregates[name]
	if !ok {
		aggregate = FunctionTimeAggregate{Sum: 0, Max: delta, Min: delta, Count: 0}
	}
	aggregate.Count++
	aggregate.Sum += delta
	if delta > aggregate.Max {
		aggregate.Max = delta
	}
	if delta < aggregate.Min {
		aggregate.Min = delta
	}

	r.aggregates[name] = aggregate
}

// Timed uses the default registry, see func (r *FunctionTimeRegistry) Timed
func Timed(name string, start time.Time) {
	defaultRegistry.Timed(name, start)
}

type metricsRow struct {
	name                 string
	count, max, min, avg int64
}

// Output prints the statistics for all aggregated events in this registry to the given file
// It overrites the given file if present
func (r *FunctionTimeRegistry) Output(filepath string) error {
	if len(r.aggregates) == 0 {
		return fmt.Errorf("No metrics were collected")
	}

	file, err := os.Create(filepath)
	if err != nil {
		return errors.Wrap(err, "Failed to open registry output file")
	}

	fmt.Fprintf(file, "Registry output for %s\n", time.Now().Local().Format("2006-01-02 15:04:05"))
	rows := make([]metricsRow, 0, len(r.aggregates))
	for name, aggregate := range r.aggregates {
		rows = append(rows, metricsRow{
			name:  name,
			count: aggregate.Count,
			avg:   int64(float64(aggregate.Sum) / float64(aggregate.Count)),
			max:   aggregate.Max,
			min:   aggregate.Min,
		})
	}

	for i, r := range rows {
		if i%10 == 0 {
			fmt.Fprintf(file, "\n%30s   %10s   %10s   %10s   %10s\n\n", "Metric Name", "Count(##)", "Avg(ns)", "Max(ns)", "Min(ns)")
		}
		fmt.Fprintf(file, "%30s   %10d   %10d   %10d   %10d\n", r.name, r.count, r.avg, r.max, r.min)
	}

	return nil
}

// Output uses the default registry, see func (r *FunctionTimeRegistry) Output
func Output(filepath string) error {
	return defaultRegistry.Output(filepath)
}
