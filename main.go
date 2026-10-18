package main

import (
	"fmt"
	"math"

	"sis/task"
)

func main() {
	svc := task.NewService(task.StatsProcessor{})

	t1, _ := svc.Create([]float64{10, 20, 30})
	t2, _ := svc.Create([]float64{5, 1, 4, 2})
	fmt.Println("created:", t1.ID, t2.ID, "state:", t1.State)

	for _, t := range []*task.Task{t1, t2} {
		if err := svc.Process(t.ID); err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Printf("task %d [%s]: %+v\n", t.ID, t.State, *t.Result)
	}

	if _, err := svc.Create(nil); err != nil {
		fmt.Println("error:", err)
	}

	if _, err := svc.Get(999); err != nil {
		fmt.Println("error:", err)
	}

	if err := svc.Process(t1.ID); err != nil {
		fmt.Println("error:", err)
	}

	t3, _ := svc.Create([]float64{math.MaxFloat64, math.MaxFloat64})
	if err := svc.Process(t3.ID); err != nil {
		fmt.Printf("task %d [%s]: %v\n", t3.ID, t3.State, err)
	}
}
