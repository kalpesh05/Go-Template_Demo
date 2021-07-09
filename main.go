package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var templates = template.Must(template.ParseFiles("template/partial/header.html", "template/partial/footer.html", "template/login.html", "template/index-logged-out.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

//A Page structure
type Page struct {
	Title        string
	ErrorMessage string
}

func Index(w http.ResponseWriter, r *http.Request) {
	// data := struct {
	// 	Title  string
	// 	Header string
	// }{
	// 	Title:  "Index Page",
	// 	Header: "Hello, World!",
	// }

	// if err := templates.ExecuteTemplate(w, "indexT.html", data); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	display(w, "index-logged-out.html", &Page{Title: "Home", ErrorMessage: ""})
}

func MyAccount(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	ErrorMessage := ""
	if r.Method == "POST" {
		via := r.FormValue("via")
		log.Println(via)
		if via == "login" {
			email := r.FormValue("email")
			password := r.FormValue("password")
			log.Println("UPDATE: email: " + email + "password" + password)
		} else if via == "register" {
			usertype := r.FormValue("account-type-radio")
			email := r.FormValue("email")
			password := r.FormValue("password")
			cpassword := r.FormValue("password-repeat-register")

			log.Println("UPDATE: usertype: " + usertype + "password" + password + "cpwd " + cpassword + " email " + email)
			if cpassword == password {
				insForm, err := db.Prepare("INSERT INTO users( email, password, type) VALUES(?,?,?)")
				if err != nil {
					panic(err.Error())
				}
				insForm.Exec(email, password, usertype)
			} else {
				ErrorMessage := "Confirm password should be match with password*"
				log.Println("UPDATE: NAME: " + ErrorMessage)
			}

		}
		defer db.Close()

	}

	display(w, "login.html", &Page{Title: "Login/Register", ErrorMessage: ErrorMessage})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/login", MyAccount)

	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("template/css/"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("template/images/"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("template/js/"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("template/fonts/"))))

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	http.Handle("/", r)
	log.Fatalln(http.ListenAndServe(":9000", nil))
}
