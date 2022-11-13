package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"DHBW_GO_Projekt/kalenderansicht"
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
				Name:     "Session Kalender-ID",
				Value:    cookieText,
				MaxAge:   300,
				Secure:   true,
				HttpOnly: false,
			}
			request.AddCookie(cookie)
			http.Redirect(writer, request, "kalender/tabellenAnsicht", http.StatusContinue)
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

func (kalender TabellenHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	kalenderansicht.TabellenHandler(writer, request)
}
