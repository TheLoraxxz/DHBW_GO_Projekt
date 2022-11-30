package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	ka "DHBW_GO_Projekt/kalenderansicht"
	"DHBW_GO_Projekt/terminfindung"
	"flag"
	"fmt"
	"html/template"
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
type ViewManagerHandler struct {
	vm             *ka.ViewManager
	viewManagerTpl *template.Template
}

var Server http.Server

func main() {
	//flags and configuration of application
	port := flag.String("port", "80", "define the port for the application. Default:80")
	adminUserName := flag.String("user", "admin", "Define Admin username for first login. Default: admin")
	adminPassword := flag.String("passw", "admin", "Define Admin Password for first login to application. Defualt: admin")
	basepath, err := os.Getwd()
	// load user data from plate and if not create a new user
	err = authentifizierung.LoadUserData(adminUserName, adminPassword, &basepath)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Coudn't load users")
	}
	err = terminfindung.LoadDataToSharedTermin(&basepath)
	if err != nil {
		fmt.Println(err)
	}
	//set a timer to write all users to plate every minute
	timerSaveData := time.NewTimer(1 * time.Minute * 15)
	go func() {
		// timer waits until one minute is over and then saves the new data
		<-timerSaveData.C
		saveAuthErr := authentifizierung.SaveUserData(&basepath)
		if saveAuthErr != nil {
			fmt.Println(saveAuthErr)
		}
		// save every 30 minutes the whole shared termine
		saveSharedErr := terminfindung.SaveSharedTermineToDisk(&basepath)
		if saveSharedErr != nil {
			fmt.Println(saveSharedErr)
		}
	}()

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
	viewManagerHdl := ViewManagerHandler{}

	// Templates für die Tabellenansicht sowie die Listenansicht erstellen
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get rooted path name corresponding to the current directory")
	}
	viewManagerHdl.viewManagerTpl = template.Must(template.New("tbl.html").ParseFiles(path+"/assets/sites/tbl.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html", path+"/assets/templates/listing.html"))
	template.Must(viewManagerHdl.viewManagerTpl.New("liste.html").ParseFiles(path+"/assets/sites/liste.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html"))
	template.Must(viewManagerHdl.viewManagerTpl.New("editor.html").ParseFiles(path+"/assets/sites/editor.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/listing.html"))
	template.Must(viewManagerHdl.viewManagerTpl.New("filterTermins.html").ParseFiles(path + "/assets/sites/filterTermins.html"))
	http.Handle("/", &root)
	http.Handle("/user/create", &createUser)
	http.Handle("/user/change", &changeUser)
	http.Handle("/user", &user)
	http.Handle("/user/view/", &viewManagerHdl)
	http.Handle("/logout", &logout)
	http.HandleFunc("/shared", AdminSiteServeHTTP)
	http.HandleFunc("/shared/create/link", CreateLinkServeHTTP)
	http.HandleFunc("/shared/create/app", ServeHTTPSharedAppCreateDate)
	http.HandleFunc("/shared/showAllLink", ShowAllLinksServeHttp)
	http.HandleFunc("/shared/public", PublicSharedWebsite)
	// start server
	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
