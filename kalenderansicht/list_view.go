package kalenderansicht

type pagesAmount uint8

const (
	fuenf     pagesAmount = 5
	zehn      pagesAmount = 10
	fuenfzehn pagesAmount = 15
)

type ListView struct {
}

func (c ListView) ErstelleKalenderEintraege(userId string) {
}

func (c ListView) WaehleDatum() {

}
func (c ListView) WaehleSeitenAnzahl(amount pagesAmount) {

}
func (c ListView) SpringSeiteWeiter() {

}
func (c ListView) SpringSeiteZuueck() {

}
