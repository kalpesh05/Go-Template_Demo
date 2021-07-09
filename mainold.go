package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Index(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {

		index, _ := template.ParseFiles("template/login.html")
		index.Execute(w, nil)
		//	http.ServeFile(w, r, "hireo/index.html")
		return
	}
	// db := dbConn()
	// selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// emp := Employee{}
	// res := []Employee{}
	// for selDB.Next() {
	// 	var id int
	// 	var name, city string
	// 	err = selDB.Scan(&id, &name, &city)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	emp.Id = id
	// 	emp.Name = name
	// 	emp.City = city
	// 	res = append(res, emp)
	// }

	// defer db.Close()
}

func BrowserJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {

		browserjob, _ := template.ParseFiles("hireo/jobs-list-layout-1.html")
		browserjob.Execute(w, nil)
		//	http.ServeFile(w, r, "hireo/index.html")
		return
	}

}

func main() {
	log.Println("Server started on: http://localhost:8080")
	t := template.Must(template.ParseFiles("hireo/header.tmpl"))
	t.Execute(os.Stdout, nil)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("hireo/css"))))

	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("hireo/js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("hireo/fonts"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("hireo/images"))))
	http.HandleFunc("/", Index)
	http.HandleFunc("/BrowserJobs", BrowserJob)
	http.ListenAndServe(":8080", nil)
}
