package memogo

import (
	"errors"
	"fmt"
	"log"
	"time"
)

/*not used*/
/*
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
*/

// GetQueueEvent - takes event from queue if it ready to be processed
func GetQueueEvent(queue Queue) (q Queue, err error) {
	q = q[:0]
	if len(queue) == 0 {
		return q, errors.New("Queue contains no event")
	}
	alone := false
	ts := 1 * 60 // +-1 minutes time shift (+1 -1 = 2 min)
	now := time.Now()
	start := now.Add(time.Second * time.Duration(ts*-1))
	end := now.Add(time.Second * time.Duration(ts))
	event1, event2 := now, now

	// Check all events in queue per pairs
	for i := 0; i < queue.Len(); i++ {
		if len(queue) >= i+1 {
			event1 = queue[i].Plan.Run
			event2 = queue[i+1].Plan.Run
			alone = false
		} else { // if no pair exist
			event1 = queue[i].Plan.Run
			alone = true
		}

		if event1.Before(start) || event1.After(end) {
			return q, err
		}

		// If skipped event1 time (T) < event2-event1 then it will be processed
		//time.Since(start).Minutes() < event2.Sub(event1).Minutes()
		GlobalConfig.LogFile.Add(fmt.Sprint("INFO: Time has come for event1: <", event1, ">, value for event2:<", event2, ">."))
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

// MemoSvc - execute notification procedure
func MemoSvc(queue Queue) (err error) {
	var q Queue

	q, err = GetQueueEvent(queue)
	if err != nil {
		if len(q) != 0 {
			return err
		}
		return nil
	}

	auth := EmailCredentials{Username: GlobalConfig.SMTPSrv.Account,
		Password: GlobalConfig.SMTPSrv.Password,
		Server:   GlobalConfig.SMTPSrv.Addr,
		Port:     GlobalConfig.SMTPSrv.Port,
		From:     GlobalConfig.SMTPSrv.From,
		FromName: GlobalConfig.SMTPSrv.FromName,
		UseTLS:   GlobalConfig.SMTPSrv.UseTLS,
	}

	for _, evt := range q {
		var memo Memo
		memo, err = GlobalTasks.GetMemo(evt.MemoID)
		if err != nil {
			log.Println("MemoSvc: GetMemo error, Memo.ID ", evt.MemoID, ":\n", err)
			return err
		}
		memoStr := ""
		for _, r := range memo.Memo {
			memoStr += r + "<br>"
		}

		//create new message (subj, msgbody)
		msg := NewHTMLMessage(memo.Subj, memoStr)

		// Collect all email recipients
		msg.Body += "<br><br>Направлено:<br>"

		if len(memo.Mails) == 0 {
			log.Println("MemoSvc: No emails found for Memo ID=", memo.ID)
			return errors.New(fmt.Sprintln("MemoSvc: No emails found for Memo ID=", memo.ID))
		}
		for _, m := range memo.Mails {
			msg.To = append(msg.To, m)
			msg.Body += m + "<br>"
		}
		// Remove doublecates if exist
		msg.To = Dedup(msg.To)

		// Sending memo
		if err := SendEmailMsg(auth, msg); err != nil {
			log.Println("MemoSvc: Error sending mail for \"", msg.To, "\":", err)
			return err
		}
		err = queue.RemoveEvt(evt)
		GlobalConfig.LogFile.Add(fmt.Sprint("Message sent. <", evt.MemoID, ",", evt.Plan.Run, "> Queue size:", len(queue)))
		if err != nil {
			log.Println("MemoSvc: Error r", err)
			return err
		}
	}
	return err
}

// Dedup - String slice deduplication (not saves the string sorting)
func Dedup(slice []string) []string {
	checked := map[string]bool{}

	//Сохраним отображение без повторяющихся элементов
	for i := range slice {
		checked[slice[i]] = true
	}

	//Перенесем отображение в результат
	result := []string{}
	for key := range checked {
		result = append(result, key)
	}

	return result
}
