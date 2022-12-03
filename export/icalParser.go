package export

/*
Zweck: Generiert die zu exportierende Ical
*/

//Mat-Nr. 8689159
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

func ParsToIcal(k []dateisystem.Termin, username string) string {
	createDirectory() // legt export Verzeichnis an
	fileForMac := "export/" + username + ".ical"

	//schreibt Ical
	p :=
		"BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:BoBoGo/DE\nMETHOD:PUBLISH\nBEGIN:VEVENT\nUID:" +
			username + "\nDTSTAMP:" +
			time.Now().Format(dateLayout) +
			"\n"

	for i := 0; i < len(k); i++ {

		p = p + "DTSTART:" +
			k[i].Date.Format(dateLayout) + "Z" +
			"\nDTEND:" + k[i].EndDate.Format(dateLayout) + "Z" +
			"\nSUMMARY:" + k[i].Title + "\nDESCRIPTION:" +
			k[i].Description + "\nCLASS:PRIVATE\n"

		switch k[i].Recurring {
		case 0:
			p = p + "RRULE:FREQ=YEARLY\n"
		case 1:
			p = p + "RRULE:FREQ=WEEKLY\n"
		case 2:
			p = p + "RRULE:FREQ=MONTHLY\n"
		case 3:
			p = p + "RRULE:FREQ=YEARLY\n"

		}
	}

	p = p + "END:VEVENT\nEND:VCALENDAR"

	writeIcal(fileForMac, p)

	return fileForMac
}

func createDirectory() {
	err := os.MkdirAll("export/", 755) //erzeugt das Exportverzeichnis, falls noch nicht existent
	if err != nil {
		fmt.Println(err)
	}
}

func writeIcal(file string, parsed string) {
	f, err := os.Create(file) //legt Export-Datei an
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	_, err2 := f.WriteString(parsed) //schreibt String
	if err2 != nil {
		log.Fatal(err2)
	}
}
