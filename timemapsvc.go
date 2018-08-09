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
func GetQueueEvent(queue Queue) (q Queue, err error) {
	q = q[:0]
	if queue.Len() == 0 {
		return q, errors.New("Queue contains no event")
	}
	alone := false
	ts := 1 * 60 // +-1 minutes time shift (+1 -1 = 2 min)
	now := time.Now()
	start := now.Add(time.Second * time.Duration(ts*-1))
	end := now.Add(time.Second * time.Duration(ts))
	event1, event2 := now, now

	for i := 0; i < queue.Len(); i++ {
		if queue.Len() >= i+1 {
			event1 = queue[i].Plan.Run
			event2 = queue[i+1].Plan.Run
			alone = false
		} else {
			event1 = queue[i].Plan.Run
			alone = true
			//event2 = queue[i].Plan.Run.Add(time.Second * time.Duration(5*60)) //5 min
		}

		if event1.Before(start) || event1.After(end) {
			return q, err
		}

		// If skipped event1 time (T) < event2-event1 then it will be processed
		//time.Since(start).Minutes() < event2.Sub(event1).Minutes()
		fmt.Println("DEBUG(event1,2):", event1, event2)
		if !event1.Before(start) && !event1.After(end) {
			q = append(q, queue[i])
			if (event2.Sub(event1).Seconds() < 60 || event1.Equal(event2)) && !alone {
				fmt.Println("EVT DIFF:", event2.Sub(event1).Seconds())
				q = append(q, queue[i+1])
			}
			// if time between event2 and event1 > 3 min - skip
			if event2.Sub(event1).Minutes() > 3 {
				return q, err
			}
		}
	}

	//	fmt.Println("EVENT2:", memoid)
	return q, err
}
