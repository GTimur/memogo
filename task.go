package memogo

/*
task.go - memo tasks realization
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// REMIND SCENARIO
// REPEAT_UNTIL_EVENT_START = Remind in Int1 before event DateStart and repeat remind every Int2 until event DateStart
// REPEAT_UNTIL_EVENT_END = Remind from DateStart to DateEnd with repeat inteval=Int2 (minutes)

// Memo definiton
type Memo struct {
	ID       int       //ID
	Date     time.Time //Date of creation
	Scenario Remind    //Notification scheduling scenario options
	Subj     string    //visible memo subject (header)
	Memo     []string  //memo body
	Mails    []string  //recepitients e-mails list
}

// Remind options
type Remind struct {
	DateStart  time.Time //Date start of event
	DateEnd    time.Time //Date end of event
	FreqBefore Freq      //Frequency of notification BEFORE event
	FreqTill   Freq      //Frequency of notification UNTIL event
	FreqAfter  Freq      //Frequency of notification AFTER event
}

//Freq - frequency of repeating notification
type Freq struct {
	Granula string // Time granula: mhd (minutes,hours,days)
	Value   int    // Time value
	Count   int    // used for AFTER options only
}

//Task defenition
type Task struct {
	ID    int
	Group string // Group tasks by each folder
	Memo  Memo
}

//ReadJSON reads JSON file into Memo struct
func (m *Memo) ReadJSON(filename string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Printf("Memo.ReadJSON: error reading JSON-file: %v\n", err)
		return err
	}

	//Read JSON file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&m)
	if err != nil {
		log.Printf("JSON decoder error: %v", err)
		return err
	}
	return err
}

//WriteJSON - writes Memo struct into JSON file
func (m *Memo) WriteJSON(filename string) error {
	var file *os.File
	err := MakeFile(filename)
	if err != nil {
		log.Printf("Memo.MakeFile: %v\n", err)
		return err
	}

	file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("Memo.WriteJSON: error writing JSON file: %v\n", err)
		return err
	}

	// write to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(&m)
	if err != nil {
		log.Printf("JSON encoder error: %v", err)
	}
	return err
}

// MakeFile - creates file if it not exists
func MakeFile(filename string) (err error) {
	if _, err = os.Stat(filename); err == nil {
		//File exist and will not be rewrited
		//fmt.Println(os.IsExist(err),err)
		return err
	}
	file, err := os.Create(filename)
	defer file.Close()
	return err
}

// TestJSON -
func TestJSON() error {
	var m Memo
	var m1 Memo
	var r Remind

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	r.DateStart = time.Date(2018, time.May, 14, 10, 0, 0, 0, loc)
	r.DateEnd = time.Date(2018, time.August, 16, 23, 59, 59, 0, loc)
	r.FreqBefore = Freq{Granula: "m", Value: -1, Count: -1}
	r.FreqTill = Freq{Granula: "m", Value: 1, Count: -1}
	r.FreqAfter = Freq{Granula: "m", Value: 1, Count: 10} // 10 times after event is end, repeat every 1 minute

	m.ID = 100
	m.Date = time.Now()
	m.Scenario = r
	m.Subj = "Срок действия сертификата Thawte для портала ibank2"
	m.Memo = []string{
		"Истекает срок действия сертификата для",
		"web-сервера ibank.ymkbank.ru",
		"контактная информация: +7(861)2100553.",
	}
	m.Mails = []string{"gtg@ymkbank.ru", "support@ymkbank.ru"}

	err = m.WriteJSON("task01.json")
	if err != nil {
		return err
	}

	err = m1.ReadJSON("task01.json")
	if err != nil {
		return err
	}
	fmt.Println("m1:", m1)
	return err
}
