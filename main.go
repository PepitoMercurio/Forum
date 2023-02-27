package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type zer struct {
}

var id int
var username string
var mail string
var password string
var zerr zer

func main() {
	//connecter = false
	database, _ := sql.Open("sqlite3", "./nraboy.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, mail TEXT, password TEXT)")
	statement.Exec()
	tmpl := template.Must(template.ParseGlob("html/*"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		mdp := r.FormValue("password")
		user := r.FormValue("username")
		adresse := r.FormValue("mail")

		bim, _ := database.Query("SELECT id, username, mail, password FROM people")
		for bim.Next() {
			bim.Scan(&id, &username, &mail, &password)
			if mdp != "" && user != "" {

				if mdp == password && user == username && adresse == mail {
					http.Redirect(w, r, "http://localhost:5550/vitrine", http.StatusSeeOther)
					fmt.Println("Bienvenue")
				} else {
					fmt.Printf("non")
				}

			}

		}
		tmpl.ExecuteTemplate(w, "login.html", zerr)
	})
	http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "page.html", zerr)
	})
	http.HandleFunc("/topic", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "topic.html", zerr)
	})
	http.HandleFunc("/vitrine", func(w http.ResponseWriter, z *http.Request) {
		tmpl.ExecuteTemplate(w, "vitrine.html", zerr)
	})

	http.HandleFunc("/posts", func(w http.ResponseWriter, z *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", zerr)
		mdp1 := z.FormValue("newpassword")
		adresse1 := z.FormValue("newmail")
		user1 := z.FormValue("newusername")
		if mdp1 != "" && user1 != "" && adresse1 != "" {
			username += user1
			password += mdp1
			mail += adresse1
			statement, _ = database.Prepare("INSERT INTO people (username, mail, password) VALUES (?, ?, ?)")
			statement.Exec(user1, adresse1, mdp1)
			rows, _ := database.Query("SELECT id, username, mail, password FROM people")
			tmpl.ExecuteTemplate(w, "pop2.html", zerr)
			for rows.Next() {
				rows.Scan(&id, &username, &mail, &password)
				fmt.Println(strconv.Itoa(id) + ": " + username + " " + mail + " " + password)
			}
		}

	})

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe("localhost:5550", nil)
}
