package export

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	var kTest []ds.Termin
	kTest = ds.AddToCache(ds.NewTerminObj("testu", "test", ds.YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC)), kTest)
	assert.Equal(t, ds.NewTerminObj("testu", "test", ds.YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC)), kTest[0])
	ParsToIcal(kTest[0], "mik")
}

func TestWriteFile(t *testing.T) {
	termin := ds.CreateNewTermin("test", "test", ds.WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik")
	terminLoaded := ds.LoadTermin("test", "mik")

	assert.Equal(t, termin, terminLoaded)
}
