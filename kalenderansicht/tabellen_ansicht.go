package kalenderansicht

import (
	"strconv"
	"time"
)

type TabellenAnsicht struct {
	Datumsanzeige time.Time
}

func (c TabellenAnsicht) JahrAnzeige() int {
	return c.Datumsanzeige.Year()
}
func (c TabellenAnsicht) JahrAnzeige2(zahlStr ...string) int {
	zahl, _ := strconv.Atoi(zahlStr[0])
	jahr := c.Datumsanzeige.Year() + zahl
	return jahr
}
func (c TabellenAnsicht) MonatsAnzeige() time.Month {
	return c.Datumsanzeige.Month()
}

func (c TabellenAnsicht) ErstelleKalenderEintraege(userId string) {

}
func (c TabellenAnsicht) ErstelleTabellenAnsicht() []int {
	return make([]int, 5)
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
func (c *TabellenAnsicht) SpringeJahr(summand int) {
	c.Datumsanzeige = c.Datumsanzeige.AddDate(summand, 0, 0)
}

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
