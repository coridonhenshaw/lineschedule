package imp

import "lineschedule/db"

type StopCallbackStruct struct {
	LineMap *map[string]int
}

func (o *StopCallbackStruct) SetLinemap(Map *map[string]int) {
	o.LineMap = Map
}

func (o *StopCallbackStruct) Line(Line *[]string) {
	var stop_id = (*Line)[(*o.LineMap)["stop_id"]]
	var stop_name = (*Line)[(*o.LineMap)["stop_name"]]

	var SQL = "INSERT INTO stops (stop_id, stop_name) VALUES (?, ?)"
	db.SQLExec(SQL, stop_id, stop_name)
}
