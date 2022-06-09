package imp

import (
	"encoding/csv"
	"fmt"
	"io"
	"lineschedule/db"
	"log"
	"os"
	"path/filepath"
)

type cboi interface {
	SetLinemap(*map[string]int)
	Line(r *[]string)
}

type CSVFileStruct struct {
	Basepath string
	LineMap  map[string]int
	Callback cboi
}

func (o *CSVFileStruct) Read(filePath string) {

	f, err := os.Open(filepath.Join(o.Basepath, filePath))
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	Line, err := r.Read()
	if err != nil {
		log.Panic(err)
	}

	o.LineMap = make(map[string]int)
	for i, v := range Line {
		o.LineMap[v] = i
	}

	o.Callback.SetLinemap(&o.LineMap)

	for {
		Line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}

		o.Callback.Line(&Line)
	}
}

func Import(Basepath string) {
	db.SQLExec(`PRAGMA synchronous = OFF`)
	db.SQLExec("PRAGMA journal_mode = OFF")
	db.SQLExec("PRAGMA temp_store = MEMORY")
	db.SQLExec("PRAGMA page_size = 32767")
	db.SQLExec("DROP TABLE IF EXISTS routes")
	db.SQLExec("DROP TABLE IF EXISTS trips")
	db.SQLExec("DROP TABLE IF EXISTS stoptimes")
	db.SQLExec("DROP TABLE IF EXISTS stops")
	db.SQLExec("DROP TABLE IF EXISTS calendar")
	db.SQLExec("DROP TABLE IF EXISTS calendar_dates")
	db.SQLExec("DROP INDEX IF EXISTS stoptimesindex")
	db.SQLExec("DROP INDEX IF EXISTS tripsindex")
	db.SQLExec("VACUUM")
	db.SQLExec("BEGIN")
	db.SQLExec("CREATE TABLE IF NOT EXISTS routes (route_id text PRIMARY KEY, route_short_name text NOT NULL, route_name text NOT NULL) ")
	db.SQLExec("CREATE TABLE IF NOT EXISTS trips (route_id text, service_id text, trip_id integer PRIMARY KEY, headsign text NOT NULL, direction integer) ")
	db.SQLExec("CREATE TABLE IF NOT EXISTS stoptimes (trip_id integer, arrival_time integer, departure_time integer, stop_id integer, stop_sequence integer) ")
	db.SQLExec("CREATE TABLE IF NOT EXISTS stops (stop_id integer PRIMARY KEY, stop_name text NOT NULL) ")
	db.SQLExec("CREATE TABLE IF NOT EXISTS calendar (service_id text, start_date integer, end_date integer, mon integer, tue integer, wed integer, thu integer, fri integer, sat integer, sun integer )")
	db.SQLExec("CREATE TABLE IF NOT EXISTS calendar_dates (service_id TEXT, date TEXT, exception_type INTEGER) ")

	var CSVFS CSVFileStruct
	CSVFS.Basepath = Basepath

	fmt.Println("Calendar")
	CSVFS.Callback = new(CalendarCallbackStruct)
	CSVFS.Read("calendar.txt")

	fmt.Println("Calendar Dates")
	CSVFS.Callback = new(CalendarDatesCallbackStruct)
	CSVFS.Read("calendar_dates.txt")

	fmt.Println("Routes")
	CSVFS.Callback = new(RouteCallbackStruct)
	CSVFS.Read("routes.txt")

	fmt.Println("Trips")
	CSVFS.Callback = new(TripCallbackStruct)
	CSVFS.Read("trips.txt")

	fmt.Println("Stops")
	CSVFS.Callback = new(StopCallbackStruct)
	CSVFS.Read("stops.txt")

	fmt.Println("Stop times")
	CSVFS.Callback = new(StopTimesCallbackStruct)
	CSVFS.Read("stop_times.txt")

	db.SQLExec("CREATE INDEX IF NOT EXISTS stoptimesindex ON stoptimes (trip_id)")
	db.SQLExec("CREATE INDEX IF NOT EXISTS tripsindex on trips (route_id)")
	db.SQLExec("COMMIT")
}
