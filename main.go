package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Article struct {
	Id        uint16
	Title     string
	Anons     string
	Full_Text string
}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "rotroturan:11111111@tcp(db4free.net:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query("SELECT * FROM `article`")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_Text)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
		//fmt.Println(fmt.Sprintf("Post: %s with id: %d", post.Title, post.Id))
	}

	t.ExecuteTemplate(w, "index", posts)
}
func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}
func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := sql.Open("mysql", "rotroturan:11111111@tcp(db4free.net:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `article` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_Text)
		if err != nil {
			panic(err)
		}
		showPost = post
		//fmt.Println(fmt.Sprintf("Post: %s with id: %d", post.Title, post.Id))
	}
	t, err := template.ParseFiles("templates/show_post.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "shows", showPost)

}
func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")
	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Not all data complete")
	} else {
		db, err := sql.Open("mysql", "rotroturan:11111111@tcp(db4free.net:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		// Set data

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `article` (`title`, `anons`, `full_text`) VALUES ('%s','%s', '%s')", title, anons, full_text))
		if err != nil {
			fmt.Println(err)
		}
		defer insert.Close()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}
func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	rtr.HandleFunc("/save_article", save_article).Methods("POST")

	http.Handle("/", rtr)
	rtr.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8080", nil)
}
func main() {

	handleFunc()
}
