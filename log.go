package memogo

import (
	"os"
	"time"
)

/*
	log.go - logfile writer
*/

// LogFile - logfile struct
// if Stale flag is set - data will be rewrited after TTLDays
type LogFile struct {
	Filename   string    // path to file
	CleanStale bool      // if Stale flag is set - data will be rewrited after TTLDays
	TTLDays    int       // TTL in days if Stale is set
	Date       time.Time //Date of creation
}

// Writeln - writes line to logfile
func (l *LogFile) Writeln() (err error) {

	return err
}

// MakeFile - creates config file if it not exists
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

/*
// fileDate - return file creation date
func (l *LogFile) fileDate() (cdate time.Time, err error) {
	file, err := os.Stat(l.Filename)
	defer file.Close()
	if err != nil {
		log.Printf("LogFile: error reading log-file: %v\n", err)
		return cdate, err
	}

	return file.cdate, err
}
*/
