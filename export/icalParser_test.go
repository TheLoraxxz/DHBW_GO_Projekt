package export

//Mat-Nr. 8689159
import (
	ds "DHBW_GO_Projekt/dateisystem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEmptyIcal(t *testing.T) {
	//setup
	ds.DeleteAll(ds.GetTermine("admin"), "admin")

	//erzeuge Ical
	ParsToIcal(ds.GetTermine("admin"), "admin")
	//String gegen den getestet werden soll
	validEmptyIcal := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:BoBoGo/DE\nMETHOD:PUBLISH\nBEGIN:VEVENT\nUID:admin\nDTSTAMP:" + time.Now().Format(dateLayout) + "\nEND:VEVENT\nEND:VCALENDAR"
	//erzeuge Pfad zur Datei
	path, _ := os.Getwd()
	path = filepath.Join(path, "export")
	path = filepath.Join(path, "admin.ical")
	buf, _ := os.ReadFile(path)
	file := string(buf)

	assert.Equal(t, validEmptyIcal, file)
	//teardown
	fmt.Println(path)
	err := os.Remove(path)
	if err != nil {
		return
	}
}

func TestLoadedIcal(t *testing.T) {
	//setup
	ds.DeleteAll(ds.GetTermine("admin"), "admin")
	ds.CreateNewTermin("testu", "test", ds.Never, time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC), time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC), false, "admin")
	//erzeuge Ical
	ParsToIcal(ds.GetTermine("admin"), "admin")
	//String gegen den getestet werden soll
	validIcal := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:BoBoGo/DE\nMETHOD:PUBLISH\nBEGIN:VEVENT\nUID:admin\nDTSTAMP:" + time.Now().Format(dateLayout) + "\nDTSTART:20221122T150205Z\nDTEND:20221122T150205Z\nSUMMARY:testu\nDESCRIPTION:test\nCLASS:PRIVATE\nEND:VEVENT\nEND:VCALENDAR"
	//erzeuge Pfad zur Datei
	path, _ := os.Getwd()
	path = filepath.Join(path, "export")
	path = filepath.Join(path, "admin.ical")
	buf, _ := os.ReadFile(path)
	file := string(buf)

	assert.Equal(t, validIcal, file)

	//teardown
	ds.DeleteAll(ds.GetTermine("admin"), "admin")
	fmt.Println(path)
	err := os.Remove(path)
	if err != nil {
		return
	}
}
