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
	port := flag.String("port", "80", "define the port for the application. Default:80")
	adminUserName := flag.String("user", "admin", "Define Admin username for first login. Default: admin")
	adminPassword := flag.String("passw", "admin", "Define Admin Password for first login to application. Defualt: admin")
	basepath, err := os.Getwd()
	// load user data from plate and if not create a new user
	err = authentifizierung.LoadUserData(adminUserName, adminPassword, &basepath)
	//set a timer to write all users to plate every minute
	timerSaveData := time.NewTimer(1 * time.Minute)
	go func() {
		// timer waits until one minute is over and then saves the new data
		<-timerSaveData.C
		errOnSave := authentifizierung.SaveUserData(&basepath)
		if errOnSave != nil {
			fmt.Println(errOnSave)
		}
	}()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Coudn't load users")
	}
	// setup server
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
	http.HandleFunc("/shared", AdminSiteServeHTTP)
	http.HandleFunc("/shared/create/link", CreateLinkServeHTTP)
	http.HandleFunc("/shared/create/app", ServeHTTPSharedAppCreateDate)
	http.HandleFunc("/shared/showAllLink", ShowAllLinksServeHttp)
	http.HandleFunc("/test", CreateTest)
	// start server
	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
