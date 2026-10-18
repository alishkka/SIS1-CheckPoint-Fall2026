package task

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrProcessingFailed  = errors.New("processing failed")
)

type State string

const (
	Pending    State = "Pending"   
	Processing State = "Processing" 
	Completed  State = "Completed" 
	Failed     State = "Failed"  
)

type Result struct {
	Count  int
	Min    float64
	Max    float64
	Sum    float64
	Mean   float64
	Median float64
}

type Task struct {
	ID     int
	Input  []float64
	State  State
	Result *Result // nil, пока задача не завершена
	Error  string  // причина ошибки (только для Failed)
}

func canChange(from, to State) bool {
	switch from {
	case Pending:
		return to == Processing
	case Processing:
		return to == Completed || to == Failed
	default:
		return false
	}
}

func (t *Task) SetState(next State) error {
	if !canChange(t.State, next) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, t.State, next)
	}
	t.State = next
	return nil
}

func ValidateInput(values []float64) error {
	if len(values) == 0 {
		return fmt.Errorf("%w: list of values is empty", ErrInvalidInput)
	}
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("%w: value %v is not a finite number", ErrInvalidInput, v)
		}
	}
	return nil
}
