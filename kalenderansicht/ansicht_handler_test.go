package kalenderansicht

import (
	"log"
	"net/http"
	"testing"
)

func testRuns(t *testing.T) {
	for true {
		http.HandleFunc("/tabellenAnsicht", tabellenHandler)
		http.HandleFunc("/listenAnsicht", listenHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

}
func TestCalendarView(t *testing.T) {
	t.Run("testRuns", testRuns)
}
