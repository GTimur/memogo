package memogo

import (
	"fmt"
	"sort"
	"time"
)

/*
  plan.go realize reaction plan by using queue

  - analyze of Tasks array
  - make queue
*/

// Plan describe NOW/NEXT/PREV run time
type Plan struct {
	Run time.Time // Run - when it needs to be executed
}

// TaskPlan describe run time of task
type TaskPlan struct {
	MemoID int64
	Plan   Plan
}

// Queue - queue of Tasks run-plans
type Queue []TaskPlan

// MakeQueue
func (q *Queue) MakeQueue() (err error) {
	var tp TaskPlan

	q.Clear()

	// map[TaskID]map[MemoID]time.Time
	for k, v := range GlobalTimeMap {
		for _, t := range v {
			tp.Plan.Run = t
			tp.MemoID = k
			*q = append(*q, tp)
		}
	}

	q.Sort()

	return err
}

// sort
func (s Queue) Less(i, j int) bool { return s[i].Plan.Run.Before(s[j].Plan.Run) }
func (s Queue) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Queue) Len() int           { return len(s) }

// SortQ - sort Q by date
func (q *Queue) Sort() {

	sort.Sort(q)
}

// Print
func (q *Queue) Print() {
	fmt.Println("FormQueue():", *q)
}

// String - return queue in string for HTML
func (q *Queue) String() string {
	var data []string

	for _, j := range *q {
		data = append(data, fmt.Sprintln(j, "</br>"))
	}

	return fmt.Sprintln(data)
}

// Clear
func (q *Queue) Clear() {
	*q = nil
}
