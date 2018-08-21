package memogo

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

/*
taskmgr.go - tasks manager realization
*/

// TasksReload - rebuilds Tasks array
func TasksReload() error {
	// clean Tasks array
	GlobalTasks = GlobalTasks[:0]
	mapID := make(map[int64]bool)   // map[memoid]bool
	fixID := make(map[int64]string) // map[memoid]path to fix duplicates
	var maxID int64
	//	groups := make(map[int]string)

	// collect all groups and all files
	var files map[string]string

	files, err := FindAllFiles(GlobalConfig.Root, []string{"*.*"})
	if err != nil {
		log.Fatalf("TaskMgr Rebuild: FindAllFiles error: %v", err)
		return err
	}

	i := 0
	for k := range files {
		var task Task
		task.ID = i
		task.Group = filepath.Dir(strings.Replace(k, GlobalConfig.Root, "", -1)) //get name of folder as name of group
		err := task.Memo.ReadJSON(k)
		if err != nil {
			GlobalConfig.LogFile.Add(fmt.Sprintf("ERROR: Task file=<%s>: %v", k, err))
			//return err
			continue
		}
		// if Draft parameter set - ignore this memo
		if task.Memo.Draft {
			continue
		}
		// check for uniq Memo.ID (must no have duplicates)
		if _, ok := mapID[task.Memo.ID]; ok {
			GlobalConfig.LogFile.Add(fmt.Sprint("TasksReload: found duplcate Memo.ID: Memo.ID=<", task.Memo.ID, "> Group=<", task.Group, "> File=<", k, ">"))
			fixID[task.Memo.ID] = k // map[memoid]path
			i++
			continue
		}
		if maxID < task.Memo.ID {
			maxID = task.Memo.ID
		}
		mapID[task.Memo.ID] = true
		i++
		GlobalTasks = append(GlobalTasks, task)
	}

	// FIX duplicates Memo.ID
	for _, v := range fixID {
		var task Task
		task.ID = i
		task.Group = filepath.Dir(strings.Replace(v, GlobalConfig.Root, "", -1)) //get name of folder as name of group
		err := task.Memo.ReadJSON(v)
		if err != nil {
			log.Fatalf("TaskMgr Rebuild error: %v", err)
			return err
		}
		task.Memo.ID = maxID + 1
		maxID++
		err = task.Memo.WriteJSON(v)
		if err != nil {
			log.Fatalf("TaskMgr Rebuild error: %v", err)
			return err
		}
		GlobalConfig.LogFile.Add(fmt.Sprint("TasksReload: fixed duplcate Memo.ID. New data: Memo.ID=<", task.Memo.ID, "> Group=<", task.Group, "> File=<", v, ">"))
		i++
		GlobalTasks = append(GlobalTasks, task)
	}

	return err
}
