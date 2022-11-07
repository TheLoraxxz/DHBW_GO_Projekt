package kalenderansicht

type pagesAmount uint8

const (
	fuenf     pagesAmount = 5
	zehn      pagesAmount = 10
	fuenfzehn pagesAmount = 15
)

type ListenAnsicht struct {
}

func (c ListenAnsicht) createCalendarEntries(userId string) {

}

func (c ListenAnsicht) waehleDatum() {

}
func (c ListenAnsicht) waehleSeitenAnzahl(amount pagesAmount) {

}
func (c ListenAnsicht) springSeiteWeiter() {

}
func (c ListenAnsicht) springSeiteZuueck() {

}
