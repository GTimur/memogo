package memogo

import "time"

/*
  Just global variables here...
*/

// GlobalConfig - global configuration (jsonconf)
var GlobalConfig Config

// QueueGlobal - global Queue (plan)
var GlobalQueue Queue

// GlobalTimeMap - store map of runtime for every memo (runmap)
var (
	GlobalTimeMap   map[int64]map[int64]time.Time
	GlobalTimeCount map[int64]map[int64]int
)

// GlobalTasks - Tasks array (taskmgr)
var GlobalTasks []Task
