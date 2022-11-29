package export

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	var kTest []ds.Termin
	ds.CreateNewTermin("testu", "test", ds.YEARLY, time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC), time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC), false, "admin")
	k := ds.CreateNewTermin("testa", "test", ds.YEARLY, time.Date(2022, 11, 22, 14, 2, 5, 0, time.UTC), time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC), false, "admin")

	kTest = ds.GetTermine("mik")

	assert.Equal(t, k, ds.FindInCacheByID(kTest, k.ID))

	ParsToIcal(kTest, "admin")
}
