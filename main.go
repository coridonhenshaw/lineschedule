package main

import (
	"database/sql"
	"encoding/xml"
	"io/ioutil"
	"log"

	"lineschedule/db"
	imp "lineschedule/importCSV"

	"github.com/integrii/flaggy"
	_ "github.com/mattn/go-sqlite3"
)

type XMLStruct struct {
	LineSchedule []PackageStruct `xml:"Journey"`
}

type PackageStruct struct {
	Output      string
	ServiceDate string
	Route       []RouteStruct
	Connection  []ConnectionStruct
	StopName    []StopNamesStruct
}

type StopNamesStruct struct {
	Stop_id string `xml:"StopID"`
	Name    string
}

type GetScheduleStruct struct {
	Route               string
	ServiceDate         string
	Direction           int
	Stops               []string
	OrientationVertical bool
}

type GetRouteInformationStruct struct {
	Route       string
	ServiceDate string
}

func main() {

	var DBName string = "gtfs.sqlite3"
	var GS GetScheduleStruct
	var GR GetRouteInformationStruct
	var XMLFile string = "package.xml"

	flaggy.String(&DBName, "db", "dbname", "Override database filename")

	subcommandImport := flaggy.NewSubcommand("import")
	subcommandImport.Description = "Import GTFS data from directory into database."
	flaggy.AttachSubcommand(subcommandImport, 1)

	subcommandGetRouteInformation := flaggy.NewSubcommand("information")
	subcommandGetRouteInformation.Description = "Get information for route"
	subcommandGetRouteInformation.String(&GR.Route, "r", "route", "Route identifier")
	subcommandGetRouteInformation.String(&GR.ServiceDate, "s", "servicedate", "Date (YYYY-MM-DD)")
	flaggy.AttachSubcommand(subcommandGetRouteInformation, 1)

	subcommandSchedule := flaggy.NewSubcommand("schedule")
	subcommandSchedule.Description = "Get schedule for route"
	subcommandSchedule.String(&GS.Route, "r", "route", "Route identifier")
	subcommandSchedule.String(&GS.ServiceDate, "s", "servicedate", "Date (YYYY-MM-DD)")
	subcommandSchedule.Int(&GS.Direction, "d", "direction", "Direction")
	subcommandSchedule.StringSlice(&GS.Stops, "st", "stops", "Show service only at specified stop IDs.")
	subcommandSchedule.Bool(&GS.OrientationVertical, "v", "vertical", "Arrange trips vertically")
	flaggy.AttachSubcommand(subcommandSchedule, 1)

	subcommandCSVSchedule := flaggy.NewSubcommand("csvschedule")
	subcommandCSVSchedule.Description = "Make line schedule CSV(s) from XML configuration"
	subcommandCSVSchedule.String(&XMLFile, "if", "file", "XML file")
	flaggy.AttachSubcommand(subcommandCSVSchedule, 1)

	flaggy.Parse()

	var err error

	db.Con, err = sql.Open("sqlite3", DBName)
	if err != nil {
		log.Panic(err)
	}
	defer db.Con.Close()

	if subcommandImport.Used {
		if len(flaggy.TrailingArguments) == 0 {
			log.Panic("No import directory provided.")
		}
		imp.Import(flaggy.TrailingArguments[0])
	} else if subcommandGetRouteInformation.Used {
		GetRouteInformation(GR.Route, GR.ServiceDate)
	} else if subcommandSchedule.Used {
		ts := GetSchedule(GS.Route, GS.ServiceDate, GS.Direction, GS.Stops)
		if GS.OrientationVertical {
			PrintScheduleVertical(ts)
		} else {
			PrintSchedule(ts)
		}
	} else if subcommandCSVSchedule.Used {
		var Packages XMLStruct

		Content, err := ioutil.ReadFile(XMLFile)
		if err != nil {
			log.Fatal(err)
		}

		err = xml.Unmarshal(Content, &Packages)
		if err != nil {
			log.Panic(err)
		}

		for _, Package := range Packages.LineSchedule {
			ts := GetTable(Package.ServiceDate, Package.Route, Package.Connection, Package.StopName)
			WriteCSV(Package.Output, ts)
		}
	}
}
