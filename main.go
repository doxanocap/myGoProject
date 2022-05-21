package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	ID                              uint16
	UserName, Title, Tags, FullText string
}

var posts = []Article{}
var showPosts = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM articles")

	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next() {
		var user Article
		err = res.Scan(&user.ID, &user.UserName, &user.Title, &user.Tags, &user.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, user)
	}
	temp.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	temp.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	yourName := r.FormValue("yourName")
	articleTitle := r.FormValue("articleTitle")
	textTags := r.FormValue("textTags")
	full_text := r.FormValue("full_text")

	if yourName == "" || articleTitle == "" || textTags == "" || full_text == "" {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	} else {

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")

		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`name`, `articleTitle`, `textTags`, `full_text`) VALUES('%s','%s','%s','%s')", yourName, articleTitle, textTags, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func errorTemplates(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/error.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Printf(err.Error())
	}
	temp.ExecuteTemplate(w, "error", nil)
}

func show_post(w http.ResponseWriter, r *http.Request) {
	// gorilla mux
	vars := mux.Vars(r)

	temp, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	// entering to data base
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// selecting data
	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE id = '%s'", vars["ID"]))
	if err != nil {
		panic(err)
	}
	showPosts = Article{}
	for res.Next() {
		var user Article
		err = res.Scan(&user.ID, &user.UserName, &user.Title, &user.Tags, &user.FullText)
		if err != nil {
			panic(err)
		}
		showPosts = user
	}

	temp.ExecuteTemplate(w, "show", showPosts)
}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/error", errorTemplates).Methods("GET")
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{ID:[0-9]+}", show_post).Methods("GET")
	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":2004", nil)
}

func main() {
	handleFunc()
}
