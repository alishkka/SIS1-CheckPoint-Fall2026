package task

import "fmt"

type Service struct {
	tasks     map[int]*Task 
	nextID    int
	processor Processor
}

func NewService(p Processor) *Service {
	return &Service{
		tasks:     make(map[int]*Task),
		nextID:    1,
		processor: p,
	}
}

func (s *Service) Create(input []float64) (*Task, error) {
	if err := ValidateInput(input); err != nil {
		return nil, err
	}

	data := make([]float64, len(input))
	copy(data, input)

	t := &Task{ID: s.nextID, Input: data, State: Pending}
	s.tasks[t.ID] = t
	s.nextID++
	return t, nil
}

func (s *Service) Get(id int) (*Task, error) {
	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("%w: id %d", ErrTaskNotFound, id)
	}
	return t, nil
}

func (s *Service) Process(id int) error {
	t, err := s.Get(id)
	if err != nil {
		return err
	}

	if err := t.SetState(Processing); err != nil {
		return err
	}

	result, err := s.processor.Process(t.Input)
	if err != nil {
		t.Error = err.Error()
		t.SetState(Failed)
		return fmt.Errorf("%w: %v", ErrProcessingFailed, err)
	}

	t.Result = &result
	return t.SetState(Completed)
}
