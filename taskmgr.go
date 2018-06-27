package memogo

import (
	"log"
	"path/filepath"
	"strings"
)

/*
taskmgr.go - tasks manager realization
*/

// TasksReload - rebuilds Tasks array
func TasksReload() error {
	var task Task

	// clean Tasks array
	GlobalTasks = GlobalTasks[0:0]
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
		GlobalTasks = append(GlobalTasks, task)
	}
	//fmt.Println(GlobalTasks)
	return err
}
