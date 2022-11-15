package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"flag"
	"log"
	"net/http"
)

type RootHandler struct{}
type TabellenHandler struct{}
type CreatUserHandler struct{}

var Server http.Server

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
	createUser := CreatUserHandler{}
	http.Handle("/", &root)
	http.Handle("/user/Create", &createUser)

	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
