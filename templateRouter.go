package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func (h RootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		err := request.ParseForm()
		if err != nil {
			return
		}
		username := request.Form.Get("user")
		password := request.Form.Get("password")
		// cookie authentifizieren checken
		isUser, cookieText := authentifizierung.AuthenticateUser(&username, &password)
		if isUser == true {
			// wenn user authentifiziert ist dann wird cookie erstellt und
			cookie := &http.Cookie{
				Name:     "session-KalenderID",
				Value:    cookieText,
				Path:     "/",
				MaxAge:   3600,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(writer, cookie)
			//redirect to new site
			http.Redirect(writer, request, "https://"+request.Host+"/user/Create", http.StatusFound)
			return
		} else {
			// wenn nicht authentifiziert ist wird weiter geleitet oder bei problemen gibt es ein 500 status
			if len(cookieText) == 0 {
				writer.WriteHeader(500)
			} else {
				http.Redirect(writer, request, "/", http.StatusContinue)
			}
		}
	}
	mainRoute, err := template.ParseFiles("./assets/sites/index.html", "./assets/templates/footer.html")
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}

}

func (c CreatUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	fmt.Fprintf(writer, "<h1>Test</h1>")
}
