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

// TasksGlobal - Tasks array
var TasksGlobal []Task

// Rebuild - rebuilds Tasks array
func TasksRebuild() error {
	var task Task

	// clean Tasks array
	TasksGlobal = TasksGlobal[0:0]
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
		task.ID = i
		task.Group = filepath.Dir(strings.Replace(k, GlobalConfig.Root, "", -1)) //get name of folder as name of group
		err := task.Memo.ReadJSON(k)
		if err != nil {
			log.Fatalf("TaskMgr Rebuild error: %v", err)
			return err
		}
		i++
		TasksGlobal = append(TasksGlobal, task)
	}

	fmt.Println(TasksGlobal)

	return err
}
