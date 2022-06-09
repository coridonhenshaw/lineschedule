package imp

import (
	"database/sql"
	"fmt"
	"lineschedule/db"
	"log"
	"strconv"
)

type StopTimesCallbackStruct struct {
	LineMap  *map[string]int
	TIDIndex int
	ATIndex  int
	DTIndex  int
	SIIndex  int
	SSIndex  int

	InsertStatement *sql.Stmt

	N int
}

func (o *StopTimesCallbackStruct) SetLinemap(Map *map[string]int) {
	o.LineMap = Map
	o.TIDIndex = (*o.LineMap)["trip_id"]
	o.ATIndex = (*o.LineMap)["arrival_time"]
	o.DTIndex = (*o.LineMap)["departure_time"]
	o.SIIndex = (*o.LineMap)["stop_id"]
	o.SSIndex = (*o.LineMap)["stop_sequence"]

	var err error
	var SQL = "INSERT INTO stoptimes (trip_id, arrival_time, departure_time, stop_id, stop_sequence) VALUES (?, ?, ?, ?, ?)"
	o.InsertStatement, err = db.Con.Prepare(SQL)
	if err != nil {
		log.Panic(err)
	}
}

func ParseInt(value string) int {
	Result, err := strconv.Atoi(value)
	if err != nil {
		log.Panic(err)
	}
	return Result
}

func CookTime(Time string) int {

	//	Time := strings.Split(TimeString, ":")

	Hours := ParseInt(Time[0:2])
	Minutes := ParseInt(Time[3:5])
	Seconds := ParseInt(Time[6:8])

	var Final int = int(Hours*60*60 + Minutes*60 + Seconds)

	return Final
}

func (o *StopTimesCallbackStruct) Line(Line *[]string) {

	var trip_id = (*Line)[o.TIDIndex]
	var arrival_time = (*Line)[o.ATIndex]
	var departure_time = (*Line)[o.DTIndex]
	var stop_id = (*Line)[o.SIIndex]
	var stop_sequence = (*Line)[o.SSIndex]

	var arrival_time_cooked = CookTime(arrival_time)
	var departure_time_cooked = CookTime(departure_time)

	//	var SQL = "INSERT INTO stoptimes (trip_id, arrival_time, departure_time, stop_id, stop_sequence) VALUES (?, ?, ?, ?, ?)"
	//	SQLExec(SQL, trip_id, arrival_time_cooked, departure_time_cooked, stop_id, stop_sequence)

	o.InsertStatement.Exec(trip_id, arrival_time_cooked, departure_time_cooked, stop_id, stop_sequence)

	if o.N%10000 == 0 {
		fmt.Printf("%v\r", o.N)
	}
	o.N++
}
