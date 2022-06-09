package imp

import (
	"database/sql"
	"lineschedule/db"
	"log"
)

type TripCallbackStruct struct {
	LineMap  *map[string]int
	TIDIndex int

	InsertStatement *sql.Stmt
}

func (o *TripCallbackStruct) SetLinemap(Map *map[string]int) {
	o.LineMap = Map

	o.TIDIndex = (*o.LineMap)["trip_id"]

	var err error
	var SQL = "INSERT INTO trips (route_id, service_id, trip_id, headsign, direction) VALUES (?, ?, ?, ?, ?)"
	o.InsertStatement, err = db.Con.Prepare(SQL)
	if err != nil {
		log.Panic(err)
	}
}

func (o *TripCallbackStruct) Line(Line *[]string) {
	var route_id = (*Line)[(*o.LineMap)["route_id"]]
	var service_id = (*Line)[(*o.LineMap)["service_id"]]
	var trip_id = (*Line)[o.TIDIndex]
	var headsign = (*Line)[(*o.LineMap)["trip_headsign"]]
	var direction = (*Line)[(*o.LineMap)["direction_id"]]

	o.InsertStatement.Exec(route_id, service_id, trip_id, headsign, direction)
}
