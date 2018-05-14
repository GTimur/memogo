package main

import (
	"fmt"
	"log"
	"memogo"
)

var globalconfig memogo.Config

func main() {
	globalconfig = memogo.Config{
		Root: "./root/",
		SMTPSrv: memogo.SrvSMTP{
			Addr:     "10.20.20.6",
			Port:     25,
			Account:  "noti",
			Password: "Bank999",
			From:     "noti@ymkbank.ru",
			FromName: "Memo GO",
			UseTLS:   false,
		},
		MgrSrv: memogo.ManagerSrv{
			Addr: "127.0.0.1",
			Port: 8000,
		},
	}

	var files map[string]string

	files, err := memogo.FindAllFiles(globalconfig.Root, []string{"*.*"})
	if err != nil {
		log.Fatalf("Main(): FindFiles error: %v", err)
	}

	fmt.Println("FILES FOUND:", files)

	err = globalconfig.MakeConfig()
	if err != nil {
		panic(err)
	}
	err = globalconfig.WriteJSON()
	if err != nil {
		panic(err)
	}

	err = memogo.TestJSON()
	if err != nil {
		log.Fatalf("TestJSON() error: %v", err)
	}
}
