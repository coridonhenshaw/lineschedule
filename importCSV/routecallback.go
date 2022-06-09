package imp

import "lineschedule/db"

type RouteCallbackStruct struct {
	LineMap *map[string]int
}

func (o *RouteCallbackStruct) SetLinemap(Map *map[string]int) {
	o.LineMap = Map
}

func (o *RouteCallbackStruct) Line(Line *[]string) {
	var route_id = (*Line)[(*o.LineMap)["route_id"]]
	var route_short_name = (*Line)[(*o.LineMap)["route_short_name"]]
	var route_long_name = (*Line)[(*o.LineMap)["route_long_name"]]

	var SQL = "INSERT INTO routes (route_id, route_short_name, route_name) VALUES (?, ?, ?)"
	db.SQLExec(SQL, route_id, route_short_name, route_long_name)
}
