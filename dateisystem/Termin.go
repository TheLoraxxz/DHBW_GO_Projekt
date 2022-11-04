package dateisystem

//Mat-Nr. 8689159
import (
	"time"
)

type repeat int

const ( //"enum" um Wiederholung anzuzeigen
	taeglich repeat = iota
	woechentlich
	monatlich
	jaehrlich
	niemals
)

const ( //genutzt zum Formatieren von time.Date() Objekten
	dateLayoutISO = "2006-01-02T15:04:05 UTC"
)

type Termin struct {
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Recurring   repeat    `json:"Recurring"`
	Date        time.Time `json:"Date"`
	EndDate     time.Time `json:"EndDate"`
}

func setTitle(t *Termin, newTitle string) {
	t.Title = newTitle
}

func setDescription(t *Termin, newDescription string) {
	t.Description = newDescription
}

func setRecurring(t *Termin, newRecurring repeat) {
	t.Recurring = newRecurring
}

func setDate(t *Termin, newDate string) {
	d, _ := time.Parse(dateLayoutISO, newDate)
	t.Date = d
}

func setEndeDate(t *Termin, newDate string) {
	d, _ := time.Parse(dateLayoutISO, newDate)
	t.EndDate = d
}
