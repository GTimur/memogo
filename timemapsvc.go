package memogo

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

/*not used*/
func GrpReader() (err error) {
	tgrp := make(map[int64]string)
	grp := make(map[string]bool)

	// for every existing task
	for _, task := range GlobalTasks {
		tgrp[task.Memo.ID] = task.Group
		grp[task.Group] = true

	}

	for kg := range grp {
		for _, task := range GlobalTasks {
			//for selected group
			if strings.EqualFold(task.Group, kg) {
				CurrentAction(task)
				fmt.Println("TASK:", task.ID, kg)
			}
		}
	}

	return err
}

func CurrentAction(task Task) (err error) {
	var minid int64
	var value time.Time

	minid = -1
	for memoid, gtm := range GlobalTimeMap {
		if memoid != task.Memo.ID {
			continue
		}
		fmt.Println(memoid)
		// initialize minid by first value for minid from map
		for id := range gtm {
			if minid < id {
				minid = id
				break
			}
		}
		// search minimal id
		now := time.Now()
		for id, t := range gtm {
			if (minid > id) && (!t.Before(now)) {
				minid = id
				value = t
			}
		}
		fmt.Println("MIN:", minid, value)
		minid = -1
	}

	return err
}

// GetEvent - takes event from queue if it ready to process
func GetQueueEvent(queue Queue) (memoid []int64, err error) {
	if queue.Len() == 0 {
		return memoid, errors.New("Queue contains no event")
	}
	ts := 1 * 60 // +-1 minutes time shift (+1 -1 = 2 min)
	now := time.Now()
	start := now.Add(time.Second * time.Duration(ts*-1))
	end := now.Add(time.Second * time.Duration(ts))
	event1, event2 := now, now

	if queue.Len() > 1 {
		event1 = queue[0].Plan.Run
		event2 = queue[1].Plan.Run
	}
	event1 = queue[0].Plan.Run
	event2 = queue[0].Plan.Run.Add(time.Second * time.Duration(5*60)) //5 min

	// If event(s) time is skipped and next event is too far, then
	// event(s) will be selected to process
	if time.Since(start).Minutes() < event2.Sub(event1).Minutes() {

	}

	fmt.Println("START:", start)
	fmt.Println("START:", end)
	fmt.Println("PROPUSH:")
	fmt.Println("DURATION:", event2.Sub(event1).Minutes())
	fmt.Println("EVENT1:", event1)
	fmt.Println("EVENT2:", event2)

	//if queue[0].Plan.Run

	return
}
