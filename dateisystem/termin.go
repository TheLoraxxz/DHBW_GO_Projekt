/*
@author: 2447899 8689159 3000685
*/
package dateisystem

/*
Zweck: Struct inklusive Seter für Termine
*/

//Mat-Nr. 8689159
import (
	"golang.org/x/crypto/bcrypt"
	"strings"
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
	Shared      bool      `json:"Shared"`
	ID          string    `json:"ID"`
}

// createID erzeugt neue ID
func createID(dat time.Time, endDat time.Time) string {

	u := time.Now().String()

	id := dat.String() + endDat.String() + u

	//generiert Hash --> gewährleistet hohe Kollisionsfreiheit bei IDs
	bytes, _ := bcrypt.GenerateFromPassword([]byte(id), 1)
	id = string(bytes)

	//Entfernt problematische Chars aus Hash
	id = strings.Replace(id, "/", "E", 99)
	id = strings.Replace(id, ".", "D", 99)

	return id
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

func (Termin) SetShared(t *Termin, shared bool) {
	t.Shared = shared
}
