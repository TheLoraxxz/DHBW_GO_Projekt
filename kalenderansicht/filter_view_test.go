package kalenderansicht

import "testing"

func TestGetEntries(t *testing.T) {

}

func TestFilterView(t *testing.T) {
	//slice mit Testterminen erstellen, benötigt viel Zeit: daher ein globales Slice
	testTermine30 = create30TestTermins()

	//Tests für Custom-Settings innerhalb der Webseite
	t.Run("testRuns GetEntries", TestGetEntries)
}
