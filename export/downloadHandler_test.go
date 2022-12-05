/*
@author: 2447899 8689159 3000685
*/
package export

import (
	"DHBW_GO_Projekt/dateisystem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadHandler(t *testing.T) {
	dateisystem.DeleteAll(dateisystem.GetTermine("admin"), "admin")

	ts := createServer(func(name, pwd string) (bool, string) {
		return true, "test" // <--- accept any request
	})

	res := doRequestWithPassword(t, ts.URL)

	//Es wird immer eine leere Ical erwartet
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal(t, "13", res.Header.Get("Content-Length"))

	path, _ := os.Getwd()
	path = filepath.Join(path, "export")
	path = filepath.Join(path, "admin.ical")

	fmt.Println(path)
	err := os.Remove(path)
	if err != nil {
		return
	}

	ts.Close()
}
