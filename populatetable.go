package main

import (
	"fmt"
	"lineschedule/db"
	"log"
	"strings"
	"time"
)

type RouteStruct struct {
	Route     string
	Direction int
	Stops     []string
}

func GetTable(ServiceDate string, RouteBlock []RouteStruct, Connections []ConnectionStruct, StopNames []StopNamesStruct) *TableStruct {
	var ts TableStruct

	ridmap := make(map[string]string)
	for _, v := range RouteBlock {
		ridmap[v.Route] = GetRouteID(v.Route)
	}

	for i := range Connections {
		Connections[i].fromrouteid = ridmap[Connections[i].FromRoute]
		Connections[i].torouteid = ridmap[Connections[i].ToRoute]
	}

	NameOverride := make(map[string]string)

	for _, v := range StopNames {
		NameOverride[v.Stop_id] = v.Name
	}

	for _, x := range RouteBlock {
		route_id := ridmap[x.Route]

		Date, err := time.Parse("2006-01-02", ServiceDate)
		if err != nil {
			log.Panic(err)
		}

		Trips := GetTripsForDate(Date, route_id, x.Direction)

		for _, v := range Trips {
			var r0 []EntryStruct
			Times := GetTimesForTrip(v.trip_id, x.Stops)
			for _, b := range Times {

				var q ConnectionStruct
				var m int
				for _, n := range Connections {
					m = n.Test(route_id, b.stop_id)
					if m != 0 {
						q = n
						break
					}
				}

				if m == -1 { //to connection
					es := EntryStruct{stop_id: q.FromStop, stop_name: ts.StopNameMap[q.FromStop], arrival_time: b.arrival_time, departure_time: b.departure_time, Valid: true, Connection: true, Visible: false}
					r0 = append(r0, es)
					es = EntryStruct{stop_id: q.ToStop, stop_name: b.stop_name, arrival_time: b.arrival_time, departure_time: b.departure_time, Valid: true, Connection: true, Visible: true}
					r0 = append(r0, es)
				} else if m == 1 { //from connection
					es := EntryStruct{stop_id: q.FromStop, stop_name: b.stop_name, arrival_time: b.arrival_time, departure_time: b.departure_time, Valid: true, Connection: true, Visible: true}
					r0 = append(r0, es)
					es = EntryStruct{stop_id: q.ToStop, stop_name: b.stop_name, arrival_time: b.arrival_time, departure_time: b.departure_time, Valid: true, Connection: true, Visible: false}
					r0 = append(r0, es)
				} else {
					es := EntryStruct{stop_id: b.stop_id, stop_name: b.stop_name, arrival_time: b.arrival_time, departure_time: b.departure_time, Valid: true, Visible: true}
					r0 = append(r0, es)
				}

			}
			for i := range r0 {
				d, v := NameOverride[r0[i].stop_id]
				if v {
					r0[i].stop_id = d
				}
			}
			ts.AddRow(v.headsign, r0)
		}

	}

	ts.Finalize()

	return &ts

}

type ConnectionStruct struct {
	FromRoute string
	FromStop  string
	ToRoute   string
	ToStop    string

	fromrouteid string
	torouteid   string
}

func (o *ConnectionStruct) Test(route_id string, stop_id string) int {
	if len(o.fromrouteid) == 0 {
		log.Panic()
	}
	if len(o.torouteid) == 0 {
		log.Panic()
	}

	if route_id == o.fromrouteid && stop_id == o.FromStop {
		return 1
	}
	if route_id == o.torouteid && stop_id == o.ToStop {
		return -1
	}

	return 0
}

func GetSchedule(Route string, ServiceDate string, Direction int, Stops []string) *TableStruct {
	route_id := GetRouteID(Route)

	Date, err := time.Parse("2006-01-02", ServiceDate)
	if err != nil {
		log.Panic(err)
	}

	Trips := GetTripsForDate(Date, route_id, Direction)

	var ts TableStruct

	for _, v := range Trips {
		var r0 []EntryStruct
		Times := GetTimesForTrip(v.trip_id, Stops)
		for _, b := range Times {
			es := EntryStruct{stop_id: b.stop_id, stop_name: b.stop_name, arrival_time: b.arrival_time, departure_time: b.departure_time, Valid: true}
			r0 = append(r0, es)
		}
		ts.AddRow(v.headsign, r0)
	}

	ts.Finalize()

	return &ts
}

