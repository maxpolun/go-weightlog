package main

import (
	"./sessions"
	"./users"
	"./util"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func setJSON(rw http.ResponseWriter) {
	header := rw.Header()
	header.Set("Content-Type", "application/json")

}

func rootHandler(rw http.ResponseWriter, req *http.Request) {
	setJSON(rw)
	default_resp := []byte(`{login:'/login'}`)
	rw.Write(default_resp)
}
func loginHandler(rw http.ResponseWriter, req *http.Request) {
	setJSON(rw)
	email := req.FormValue("email")
	pw := req.FormValue("password")
	if len(email) == 0 || len(pw) == 0 {
		http.Error(rw, "{errors:[\"no_username_password\"]}", 400)
		return
	}
	db := util.GetDB()
	defer db.Close()
	u, err := users.GetByEmail(db, email)
	if err != nil {
		log.Printf("error while looking up user: %v", err)
		http.Error(rw, "{errors:[\"bad_login\"]}", 400)
		return
	}
	if !u.Verify(pw) {
		http.Error(rw, "{errors:[\"bad_login\"]}", 400)
		return
	}
	s, _ := sessions.New(u.Id)
	err = s.Save(db)
	if err != nil {
		log.Printf("error while saving session in database: %v", err)
		http.Error(rw, "{errors:[\"internal_error\"]}", 500)
		return
	}
	rw.Write([]byte(fmt.Sprintf(`{errors:[], session_id:"%s"}`, s.Id)))
}
func registerHandler(rw http.ResponseWriter, req *http.Request) {
	setJSON(rw)
	email := req.FormValue("email")
	pw := req.FormValue("password")
	pw_conf := req.FormValue("password_confirmation")
	if email == "" || pw == "" || pw_conf == "" {
		http.Error(rw, "{errors:[\"no_email_password\"]}", 400)
		return
	}
	if pw != pw_conf {
		http.Error(rw, "{errors:[\"pw_conf_bad_match\"]}", 400)
		return
	}
	user, err := users.New(email, pw)
	if err != nil {
		http.Error(rw, "{errors:[\"hash_error\"]}", 500)
		return
	}
	db := util.GetDB()
	defer db.Close()
	err = user.Save(db)
	if err != nil {
		http.Error(rw, "{errors:[\"db_error\"]}", 500)
		return
	}
	rw.Write([]byte(fmt.Sprintf(`{errors:[], user:{id:%v, email:%v}}`, user.Id, user.Email)))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET").Name("root")
	router.HandleFunc("/login", loginHandler).Methods("POST").Name("login")
	router.HandleFunc("/register", registerHandler).Methods("POST").Name("register")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	} else {
		port = ":" + port
	}
	log.Printf("starting api server on port %v", port)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(port, nil))
}
