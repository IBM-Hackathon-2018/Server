package main

import (
	"net/http"
	"flag"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/mux"
	"log"
	"golang.org/x/crypto/bcrypt"
)

var (
	addr = flag.String("addr", "127.0.0.1:8000", "http service address")
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))
	)


func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/"
	//Todo querry database
	if name == "a" && pass == "a" {
		setSession(name, response)
		redirectTarget = "/internal"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func registerHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	key := request.FormValue("secret_key")
	redirectTarget := "/register"
	if name != "" && pass != "" && key =="secret_key"{
		setSession(name, response)
		redirectTarget = "/"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		http.ServeFile(response, request, "web/internal.html")
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "web/login_screen.html")
}

func uploadHandler(response http.ResponseWriter, request *http.Request) {
	//Todo upload file or hash somewhere
}


func main() {
	bitSignApi := newBitSignApi()
	r := mux.NewRouter()
	r.HandleFunc("/", indexPageHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/logout", logoutHandler)

	r.HandleFunc("/internal", internalPageHandler)
	r.HandleFunc("/upload", uploadHandler)

	//API
	r.HandleFunc("/documents", bitSignApi.listDocuments).Methods("GET")
	r.HandleFunc("/documents", bitSignApi.uploadDocuments).Methods("POST")
	r.HandleFunc("/documents/{id}", bitSignApi.showDocument).Methods("GET")
	r.HandleFunc("/documents/{id}/persons/{personID}", bitSignApi.addSignee).Methods("PUT")
	r.HandleFunc("/documents/{id}/persons/{personID}", bitSignApi.removeSignee).Methods("Delete")
	r.HandleFunc("/documents/{id}/sign", bitSignApi.signDocument).Methods("PUT")

	r.HandleFunc("/register", uploadHandler)

	log.Fatal(http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", r))


}