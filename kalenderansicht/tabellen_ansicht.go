package kalenderansicht

import (
	"time"
)

type TabellenAnsicht struct {
	Datumsanzeige time.Time
}

func (c TabellenAnsicht) JahrAnzeige() int {
	return c.Datumsanzeige.Year()
}
func (c TabellenAnsicht) MonatsAnzeige() time.Month {
	return c.Datumsanzeige.Month()
}

func (c TabellenAnsicht) ErstelleKalenderEintraege(userId string) {

}

func (c *TabellenAnsicht) SpringMonatVor() {
	switch day := uint8(c.Datumsanzeige.Day()); {
	case day < 29:
		c.Datumsanzeige = c.Datumsanzeige.AddDate(0, 1, 0)
	case day >= 29:
		c.Datumsanzeige = c.Datumsanzeige.AddDate(0, 0, 5)
	}
}

func (c *TabellenAnsicht) SpringMonatZurueck() {
	switch day := uint8(c.Datumsanzeige.Day()); {
	case day < 29:
		c.Datumsanzeige = c.Datumsanzeige.AddDate(0, -1, 0)
	case day >= 29:
		c.Datumsanzeige = c.Datumsanzeige.AddDate(0, 0, -31)
	}
}

//Date(year int, month Month, day int, hour int, min int, sec int, nsec int, loc *Location) Time

func (c *TabellenAnsicht) WaehleMonat(monat time.Month) {
	jahr := c.Datumsanzeige.Year()
	c.Datumsanzeige = time.Date(
		jahr,
		monat,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func (c *TabellenAnsicht) SpringZuHeute() {
	c.Datumsanzeige = time.Now()
}
