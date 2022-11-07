package kalenderansicht

import (
	"time"
)

type TabellenKalendar struct {
	Datumsanzeige time.Time
}

func (c TabellenKalendar) JahrAnzeige() int {
	return c.Datumsanzeige.Year()
}
func (c TabellenKalendar) MonatsAnzeige() time.Month {
	return c.Datumsanzeige.Month()
}

func (c TabellenKalendar) createCalendarEntries(userId string) {

}

func (c *TabellenKalendar) springMonatVor() {
	switch day := uint8(c.Datumsanzeige.Day()); {
	case day < 29:
		c.Datumsanzeige = c.Datumsanzeige.AddDate(0, 1, 0)
	case day > 30:
		c.Datumsanzeige = c.Datumsanzeige.AddDate(0, 0, 5)
	}
}

func (c *TabellenKalendar) springMonatZurueck() {
	switch day := uint8(c.Datumsanzeige.Day()); {
	case day < 29:
		c.Datumsanzeige = c.Datumsanzeige.AddDate(0, -1, 0)
	case day > 30:
		c.Datumsanzeige = c.Datumsanzeige.AddDate(0, 0, -31)
	}
}

func (c *TabellenKalendar) waehleMonat(month time.Time) {

}

func (c *TabellenKalendar) springZuHeute() {

}
