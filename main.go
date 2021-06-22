package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"html/template"
	"net/http"
	"os"
)

const TemplatesPath = "templates/"
const StaticPath = "static/"

type Article struct {
	ID                    uint
	Title, Brief, Content string
}

func main() {
	handlerFunc()
}

func handlerFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save-article", saveArticle).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", getPost).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/"+StaticPath, http.StripPrefix("/"+StaticPath, http.FileServer(http.Dir(StaticPath))))
	http.ListenAndServe(":"+os.Getenv("APP_PORT"), nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(TemplatesPath+"layout.gohtml", TemplatesPath+"index.gohtml")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM articles")
	if err != nil {
		panic(err)
	}

	var posts []Article
	for res.Next() {
		var post Article
		err = res.Scan(&post.ID, &post.Title, &post.Brief, &post.Content)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "layout", posts)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles(TemplatesPath+"layout.gohtml", TemplatesPath+"post.gohtml")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE id = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	var post Article
	for res.Next() {
		err = res.Scan(&post.ID, &post.Title, &post.Brief, &post.Content)
		if err != nil {
			panic(err)
		}
	}

	t.ExecuteTemplate(w, "layout", post)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(TemplatesPath+"layout.gohtml", TemplatesPath+"create.gohtml")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "layout", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	brief := r.FormValue("brief")
	text := r.FormValue("text")

	if title == "" || brief == "" || text == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
		return
	}

	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO articles (title, brief, content) VALUES ('%s','%s', '%s')", title, brief, text))
	if err != nil {
		panic(err)
	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
