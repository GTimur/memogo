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

// Tasks array
var Tasks []Task

// Rebuild - rebuilds Tasks array
func Rebuild() error {
	var task Task
	//	groups := make(map[int]string)

	// collect all groups and all files
	var files map[string]string

	files, err := FindAllFiles(GlobalConfig.Root, []string{"*.*"})
	if err != nil {
		log.Fatalf("TaskMgr Rebuild: FindFiles error: %v", err)
		return err
	}

	i := 0
	for k := range files {
		task.ID = i
		task.Group = filepath.Dir(strings.Replace(k, GlobalConfig.Root, "", -1))
		task.Memo.ReadJSON(k)
		i++
		Tasks = append(Tasks, task)
	}

	fmt.Println(Tasks)

	return err
}
