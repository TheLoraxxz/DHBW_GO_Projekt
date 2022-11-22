package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type RootHandler struct{}
type ChangeUserHandler struct{}
type UserHandler struct{}
type LogoutHandler struct{}
type CreatUserHandler struct{}

var Server http.Server

func main() {
	//flags and configuration of application
	port := flag.String("port", "80", "define the port for the application")
	adminUserName := flag.String("user", "admin", "Define Admin username for first login")
	adminPassword := flag.String("passw", "admin", "Define Admin Password for first login to application")
	basepath, err := os.Getwd()
	err = authentifizierung.LoadUserData(adminUserName, adminPassword, &basepath)
	timerSaveData := time.NewTimer(1 * time.Minute)
	go func() {
		<-timerSaveData.C
		err := authentifizierung.SaveUserData(&basepath)
		if err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Coudn't load users")
	}
	Server = http.Server{
		Addr: ":" + *port,
	}

	//http handles
	//hier weitere handler hinzufügen in ähnlicher fashion für die verschiedenen Templates
	root := RootHandler{}
	createUser := CreatUserHandler{}
	changeUser := ChangeUserHandler{}
	user := UserHandler{}
	logout := LogoutHandler{}
	http.Handle("/", &root)
	http.Handle("/user/create", &createUser)
	http.Handle("/user/change", &changeUser)
	http.Handle("/user", &user)
	http.Handle("/logout", &logout)

	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
