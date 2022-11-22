package export

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	var kTest []ds.Termin
	kTest = ds.AddToCache(ds.NewTerminObj("testu", "test", ds.YEARLY, time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC), time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC)), kTest)
	kTest = ds.AddToCache(ds.NewTerminObj("testa", "test", ds.YEARLY, time.Date(2022, 11, 22, 14, 2, 5, 0, time.UTC), time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC)), kTest)
	assert.Equal(t, ds.NewTerminObj("testu", "test", ds.YEARLY, time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC), time.Date(2022, 11, 22, 15, 2, 5, 0, time.UTC)), kTest[0])
	ParsToIcal(kTest, "mik")
}
