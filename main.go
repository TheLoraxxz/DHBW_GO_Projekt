package main

import (
	"log"
	"net/http"
)

type RootHandler struct{}
type LoginHandler struct{}

var Server = http.Server{
	Addr: ":80",
}

func main() {
	//hier weitere handler hinzufügen in ähnlicher fashion für die verschiedenen Templates
	root := RootHandler{}
	http.Handle("/", &root)
	http.HandleFunc("/tabellenAnsicht", TableHandler)
	http.HandleFunc("/listenAnsicht", ListHandler)
	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
