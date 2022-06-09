package imp

import (
	"lineschedule/db"
	"log"
	"time"
)

type CalendarCallbackStruct struct {
	LineMap *map[string]int
}

func (o *CalendarCallbackStruct) SetLinemap(Map *map[string]int) {
	o.LineMap = Map
}

func (o *CalendarCallbackStruct) Line(Line *[]string) {
	var service_id = (*Line)[(*o.LineMap)["service_id"]]
	var start_date = (*Line)[(*o.LineMap)["start_date"]]
	var end_date = (*Line)[(*o.LineMap)["end_date"]]
	var monday = (*Line)[(*o.LineMap)["monday"]]
	var tuesday = (*Line)[(*o.LineMap)["tuesday"]]
	var wednesday = (*Line)[(*o.LineMap)["wednesday"]]
	var thursday = (*Line)[(*o.LineMap)["thursday"]]
	var friday = (*Line)[(*o.LineMap)["friday"]]
	var saturday = (*Line)[(*o.LineMap)["saturday"]]
	var sunday = (*Line)[(*o.LineMap)["sunday"]]

	var err error
	sdd, err := time.Parse("20060102", start_date)
	if err != nil {
		log.Panic(err)
	}
	sd := sdd.Unix()

	edd, err := time.Parse("20060102", end_date)
	if err != nil {
		log.Panic(err)
	}
	ed := edd.Unix()

	var SQL = "INSERT INTO calendar (service_id, start_date, end_date, mon, tue, wed, thu, fri, sat, sun) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	db.SQLExec(SQL, service_id, sd, ed, monday, tuesday, wednesday, thursday, friday, saturday, sunday)
}
