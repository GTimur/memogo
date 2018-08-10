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

// StringByID - return queue by grups in string for HTML
func (q *Queue) StringByID() string {
	var data []string
	var id map[int64]bool //memoID
	id = make(map[int64]bool)

	// Find every memoID
	for _, j := range *q {
		id[j.MemoID] = true
	}

	fmt.Println("ID=", id)

	// Generate data for every memoID
	for k := range id {
		data = append(data, fmt.Sprintln("ID:", k, "</br>"))
		for _, j := range *q {
			if k == j.MemoID {
				data = append(data, fmt.Sprintln(j, "</br>"))
			}
		}
	}

	return fmt.Sprintln(data)
}

// Clear
func (q *Queue) Clear() {
	*q = nil
}

// RemoveEvt - removes event from queue
func (q *Queue) RemoveEvt(event TaskPlan) (err error) {
	if len(*q) == 0 {
		return err
	}

	var tp []TaskPlan
	tp = *q

	for i, j := range *q {
		if j.MemoID == event.MemoID && j.Plan.Run.Equal(event.Plan.Run) {
			fmt.Println("Event to remove:", event)
			tp = append(tp[:i], tp[i+1:]...)
		}
	}
	*q = tp

	return err
}
