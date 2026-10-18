package task

import (
	"errors"
	"math"
	"testing"
)

func TestCreateSuccess(t *testing.T) {
	svc := NewService(StatsProcessor{})

	got, err := svc.Create([]float64{10, 20, 30})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
	if got.State != Pending {
		t.Errorf("State = %s, want Pending", got.State)
	}
	if got.Result != nil {
		t.Errorf("Result = %v, want nil", got.Result)
	}
}

func TestCreateInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
	}{
		{"nil slice", nil},
		{"empty slice", []float64{}},
		{"NaN", []float64{1, math.NaN()}},
		{"+Inf", []float64{math.Inf(1)}},
		{"-Inf", []float64{math.Inf(-1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(StatsProcessor{})
			_, err := svc.Create(tt.input)
			if !errors.Is(err, ErrInvalidInput) {
				t.Errorf("error = %v, want ErrInvalidInput", err)
			}
		})
	}
}

func TestGetExisting(t *testing.T) {
	svc := NewService(StatsProcessor{})
	created, _ := svc.Create([]float64{5})

	got, err := svc.Get(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %d, want %d", got.ID, created.ID)
	}
}

func TestGetNotFound(t *testing.T) {
	svc := NewService(StatsProcessor{})

	_, err := svc.Get(42)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("error = %v, want ErrTaskNotFound", err)
	}
}

func TestProcessSuccess(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		want  Result
	}{
		{"odd count", []float64{10, 20, 30}, Result{Count: 3, Min: 10, Max: 30, Sum: 60, Mean: 20, Median: 20}},
		{"even count", []float64{4, 1, 3, 2}, Result{Count: 4, Min: 1, Max: 4, Sum: 10, Mean: 2.5, Median: 2.5}},
		{"one value", []float64{7}, Result{Count: 1, Min: 7, Max: 7, Sum: 7, Mean: 7, Median: 7}},
		{"negative", []float64{-5, 5}, Result{Count: 2, Min: -5, Max: 5, Sum: 0, Mean: 0, Median: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(StatsProcessor{})
			created, _ := svc.Create(tt.input)

			if err := svc.Process(created.ID); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			stored, _ := svc.Get(created.ID)
			if stored.State != Completed {
				t.Errorf("State = %s, want Completed", stored.State)
			}
			if stored.Result == nil {
				t.Fatal("Result is nil")
			}
			if *stored.Result != tt.want {
				t.Errorf("Result = %+v, want %+v", *stored.Result, tt.want)
			}
		})
	}
}

func TestProcessNotFound(t *testing.T) {
	svc := NewService(StatsProcessor{})

	err := svc.Process(1)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("error = %v, want ErrTaskNotFound", err)
	}
}

func TestProcessTwice(t *testing.T) {
	svc := NewService(StatsProcessor{})
	created, _ := svc.Create([]float64{1, 2, 3})
	svc.Process(created.ID)

	err := svc.Process(created.ID)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("error = %v, want ErrInvalidTransition", err)
	}
	if created.State != Completed {
		t.Errorf("State = %s, want Completed", created.State)
	}
}

func TestSetState(t *testing.T) {
	tests := []struct {
		from, to State
		wantErr  bool
	}{
		{Pending, Processing, false},
		{Processing, Completed, false},
		{Processing, Failed, false},
		{Pending, Completed, true},
		{Processing, Pending, true},
		{Completed, Processing, true},
		{Failed, Processing, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+" to "+string(tt.to), func(t *testing.T) {
			task := &Task{State: tt.from}
			err := task.SetState(tt.to)

			if tt.wantErr && !errors.Is(err, ErrInvalidTransition) {
				t.Errorf("error = %v, want ErrInvalidTransition", err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestProcessFailed(t *testing.T) {
	svc := NewService(StatsProcessor{})
	created, _ := svc.Create([]float64{math.MaxFloat64, math.MaxFloat64})

	err := svc.Process(created.ID)
	if !errors.Is(err, ErrProcessingFailed) {
		t.Fatalf("error = %v, want ErrProcessingFailed", err)
	}
	if created.State != Failed {
		t.Errorf("State = %s, want Failed", created.State)
	}
	if created.Result != nil {
		t.Errorf("Result = %v, want nil", created.Result)
	}
}

func TestIndependentTasks(t *testing.T) {
	svc := NewService(StatsProcessor{})
	a, _ := svc.Create([]float64{1, 2, 3})
	b, _ := svc.Create([]float64{10, 20})

	if a.ID == b.ID {
		t.Fatalf("IDs must be different, got %d and %d", a.ID, b.ID)
	}

	svc.Process(a.ID)

	if a.Result == nil || a.Result.Sum != 6 {
		t.Errorf("a.Result = %v, want sum 6", a.Result)
	}
	if b.State != Pending || b.Result != nil {
		t.Errorf("b must stay untouched: state=%s result=%v", b.State, b.Result)
	}

	svc.Process(b.ID)

	if b.Result == nil || b.Result.Sum != 30 {
		t.Errorf("b.Result = %v, want sum 30", b.Result)
	}
	if a.Result.Sum != 6 {
		t.Errorf("a.Result changed: sum = %v, want 6", a.Result.Sum)
	}
}

func TestInputIsCopied(t *testing.T) {
	svc := NewService(StatsProcessor{})
	input := []float64{1, 2, 3}
	created, _ := svc.Create(input)

	input[0] = 999

	if created.Input[0] != 1 {
		t.Errorf("Input[0] = %v, want 1", created.Input[0])
	}
}
