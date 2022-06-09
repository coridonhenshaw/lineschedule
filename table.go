package main

import (
	"log"
	"sort"
)

type EntryStruct struct {
	stop_id        string
	stop_name      string
	arrival_time   int
	departure_time int
	Valid          bool
	Connection     bool
	Visible        bool
}

type ColumnStruct struct {
	stop_id    string
	Connection bool
}

type RowStruct struct {
	Name  string
	Entry []EntryStruct
}

type StopStruct struct {
	Name    string
	stop_id string
	Count   int
}

type TableStruct struct {
	Rows  []RowStruct
	Stops []StopStruct

	StopNameMap  map[string]string
	StopCountMap map[string]int
	Columns      []ColumnStruct
	Scratch      []RowStruct

	colmap map[string]int
}

func (a TableStruct) Len() int      { return len(a.Rows) }
func (a TableStruct) Swap(i, j int) { a.Rows[i], a.Rows[j] = a.Rows[j], a.Rows[i] }
func (a TableStruct) Less(i, j int) bool {
	for x := range a.Rows[i].Entry {
		if a.Rows[i].Entry[x].Valid && a.Rows[j].Entry[x].Valid {
			return a.Rows[i].Entry[x].arrival_time < a.Rows[j].Entry[x].arrival_time
		}
	}

	var I1, T1 int
	for x := range a.Rows[i].Entry {
		if a.Rows[i].Entry[x].Valid {
			T1 += a.Rows[i].Entry[x].arrival_time
			I1++
		}
	}
	var I2, T2 int
	for x := range a.Rows[j].Entry {
		if a.Rows[j].Entry[x].Valid {
			T2 += a.Rows[j].Entry[x].arrival_time
			I2++
		}
	}

	T1 /= I1
	T2 /= I2

	return T1 < T2

	// for x := range a.Rows[i].Entry {
	// 	if a.Rows[i].Entry[x].Valid || a.Rows[j].Entry[x].Valid {
	// 		return a.Rows[i].Entry[x].arrival_time < a.Rows[j].Entry[x].arrival_time
	// 	}
	// }
	// return false
}

func (o *TableStruct) InColumns(Key string) int {
	i, v := o.colmap[Key]
	if v {
		return i
	}
	return -1
}

func (o *TableStruct) InsertColumn(Entry EntryStruct, Position int) {

	if Position < len(o.Columns) && o.Columns[Position].Connection == true {
		for {
			Position++
			if Position >= len(o.Columns) {
				break
			}

			if o.Columns[Position].Connection != true {
				break
			}
		}
	}

	var Result []ColumnStruct

	var NewColumn ColumnStruct
	NewColumn.stop_id = Entry.stop_id
	NewColumn.Connection = Entry.Connection

	Result = append(Result, o.Columns[0:Position]...)
	Result = append(Result, NewColumn)
	Result = append(Result, o.Columns[Position:]...)

	o.Columns = Result

	if o.colmap == nil {
		o.colmap = make(map[string]int)
	}

	o.colmap[Entry.stop_id] = Position
}

func (o *TableStruct) AddRow(Name string, Keys []EntryStruct) {
	if o.StopNameMap == nil {
		o.StopNameMap = make(map[string]string)
		o.StopCountMap = make(map[string]int)
	}

	var Row RowStruct
	Row.Name = Name

	var LastInsertPosition int
	for _, v := range Keys {
		FoundColumn := o.InColumns(v.stop_id)
		if FoundColumn == -1 {
			o.InsertColumn(v, LastInsertPosition)
			LastInsertPosition++
		} else {
			LastInsertPosition = FoundColumn + 1
		}
		o.StopNameMap[v.stop_id] = v.stop_name
		o.StopCountMap[v.stop_id]++

		Row.Entry = append(Row.Entry, v)
	}
	o.Scratch = append(o.Scratch, Row)
}

func (o *TableStruct) Finalize() {
	Width := len(o.Columns)
	o.Stops = make([]StopStruct, Width)
	o.Rows = make([]RowStruct, len(o.Scratch))
	for i := range o.Rows {
		o.Rows[i].Entry = make([]EntryStruct, Width)
		o.Rows[i].Name = o.Scratch[i].Name
	}

	ColMap := make(map[string]int)
	for i, v := range o.Columns {
		ColMap[v.stop_id] = i
		St := StopStruct{Name: o.StopNameMap[v.stop_id], stop_id: v.stop_id, Count: o.StopCountMap[v.stop_id]}

		o.Stops[i] = St
	}

	for CurrRow, x := range o.Scratch {
		for _, y := range x.Entry {
			Col, Valid := ColMap[y.stop_id]
			if !Valid {
				log.Panic()
			}
			o.Rows[CurrRow].Entry[Col] = y
		}
	}

	// for i := range o.Rows {
	// 	var Last int
	// 	for j := range o.Rows[i].Entry {
	// 		var e = &o.Rows[i].Entry[j]

	// 		if (*e).Valid {
	// 			Last = (*e).arrival_time
	// 			continue
	// 		}

	// 		if Last == 0 {
	// 			//was j < len
	// 			for k := j; k < len(o.Rows[i].Entry); k++ {
	// 				if o.Rows[i].Entry[k].Valid && !o.Rows[i].Entry[k].Connection {
	// 					Last = o.Rows[i].Entry[k].arrival_time
	// 					break
	// 				}
	// 			}
	// 		}
	// 		(*e).arrival_time = Last
	// 	}
	// }

	// for _, v := range o.Rows {
	// 	fmt.Println(v.Name, v.Entry)
	// }

	sort.Stable(*o)
}
