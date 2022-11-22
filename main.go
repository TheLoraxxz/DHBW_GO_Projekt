package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	ka "DHBW_GO_Projekt/kalenderansicht"
	"flag"
	"html/template"
	"log"
	"net/http"
)

type RootHandler struct{}
type ViewmanagerHandler struct {
	vm             *ka.ViewManager
	cookie         string
	viewmanagerTpl *template.Template
}
type CreatUserHandler struct {
}

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
	viewmanagerHandler := ViewmanagerHandler{}

	http.Handle("/", &root)
	http.Handle("/user/Create", &createUser)

	http.Handle("/user/view/", &viewmanagerHandler)
	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
