package dateisystem

//Mat-Nr. 8689159
import (
	"time"
)

type Repeat int

const ( //"enum" um Wiederholung anzuzeigen
	DAILY Repeat = iota
	WEEKLY
	MONTHLY
	YEARLY
	Never
)

const ( //genutzt zum Formatieren von time.Date() Objekten
	dateLayoutISO = "2006-01-02T15:04:05 UTC"
)

type Termin struct {
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Recurring   Repeat    `json:"Recurring"`
	Date        time.Time `json:"Date"`
	EndDate     time.Time `json:"EndDate"`
}

func (Termin) SetTitle(t *Termin, newTitle string) {
	t.Title = newTitle
}

func (Termin) SetDescription(t *Termin, newDescription string) {
	t.Description = newDescription
}

func (Termin) SetRecurring(t *Termin, newRecurring Repeat) {
	t.Recurring = newRecurring
}

func (Termin) SetDate(t *Termin, newDate time.Time) {
	t.Date = newDate
}

func (Termin) SetEndeDate(t *Termin, newDate time.Time) {
	t.EndDate = newDate
}
