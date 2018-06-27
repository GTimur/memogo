package memogo

import "time"

/*
  plan.go realize reaction plan by using queue

  - analyze of Tasks array
  - make queue
*/

// Plan describe NOW/NEXT/PREV run time
type Plan struct {
	Run  time.Time // Run - when it needs to be executed
	Next time.Time
	Prev time.Time
}

// TaskPlan describe run time of task
type TaskPlan struct {
	Task Task
	Plan Plan
}

// Queue - queue of Tasks run-plans
type Queue []TaskPlan

// Analize - analyze Tasks
func (q *Queue) Analize() (err error) {
	QueueGlobal = QueueGlobal[0:0]

	//for r, v := range TasksGlobal {

	//}

	return
}
