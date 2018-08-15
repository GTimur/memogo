package memogo

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

/*
	log.go - logfile writer
*/

// LOGFILE - constant path to log file
const LOGFILE = "./logfile.log"

// LOGFILEDAT - constant path to log file service info
const LOGFILEDAT = "./logfile.json"

// LogFile - logfile struct
// if Stale flag is set - data will be rewrited after TTLDays
type LogFile struct {
	Filename   string // path to file
	CleanStale bool   // if Stale flag is set - data will be rewrited after TTLDays
	TTLDays    int    // TTL in days if Stale is set
}

// LogDatFile - struct for control-file
type LogDatFile struct {
	Filename string
	Date     time.Time //Date of creation
}

// Add - writes line to logfile
func (l *LogFile) Add(line string) (err error) {
	prefix := time.Now().Format("2006-01-02 15:04:05")
	file, err := os.OpenFile(l.Filename, os.O_APPEND, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("LogFile Add: %v\n", err)
		return err
	}

	_, err = file.WriteString(prefix + "\t" + line + "\n")
	if err != nil {
		log.Printf("LogFile WriteString: %v\n", err)
		return err
	}

	return err
}

// ClearByTTL - clear logfile if TTL expired and update datfile
func (l *LogFile) ClearByTTL(dat LogDatFile) (err error) {
	// if no cleanstale needed - finish here
	if !l.CleanStale {
		return nil
	}
	var date time.Time
	date, err = dat.ReadLogDate()
	if err != nil {
		log.Println("ReadLogDate:", err)
		return err
	}

	// if TTL not expired exit
	if (time.Since(date).Hours() / 24.0) <= float64(l.TTLDays) {
		return err
	}

	l.Add("LogFile cleared cause of TTL, previous DATE:" + date.Format("2006-01-02 15:04:05"))
	GlobalLogDatFile.Date = time.Now()

	err = GlobalLogDatFile.MakeNewFile()
	if err != nil {
		return err
	}

	err = GlobalLogDatFile.WriteJSON()
	if err != nil {
		return err
	}

	return err
}

// Init - new logfile, new logdatefile with date
// Create files if files not exists, sets GlobalLogDatFile.date
func (l *LogFile) Init() (err error) {
	//str := fmt.Sprintf("%02d-%02d-%04dT00:00:01+00:00", date.Day(), date.Month(), date.Year())
	// rewrite log-file
	err = l.MakeFile()
	if err != nil {
		return err
	}
	// rewrite logdatfile
	err = GlobalLogDatFile.MakeFile()
	if err != nil {
		return err
	}

	err = GlobalLogDatFile.WriteJSON()
	if err != nil {
		log.Printf("InitLog error: %v\n", err)
		return err
	}

	l.Add("LogFile initialization completed.")

	return err
}

// ReadLogDate - read date from file
func (l *LogDatFile) ReadLogDate() (date time.Time, err error) {
	err = l.ReadJSON()
	return l.Date, err
}

//writeJSON - file must be exist and accessible
func (l *LogDatFile) writeJSON() (err error) {
	file, err := os.OpenFile(l.Filename, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("LogDatFile: error writing JSON file: %v\n", err)
		return err
	}

	//Готовим данные JSON (конвертируем в экспортируемый вид)
	jsonDat := l

	// пишем в файл
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(&jsonDat)
	if err != nil {
		log.Printf("LogDatFile: JSON encoder error: %v", err)
	}
	return err
}

// WriteJSON - writes struct to JSON file
// filepath is set in LOGFILEDAT constant
func (l *LogDatFile) WriteJSON() (err error) {
	err = l.writeJSON()
	return
}

//Read from Json-config file
func (l *LogDatFile) readJSON() (err error) {
	file, err := os.Open(l.Filename)
	defer file.Close()
	if err != nil {
		log.Printf("LogDatFile: error reading JSON file: %v\n", err)
		return err
	}

	//Готовим для импорта структуру JSON
	var jsonDat LogDatFile

	//Читаем файл JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonDat)
	if err != nil {
		log.Printf("LogDatFile: JSON decoder error: %v", err)
		return err
	}

	*l = jsonDat

	return err
}

// ReadJSON - reads JSON file into struct
func (l *LogDatFile) ReadJSON() (err error) {
	err = l.readJSON()
	return
}

// MakeFile - creates file if it not exists
func (l *LogDatFile) MakeFile() (err error) {
	if _, err = os.Stat(l.Filename); err == nil {
		//File exist and will not be rewrited
		//fmt.Println(os.IsExist(err),err)
		return err
	}
	file, err := os.Create(l.Filename)
	defer file.Close()
	return err
}

// MakeNewFile - creates file
func (l *LogDatFile) MakeNewFile() (err error) {
	file, err := os.Create(l.Filename)
	defer file.Close()
	return err
}

// MakeFile - creates file if it not exists
func (l *LogFile) MakeFile() (err error) {
	if _, err = os.Stat(l.Filename); err == nil {
		//File exist and will not be rewrited
		//fmt.Println(os.IsExist(err),err)
		return err
	}
	file, err := os.Create(l.Filename)
	defer file.Close()
	return err
}

// IsExist - check if file exist
func (l *LogFile) IsExist() (exist bool) {
	if _, err := os.Stat(l.Filename); os.IsNotExist(err) {
		return false //file not exist
	}
	return true
}

// IsExist - check if file exist
func (l *LogDatFile) IsExist() (exist bool) {
	if _, err := os.Stat(l.Filename); os.IsNotExist(err) {
		return false //file not exist
	}
	return true
}

/*
// makeFile - creates/rewrites log file
func (l *LogFile) makeFile() (err error) {
	file, err := os.Create(l.Filename)
	defer file.Close()
	return err
}

// makeDatFile - creates|rewrites log dat-file
func (l *LogFile) makeDatFile() (err error) {
	file, err := os.Create(LOGFILEDAT)
	defer file.Close()
	return err
}

*/
