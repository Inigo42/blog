package main

import (
	"fmt"

	//"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var templates *template.Template

type Payload struct {
	Stuff [3]string
	Link  string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	words := [3]string{"first title", "second title", "third title"}
	link2xkcd := "https://imgs.xkcd.com/comics/hug_count.png"
	payload := Payload{Stuff: words, Link: link2xkcd}
	templates.ExecuteTemplate(w, "index.html", payload)
}

func testingHandler(w http.ResponseWriter, r *http.Request) {
	stuff := [3]string{"blah1", "blah2", "blah3"}
	link2xkcd := "https://imgs.xkcd.com/comics/hug_count.png"

	payload := Payload{Stuff: stuff, Link: link2xkcd}
	templates.ExecuteTemplate(w, "testing.html", payload)
}

func basicArticleHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "article.html", nil)
}

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	fmt.Println("THIS WILL BE AN AWESOME BLOG!!!")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/testing", testingHandler)
	http.HandleFunc("/test_article", basicArticleHandler)
	http.ListenAndServe(":8080", nil)
}
