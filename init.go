package memogo

import (
	"fmt"
	"log"
	"time"
)

/*
  init.go - initialization
*/

const (
	BANNER  = "MemoGO reminder utility"
	VERSION = "v.0.2.9 by GTG (C) 2018"
)

// Banner - print program banner
func Banner() {
	fmt.Println(BANNER + " " + VERSION)
	fmt.Println("USAGE: Flag -h shows OPTIONS and HELP information.")
}

func InitConfig() {
	err := GlobalConfig.ReadJSON()
	if err != nil {
		log.Fatalln("Error: Config file (", CONFIGFILE, ") not found.\n Program terminated.")
	}
}

func InitLog() {
	// If logfiles pair exist
	var logExist, datExist bool
	var err error

	GlobalLogDatFile.Filename = LOGFILEDAT

	logExist = GlobalConfig.LogFile.IsExist()
	datExist = GlobalLogDatFile.IsExist()

	if !datExist {
		GlobalLogDatFile.Date = time.Now() //.Add(time.Duration(60) * time.Second)
	} else {
		GlobalLogDatFile.Date, err = GlobalLogDatFile.ReadLogDate()
		if err != nil {
			log.Fatalln("InitLog ReadLogDate error:", err)
		}
	}

	// if logfile or datfile not exist - reinit all log files
	if !logExist || !datExist {
		fmt.Println("INFO: Logfile (" + GlobalConfig.LogFile.Filename + ") or control file (" + LOGFILEDAT + ") not found and will be created.")
		// logfile initialization
		err := GlobalConfig.LogFile.Init()
		if err != nil {
			log.Fatalln("LogFile init:", err)
		}
	}

	// ClearByTTL - clear logfile if TTL expired
	err = GlobalConfig.LogFile.ClearByTTL(GlobalLogDatFile)
	if err != nil {
		log.Fatalln("ClearByTTL:", err)
	}

	if !GlobalConfig.LogFile.CleanStale {
		fmt.Println("INFO: Logfile cleaning is turned off (CleanStale=false)")
	}
}

func InitEvents() {
	// read tasks from disk and rebuild GlobalTask
	err := TasksReload()
	if err != nil {
		log.Fatal("InitEvents (Tasks) error:", err)
	}

	// read GlobalTask array and rebuild GlobalTimeMap
	err = BuildTimeMap()
	if err != nil {
		log.Fatal("InitEvents (Map) error:", err)
	}

	// read GlobalTimeMap and build GlobalQueue
	err = GlobalQueue.MakeQueue()
	if err != nil {
		log.Fatal("InitEvents (Queue) error:", err)
	}
}
