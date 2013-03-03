package main

import (
	"./sessions"
	"./users"
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"log"
	"net/http"
	"os"
)

var db_url string

func init() {
	db_url = os.Getenv("DATABASE_URL")
	if db_url == "" {
		db_url = "user=weightlog dbname=weightlog password=weightlog sslmode=disable"
	}
}

func getDB() *sql.DB {
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		panic(err)
	}
	return db
}

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
	db := getDB()
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
	db := getDB()
	defer db.Close()
	err = user.Save(db)
	if err != nil {
		http.Error(rw, "{errors:[\"db_error\"]}", 500)
		return
	}
	rw.Write([]byte(fmt.Sprintf(`{errors:[], user:{id:%v, email:%v}}`, user.Id, user.Email)))
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	} else {
		port = ":" + port
	}
	log.Printf("starting api server on port %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