func GetRouteID(Route string) string {
	if len(Route) == 0 {
		log.Panic("No route name provided.")
	}

	SQL := `SELECT route_id FROM routes WHERE LTRIM(route_short_name, "0") == ?`

	var route_id string

	err := db.Con.QueryRow(SQL, Route).Scan(&route_id)
	if err != nil {
		log.Panic(err)
	}

	return route_id
}

func GetDirections(route_id string) []string {
	var Directions []string

	SQL := `SELECT direction, direction_id FROM directions WHERE route_id == ?`

	Rows, err := db.Con.Query(SQL, route_id)
	if err != nil {
		log.Panic(err)
	}
	defer Rows.Close()

	for Rows.Next() {
		var direction string
		if err := Rows.Scan(&direction); err != nil {
			log.Panic(err)
		}
		Directions = append(Directions, direction)
	}
	if err = Rows.Err(); err != nil {
		log.Panic(err)
	}

	return Directions
}

type TripStruct struct {
	trip_id   string
	headsign  string
	direction int
}

func GetTripsForDate(Date time.Time, route_id string, direction int) []TripStruct {
	var Trips []TripStruct

	Weekday := Date.Weekday()
	WeekdayCooked := strings.ToLower(Weekday.String())[0:3]

	SQL := `SELECT trip_id, headsign, direction FROM trips
	 		LEFT JOIN calendar ON trips.service_id == calendar.service_id
	 		WHERE %v=1 AND route_id = ? AND direction = ? AND start_date < ? AND end_date > ?
			UNION
			SELECT trip_id, headsign, direction FROM trips LEFT JOIN calendar_dates ON trips.service_id == calendar_dates.service_id
			WHERE date = ? AND route_id = ? AND direction = ? AND exception_type == 1
			EXCEPT
			SELECT trip_id, headsign, direction FROM trips LEFT JOIN calendar_dates ON trips.service_id == calendar_dates.service_id
			WHERE date = ? AND route_id = ? AND direction = ? AND exception_type == 2`

	SQL = fmt.Sprintf(SQL, WeekdayCooked)

	var CookedDate = Date.Unix()
	DateCooked := Date.Format("20060102")
	Rows, err := db.Con.Query(SQL, route_id, direction, CookedDate, CookedDate, DateCooked, route_id, direction, DateCooked, route_id, direction)
	if err != nil {
		log.Panic(err)
	}
	defer Rows.Close()

	var Trip TripStruct

	for Rows.Next() {
		if err := Rows.Scan(&Trip.trip_id, &Trip.headsign, &Trip.direction); err != nil {
			log.Panic(err)
		}
		Trips = append(Trips, Trip)
	}
	if err = Rows.Err(); err != nil {
		log.Panic(err)
	}

	return Trips
}

type TimesStruct struct {
	arrival_time   int
	departure_time int
	stop_name      string
	stop_id        string
}

func GetTimesForTrip(trip_id string, Stops []string) []TimesStruct {

	var TimesList []TimesStruct

	var StopSQL string
	if len(Stops) > 0 {
		StopSQL = " AND stops.stop_ID IN (" + strings.Join(Stops, ", ") + ") "
	}

	SQL := `SELECT arrival_time, departure_time, stop_name, stops.stop_id FROM stoptimes
	 LEFT JOIN stops ON stoptimes.stop_id == stops.stop_id
	 WHERE trip_id = ?` + StopSQL + `ORDER BY stop_sequence`

	Rows, err := db.Con.Query(SQL, trip_id)
	if err != nil {
		log.Panic(err)
	}
	defer Rows.Close()

	var Times TimesStruct

	for Rows.Next() {
		if err := Rows.Scan(&Times.arrival_time, &Times.departure_time, &Times.stop_name, &Times.stop_id); err != nil {
			log.Panic(err)
		}
		TimesList = append(TimesList, Times)
	}
	if err = Rows.Err(); err != nil {
		log.Panic(err)
	}
	return TimesList
}
