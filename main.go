package main

import (
	"fmt"
	"log"
	"memogo"
	"time"
)

// REMIND SCENARIO
// REPEAT_UNTIL_EVENT_START = Remind in Int1 before event DateStart and repeat remind every Int2 until event DateStart
// REPEAT_UNTIL_EVENT_END = Remind from DateStart to DateEnd with repeat inteval=Int2 (minutes)

// Remind options
type Remind struct {
	DateStart time.Time //Date start of event
	DateEnd   time.Time //Date end of event
	Int1      int       //Frequency of reminder messages (minutes)
	Int2      int       //Frequency of reminder messages (minutes)
}

// Memo definiton
type Memo struct {
	ID       int       //ID
	Date     time.Time //Date of creation
	Reminder Remind    //Remind scenario options
	Subj     string    //visible memo subject
	Memo     []string  //memo body
	Mails    []string  //emails list
}

var globalconfig memogo.Config

func main() {
	memo := Memo{}
	rem := Remind{}
	//smtp :=
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

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	rem.DateStart = time.Date(2018, time.May, 11, 10, 0, 0, 0, loc)
	rem.DateEnd = time.Date(2018, time.August, 16, 23, 59, 59, 0, loc)
	rem.Int1 = 5 //in * minutes
	rem.Int2 = 1 //every * minutes

	memo.ID = 100
	memo.Date = time.Now()
	memo.Subj = "Сертификат thawte ibank.ymkbank.ru истекает 17.08.2018"
	memo.Memo = []string{"Истекает срок действия сертификата для", "web-сервера ibank.ymkbank.ru"}
	memo.Mails = []string{"gtg@ymkbank.ru"}
	memo.Reminder = rem

	fmt.Println(memo)

	files, err = memogo.FindAllFiles(globalconfig.Root, []string{"*.*"})
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
}
