package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

func UncookTime(StopTime int) string {
	h := StopTime / (60 * 60)
	StopTime -= h * (60 * 60)
	m := StopTime / 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

func GetRouteInformation(Route string, ServiceDate string) {
	fmt.Print("Service on route ", Route, " for ", ServiceDate, ":\n\n")

	for i := 0; i < 2; i++ {
		ts := GetSchedule(Route, ServiceDate, i, nil)
		fmt.Print("Direction ", i, ":\n")
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(w, "Stop ID \tTrips at stop \tStop name")

		l := len(ts.Rows)
		for _, v := range ts.Stops {
			fmt.Fprintf(w, "%v\t%v/%v\t%v\n", v.stop_id, v.Count, l, v.Name)
		}
		w.Flush()
		fmt.Println()
	}

}

func PrintSchedule(ts *TableStruct) {

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "\t")
	for i := range ts.Columns {
		fmt.Fprint(w, ts.Stops[i].stop_id, ": ", ts.Stops[i].Name, "\t")
	}
	fmt.Fprintln(w, "")

	for _, v := range ts.Rows {
		fmt.Fprint(w, v.Name, "\t")
		for _, b := range v.Entry {
			if b.Valid {
				fmt.Fprint(w, UncookTime(b.departure_time))
			} else {
				fmt.Fprint(w, "-")
			}
			fmt.Fprint(w, "\t")
		}
		fmt.Fprintln(w, "")
	}
}

func PrintScheduleVertical(ts *TableStruct) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	var VariableHeadsigns bool
	var LastHeadsign = ts.Rows[0].Name
	for _, v := range ts.Rows {
		if LastHeadsign != v.Name {
			VariableHeadsigns = true
			break
		}
	}

	if !VariableHeadsigns {
		fmt.Fprintf(w, "%v\n", LastHeadsign)
	} else {
		fmt.Fprintf(w, "\t")
		for _, v := range ts.Rows {
			fmt.Fprint(w, v.Name, "\t")
		}
		fmt.Fprintln(w, "")
	}

	for i, v := range ts.Stops {

		fmt.Fprint(w, v.stop_id, ": ", v.Name, "\t")

		for _, b := range ts.Rows {
			if b.Entry[i].Valid {
				fmt.Fprint(w, UncookTime(b.Entry[i].departure_time))
			}
			fmt.Fprint(w, "\t")
		}
		fmt.Fprintln(w, "")
	}
	w.Flush()
}

func WriteCSV(Filename string, ts *TableStruct) {
	file, err := os.Create(Filename)
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	Row := make([]string, 1)

	for i := range ts.Columns {
		Row = append(Row, ts.Stops[i].stop_id)
	}
	w.Write(Row)

	Row = make([]string, 1)

	for i := range ts.Columns {
		Row = append(Row, ts.Stops[i].Name)
	}
	w.Write(Row)

	for _, v := range ts.Rows {
		Row = make([]string, len(ts.Columns)+1)
		Row[0] = v.Name
		for i, b := range v.Entry {
			if b.Valid && b.Visible {
				Row[1+i] = UncookTime(b.arrival_time)
				//			}
			} else {
				//				Row[1+i] = "F:" + UncookTime(b.arrival_time)
			}
		}
		w.Write(Row)
	}
}
