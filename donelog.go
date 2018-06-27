package memogo

import "time"

/*
  donelog - realization of logfile which contain information about date and time of successful processing(reminding) for every task

  +------------------------------------------
  | Task ID | DD | MM | RRRR | HH | MI | SS |
  +------------------------------------------

*/

// LogRow - record of donelog
type LogRow struct {
	TaskID int
	DD     int
	MM     int
	RRRR   int
	HH     int
	MI     int
	SS     int
}

// DoneLog - array of records
type DoneLog []LogRow

// Row returns LogRow record by time and id
func Row(t time.Time, id int) LogRow {
	var lr LogRow

	lr.TaskID = id
	lr.DD = t.Day()
	lr.MM = int(t.Month())
	lr.RRRR = t.Year()
	lr.HH = t.Hour()
	lr.MI = t.Minute()
	lr.SS = t.Second()

	return lr
}
