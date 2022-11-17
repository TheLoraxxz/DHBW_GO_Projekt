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
				Name:     "SessionID-Kalender",
				Value:    cookieText,
				Path:     "/",
				MaxAge:   3600,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(writer, cookie)
			//redirect to new site
			http.Redirect(writer, request, "https://"+request.Host+"/user/create", http.StatusFound)
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

func (createUser CreatUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// get cookie
	cookie, err := request.Cookie("SessionID-Kalender")
	//if cookie is not existing it returns back to the host
	if err != nil {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	//if it is not allowed then continue with normal website else redirect to root
	isAllowed, _ := authentifizierung.CheckCookie(&cookie.Value)
	if isAllowed {
		if request.Method == "POST" {
			err := request.ParseForm()
			if err != nil {
				return
			}
			user := request.Form.Get("newUsername")
			password := request.Form.Get("newPassword")
			err = authentifizierung.CreateUser(&user, &password)
			if err != nil {
				return
			}
			http.Redirect(writer, request, "https://"+request.Host+"/user", http.StatusContinue)

		}
		mainRoute, err := template.ParseFiles("./assets/sites/user-create.html", "./assets/templates/footer.html", "./assets/templates/header.html")
		if err != nil {
			log.Fatal("Coudnt export Parsefiles")
			return
		}
		err = mainRoute.Execute(writer, nil)
		if err != nil {
			log.Fatal("Coudnt Execute Parsefiles")
			return
		}
	} else {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
	}

}

func (changeUser ChangeUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("SessionID-Kalender")
	//if cookie is not existing it returns back to the host
	if err != nil {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	//if it is not allowed then continue with normal website else redirect to root
	isAllowed, user := authentifizierung.CheckCookie(&cookie.Value)
	if isAllowed {
		if request.Method == "POST" {
			//if post request it actually parses the form and trys to change the password and create a new cookie
			err := request.ParseForm()
			if err != nil {
				return
			}
			password := request.Form.Get("oldPassword")
			newPassword := request.Form.Get("newPassword")
			cookies, err := authentifizierung.ChangeUser(&user, &password, &newPassword)
			if err != nil {
				return
			}
			cookie := &http.Cookie{
				Name:     "SessionID-Kalender",
				Value:    cookies,
				Path:     "/",
				MaxAge:   3600,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(writer, cookie)
			http.Redirect(writer, request, "https://"+request.Host+"/user", http.StatusContinue)

			return

		}
		mainRoute, err := template.ParseFiles("./assets/sites/user-change.html", "./assets/templates/footer.html", "./assets/templates/header.html")
		if err != nil {
			log.Fatal("Coudnt export Parsefiles")
			return
		}
		err = mainRoute.Execute(writer, nil)
		if err != nil {
			log.Fatal("Coudnt Execute Parsefiles")
			return
		}
	} else {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
	}
}

func (user UserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("SessionID-Kalender")
	//if cookie is not existing it returns back to the host
	if err != nil {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	//if it is not allowed then continue with normal website else redirect to root
	isAllowed, _ := authentifizierung.CheckCookie(&cookie.Value)
	if isAllowed {
		mainRoute, err := template.ParseFiles("./assets/sites/user.html", "./assets/templates/footer.html", "./assets/templates/header.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = mainRoute.Execute(writer, nil)
		if err != nil {
			log.Fatal("Coudnt Execute Parsefiles")
			return
		}
	} else {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
	}
}
