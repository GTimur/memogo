package memogo

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func InitGlobalTimeMap() {
	GlobalTimeMap = make(map[int64]map[int64]time.Time) // map[TaskID]map[MemoID]time.Time
	GlobalTimeCount = make(map[int64]map[int64]int)     // map[TaskID]map[MemoID]int
}

func BuildTimeMap() (err error) {
	//cleaning maps
	InitGlobalTimeMap()

	//filling data into maps
	for _, task := range GlobalTasks {
		err = BuildTimeMapTask(task)
		if err != nil {
			//fmt.Printf("BuildTimeMap error: %v\n", err)
			GlobalConfig.LogFile.Add(fmt.Sprintf("BuildTimeMap error: %v", err))
			continue
		}
	}
	return nil //err
}

//BuildTimeMapTask - fills GlobalTimeMap for task
func BuildTimeMapTask(t Task) (err error) {
	g := 0                                         //time granula
	c := 0                                         //repeat count
	before := t.Memo.Scenario.FreqBefore.Value > 0 //Уведомлять до начала события
	after := t.Memo.Scenario.FreqAfter.Value > 0   //Уведомлять после окончания события
	till := t.Memo.Scenario.FreqTill.Value > 0     //Уведомлять на протяжении события

	start := t.Memo.Scenario.DateStart
	end := t.Memo.Scenario.DateEnd

	now := time.Now()  //now basis time
	next := time.Now() //temporary variable for calculations
	var id int64
	id = 0
	timemap := make(map[int64]time.Time) // [map_id]time
	timecount := make(map[int64]int)     // [map_id]count

	if end.Before(now) && !after {
		return errors.New(fmt.Sprint("Event <Group:", t.Group, ",ID:", t.ID, "> outdated:", end.String()))
	}

	// Remind before event, starts and ends before current date
	if before && !start.Before(now) && !end.Before(now) {
		g = Gran(t.Memo.Scenario.FreqBefore.Granula)

		if t.Memo.Scenario.FreqBefore.Count <= 0 {
			c = 1
		} else {
			c = t.Memo.Scenario.FreqBefore.Count
		}

		// calc next reminding dates
		next = start.Add(time.Second * time.Duration(-1*t.Memo.Scenario.FreqBefore.Value*g*60))
		for !next.After(start) {
			if next.Before(now) {
				next = next.Add(time.Second * time.Duration(c*g*60))
				continue
			}
			timemap[id] = next
			GlobalTimeMap[int64(t.Memo.ID)] = timemap
			timecount[id] = c
			GlobalTimeCount[int64(t.Memo.ID)] = timecount
			id++
			next = next.Add(time.Second * time.Duration(c*g*60)) // use c as value for BEFORE situation
		}
	}

	// Remind till event goes, if event not ended now
	if till && !end.Before(now) {
		g = Gran(t.Memo.Scenario.FreqTill.Granula)

		if t.Memo.Scenario.FreqTill.Count == 0 {
			c = 1
		} else if t.Memo.Scenario.FreqTill.Count < 0 {
			c = 65535 //infinite repeat (until other event cause end)
		} else {
			c = t.Memo.Scenario.FreqTill.Count
		}

		next = start.Add(time.Second * time.Duration(t.Memo.Scenario.FreqTill.Value*g*60))
		actual := !next.After(end)
		next = start
		/*if !start.After(now) {
			next = NextTime(start, t.Memo.Scenario.FreqTill)
		}*/

		for actual && c > 0 {
			if next.Before(now) {
				next = next.Add(time.Second * time.Duration(t.Memo.Scenario.FreqTill.Value*g*60))
				c--
				continue
			}
			timemap[id] = next
			GlobalTimeMap[int64(t.Memo.ID)] = timemap
			timecount[id] = c
			GlobalTimeCount[int64(t.Memo.ID)] = timecount
			id++
			c--
			next = next.Add(time.Second * time.Duration(t.Memo.Scenario.FreqTill.Value*g*60))
			actual = !next.After(end)
		}
	}

	// Remind after event was ended, if last counted remind time is not come (until now)
	if after && !end.Add(time.Second*time.Duration(t.Memo.Scenario.FreqAfter.Value*Gran(t.Memo.Scenario.FreqAfter.Granula)*60*t.Memo.Scenario.FreqAfter.Count)).Before(now) {
		g = Gran(t.Memo.Scenario.FreqAfter.Granula)

		if t.Memo.Scenario.FreqAfter.Count == 0 {
			c = 1
		} else if t.Memo.Scenario.FreqAfter.Count < 0 {
			c = 65535 //infinite repeat (until other event cause end)
		} else {
			c = t.Memo.Scenario.FreqAfter.Count
		}

		next = end.Add(time.Second * time.Duration(t.Memo.Scenario.FreqAfter.Value*g*60))
		for i := 0; i <= c; i++ {
			if next.Before(now) {
				next = next.Add(time.Second * time.Duration(t.Memo.Scenario.FreqAfter.Value*g*60))
				continue
			}
			timemap[id] = next
			GlobalTimeMap[int64(t.Memo.ID)] = timemap
			timecount[id] = c
			GlobalTimeCount[int64(t.Memo.ID)] = timecount
			id++
			next = next.Add(time.Second * time.Duration(t.Memo.Scenario.FreqAfter.Value*g*60))
		}

	}

	return err
}

//NextTime calculate nearest to NOW time point by granula in seconds
//it saves granula shift for time intervals
func NextTime(point time.Time, freq Freq) time.Time {
	now := time.Now()
	start := point
	g := Gran(freq.Granula)
	gransec := g * freq.Value * 60 // Granula in seconds

	// point in PAST
	if start.Before(now) {
		for start.Before(now) {
			start = start.Add(time.Second * time.Duration(gransec))
		}
		return start
	}

	// point in FUTURE

	for !start.Before(now) {
		start = start.Add(time.Second * time.Duration(gransec*(-1)))
	}

	return start
}

//Gran returns actual value for time granula in minutes
func Gran(granula string) int {
	g := 1
	upG := strings.ToUpper(granula)

	if strings.EqualFold(upG, "D") {
		g = 1440 //day
	} else if strings.EqualFold(upG, "H") {
		g = 60 //hour
	} else if strings.EqualFold(upG, "M") {
		g = 1 //minute
	}

	return g
}
