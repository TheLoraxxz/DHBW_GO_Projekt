package dateisystem

/*
Zweck: UID für Termine bereitstellen
*/
//ToDo Dry bekämpfen
//Mat-Nr. 8689159
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// createID erzeugt neue txt, die ID hält
func createID() {

	//erzeugt Pfad zur Datei
	path := processedPath()

	//erzeuge txt
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	//Initialisiere ID
	id := 0

	//schreibe ID
	_, err2 := f.WriteString(strconv.Itoa(id))

	if err2 != nil {
		log.Fatal(err2)
	}
}

// getID liefert aktuelle ID
func getID() (id int) {

	//erzeugt Pfad zur Datei
	path := processedPath()

	//öffnet das Verzeichnis des Users
	f, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err = strconv.Atoi(string(f))

	if id < 0 {
		return 0
	}

	// erhalte aktuelle ID
	return id
}

// incrementID erhöht id des Users um 1
func incrementID() {

	//Initialisiere ID
	id := getID() + 1

	//erzeugt Pfad zur Datei
	path := processedPath()

	//lösche aktuelle txt
	err := os.Remove(path)
	if err != nil {
		return
	}

	//erzeuge txt
	f, err2 := os.Create(path)

	if err2 != nil {
		log.Fatal(err2)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	//schreibe ID
	_, err3 := f.WriteString(strconv.Itoa(id))

	if err3 != nil {
		log.Fatal(err3)
	}
}

func decrementID() {
	//Initialisiere ID
	id := getID() - 1

	//erzeugt Pfad zur Datei
	path := processedPath()

	//lösche aktuelle txt
	err := os.Remove(path)
	if err != nil {
		return
	}

	//erzeuge txt
	f, err2 := os.Create(path)

	if err2 != nil {
		log.Fatal(err2)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	//schreibe ID
	_, err3 := f.WriteString(strconv.Itoa(id))

	if err3 != nil {
		log.Fatal(err3)
	}
}

func processedPath() (path string) {
	path = GetDirectory("mik")
	path = filepath.Dir(path)
	path = filepath.Join(path, "id.txt")

	return path
}
