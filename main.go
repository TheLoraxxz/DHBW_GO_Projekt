package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"DHBW_GO_Projekt/export"
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

var errorconfigs = make(map[string]string)
var Server http.Server

func main() {
	//flags and configuration of application
	port := flag.String("port", "443", "define the port for the application. Default: 443")
	adminUserName := flag.String("user", "admin", "Define Admin username for first login. Default: admin")
	adminPassword := flag.String("passw", "admin", "Define Admin Password for first login to application. Defualt: admin")
	basepath, err := os.Getwd()
	// load user data from plate and if not create a new user
	err = authentifizierung.LoadUserData(adminUserName, adminPassword, &basepath)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Coudn't load users")
	}
	terminfindung.LoadDataToSharedTermin(&basepath)
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
	setErrorconfigs()
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
	http.HandleFunc("/download", export.WrapperAuth(export.AuthenticatorFunc(export.CheckUserValid), export.DownloadHandler))
	http.HandleFunc("/error", ErrorSite_ServeHttp)
	// start server
	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
func setErrorconfigs() {
	errorconfigs["shared_admin_WrongSelected"] = "Falsches Datum selektiert"
	errorconfigs["emptyError"] = "Interner Error Problem"
	errorconfigs["wrongAuthentication"] = "Falsche Authentifizierung / falsche Daten eingegeben"
	errorconfigs["shared_wrong_terminId"] = "Konnte nicht den Termin Findung."
	errorconfigs["internal"] = "Interner Server error"
	errorconfigs["shared_coudntCreatePerson"] = "Person schon vorhanden oder falsche Zeichen enthalten"
	errorconfigs["wrong_date_format"] = "Falsches Datenformat eingegeben"
	errorconfigs["dateIsAfter"] = "Das Startdatum ist nach dem Enddatum - bitte ändern"
	errorconfigs["No_valid_repetition"] = "Die Termin-Wiederholung ist ungültig"
	errorconfigs["Missing_title"] = "Es wurde kein Termin-Titel eingegeben."
	errorconfigs["Missing_description"] = "Es wurde keine Termin-Beschreibung eingegeben."
	errorconfigs["wrong_editing_mode"] = "Dieser Bearbeitungsmodus ist ungültig."
	errorconfigs["Unvalid_Entries_Per_Page"] = "Ungültige Anzahl an Einträgen eingegeben."
}

// ErrorSite_ServeHttp
// this site renders error --> has type and link as get input and then renders the according errorconfigs
// for using errorconfigs pls see above
func ErrorSite_ServeHttp(writer http.ResponseWriter, request *http.Request) {
	type errorConfig struct {
		Text string
		Link string
	}
	var config errorConfig
	typeErr := request.URL.Query().Get("type")
	link := request.URL.Query().Get("link")
	if val, ok := errorconfigs[typeErr]; ok {
		config = errorConfig{
			Text: val,
			Link: "https://" + request.Host + link,
		}
	} else {
		config = errorConfig{
			Text: errorconfigs["emptyError"],
			Link: "https://" + request.Host,
		}
	}
	errorRoute.Execute(writer, config)
	return

}
