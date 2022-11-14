package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"flag"
	"log"
	"net/http"
)

type RootHandler struct{}
type TabellenHandler struct{}

var Server = http.Server{
	Addr: ":80",
}

func main() {
	//flags and configuration of application
	port := flag.String("port", "80", "define the port for the application")
	adminUserName := flag.String("user", "admin", "Define Admin username for first login")
	adminPassword := flag.String("passw", "admin", "Define Admin Password for first login to application")
	authentifizierung.CreateUser(adminUserName, adminPassword)
	Server = http.Server{
		Addr: ":" + *port,
	}

	//http handles
	//hier weitere handler hinzufügen in ähnlicher fashion für die verschiedenen Templates
	root := RootHandler{}
	tabelle := TabellenHandler{}
	http.Handle("/", &root)
	http.Handle("/kalender/", &tabelle)

	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
