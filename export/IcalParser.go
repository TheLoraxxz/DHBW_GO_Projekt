package export

import (
	"DHBW_GO_Projekt/dateisystem"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	dateLayout = "20060102T150405"
)

func ParsToIcal(t dateisystem.Termin, username string) {
	checkForDirectory()
	file := "export/" + t.Title + ".ics"
	fileForMac := "export/" + t.Title + ".ical"
	p :=
		"BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:BoBoGo/DE\nMETHOD:PUBLISH\nBEGIN:VEVENT\nUID:" +
			username + "\nDTSTAMP:" +
			time.Now().Format(dateLayout) +
			"\n" + "DTSTART;" +
			t.Date.Format(dateLayout) +
			"\nDTEND;" + t.EndDate.Format(dateLayout) +
			"\nSUMMARY:" + t.Title + "\nDESCRIPTION:" +
			t.Description + "\nCLASS:PRIVATE\n"

	switch t.Recurring {
	case t.Recurring, 0:
		p = p + "RRULE:FREQ=YEARLY" + "\nEND:VEVENT\nEND:VCALENDAR"
	case t.Recurring, 1:
		p = p + "RRULE:FREQ=WEEKLY" + "\nEND:VEVENT\nEND:VCALENDAR"
	case t.Recurring, 2:
		p = p + "RRULE:FREQ=MONTHLY" + "\nEND:VEVENT\nEND:VCALENDAR"
	case t.Recurring, 3:
		p = p + "RRULE:FREQ=YEARLY" + "\nEND:VEVENT\nEND:VCALENDAR"
	case t.Recurring, 4:
		p = p + "\nEND:VEVENT\nEND:VCALENDAR"
	}

	writeI(file, p)
	writeI(fileForMac, p)
}

func checkForDirectory() {
	err := os.MkdirAll("export/", 755) //erzeugt das Exportverzeichnis, falls noch nicht existent
	if err != nil {
		fmt.Println(err)
	}
}

func writeI(file string, parsed string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.WriteString(parsed)
	if err2 != nil {
		log.Fatal(err2)
	}
}
