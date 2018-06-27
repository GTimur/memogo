package memogo

import (
	"fmt"
	"strings"
	"time"
)

func Reader() (err error) {
	tgrp := make(map[int64]string)
	grp := make(map[string]bool)

	// for every existing task
	for _, task := range GlobalTasks {
		tgrp[task.Memo.ID] = task.Group
		grp[task.Group] = true

	}

	for kg, _ := range grp {
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
