/*
@author: 2447899 8689159 3000685
*/
package dateisystem

/*
Zweck: für das Dateisystem relative Pfade bereitzustellen um zu garantieren, das aus anderen Modulen korrekt auf Dateien zugegriffen werden kann
*/

//Mat-Nr. 8689159
import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// GetDirectory liefert Pfad zu /dateisystem anhand des usernames
func GetDirectory(username string) (path string) {
	wd, err := os.Getwd() // lädt aktuelles Verzeichnis
	if err != nil {
		fmt.Println(err)
	}
	if filepath.Base(wd) != "DHBW_GO_Projekt" { // prüft ob der Projektordner dem wd entspricht, falls nein wird das letzt Pfad Element gelöscht
		wd = filepath.Dir(wd)
	}
	wd = filepath.Join(wd, "dateisystem") // geht in Modul Dateisystem
	path = filepath.Join(wd, username)    //formt Dateipfad, indem username an das wd angehangen wird

	if _, err := os.Stat(path); err == nil {
		return path
	} else if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 755) //erzeugt Verzeichnis passend zum User, falls noch nicht existent
		if err != nil {
			fmt.Println(err)
		}
		return path
	}
	return path
}

// getFile liefert Pfad zur Datei anhand des Titels und des Usernamens, relevant zum Suchen und speichern von Termine
func getFile(title string, username string) string {
	path := GetDirectory(username)
	fileName := title + ".json"

	path = filepath.Join(path, fileName)
	return path
}
