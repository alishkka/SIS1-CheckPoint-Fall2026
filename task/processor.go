package task

import (
	"errors"
	"math"
	"sort"
)

type Processor interface {
	Process(input []float64) (Result, error)
}

type StatsProcessor struct{}

func (StatsProcessor) Process(input []float64) (Result, error) {
	if err := ValidateInput(input); err != nil {
		return Result{}, err
	}

	sorted := make([]float64, len(input))
	copy(sorted, input)
	sort.Float64s(sorted)

	sum := 0.0
	for _, v := range sorted {
		sum += v
	}
	if math.IsInf(sum, 0) {
		return Result{}, errors.New("sum is too large")
	}

	n := len(sorted)
	var median float64
	if n%2 == 1 {
		median = sorted[n/2]
	} else {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	}

	return Result{
		Count:  n,
		Min:    sorted[0],
		Max:    sorted[n-1],
		Sum:    sum,
		Mean:   sum / float64(n),
		Median: median,
	}, nil
}
