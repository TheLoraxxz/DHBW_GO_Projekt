package kalenderansicht

type pagesAmount uint8

const (
	fuenf     pagesAmount = 5
	zehn      pagesAmount = 10
	fuenfzehn pagesAmount = 15
)

type ListenAnsicht struct {
}

func (c ListenAnsicht) ErstelleKalenderEintraege(userId string) {
}

func (c ListenAnsicht) WaehleDatum() {

}
func (c ListenAnsicht) WaehleSeitenAnzahl(amount pagesAmount) {

}
func (c ListenAnsicht) SpringSeiteWeiter() {

}
func (c ListenAnsicht) SpringSeiteZuueck() {

}
