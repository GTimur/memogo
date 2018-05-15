package memogo

import (
	"fmt"
	"log"
	"path"
	"strings"
)

/*
taskmgr.go - tasks manager realization
*/

//var Tasks []Task

// Rebuild - collect files
func Rebuild() error {
	//var task Task
	groups := make(map[int]string)

	// collect all groups and all files
	var files map[string]string

	files, err := FindAllFiles(GlobalConfig.Root, []string{"*.*"})
	if err != nil {
		log.Fatalf("TaskMgr Rebuild: FindFiles error: %v", err)
		return err
	}

	i := 0
	for k := range files {
		groups[i] = path.Dir(strings.Replace(k, GlobalConfig.Root, "", -1))
		i++
		fmt.Println(path.Dir(strings.Replace(k, GlobalConfig.Root, "", -1)))
	}

	return err
}
