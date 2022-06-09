package imp

import "lineschedule/db"

type CalendarDatesCallbackStruct struct {
	LineMap *map[string]int
}

func (o *CalendarDatesCallbackStruct) SetLinemap(Map *map[string]int) {
	o.LineMap = Map
}

func (o *CalendarDatesCallbackStruct) Line(Line *[]string) {
	var service_id = (*Line)[(*o.LineMap)["service_id"]]
	var date = (*Line)[(*o.LineMap)["date"]]
	var exception_type = (*Line)[(*o.LineMap)["exception_type"]]

	var SQL = "INSERT INTO calendar_dates (service_id, date, exception_type) VALUES (?, ?, ?)"
	db.SQLExec(SQL, service_id, date, exception_type)
}
