package export

import (
	"DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	var kTest []dateisystem.Termin
	kTest = dateisystem.AddToCache(dateisystem.NewTerminObj("testu", "test", dateisystem.YEARLY, "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC"), kTest)
	assert.Equal(t, dateisystem.NewTerminObj("testu", "test", dateisystem.YEARLY, "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC"), kTest[0])
	ParsToIcal(kTest[0], "mik")
}
