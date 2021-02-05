package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	//"github.com/gorilla/mux"
	"html/template"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
)

//----------------database and whatnot------------------
type Payload struct {
	Stuff [3]string
	Link  string
}

var err error

type Article struct {
	Title   string   `json:"title"`
	Index   int      `json:index`
	Author  string   `json:author`
	Content []string `json:content`
}

func addArticle(db *leveldb.DB, article Article, key string) {
	testArticle := article
	superBytes := new(bytes.Buffer)
	json.NewEncoder(superBytes).Encode(testArticle)
	db.Put([]byte(key), superBytes.Bytes(), nil)
}

func deleteArticle(db *leveldb.DB, key string) {
	db.Delete([]byte(key), nil)
}

func fetchArticle(db *leveldb.DB, key string, article *Article) {
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		panic("OH FUCKAROOONIE DOOO, we were not able to fetch the article")
	}
	json.Unmarshal(data, &article)
}

//-----------------http-----------------------------------------
var templates *template.Template

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
	db, err := leveldb.OpenFile("farticles", nil)
	if err != nil {
		panic("fuck")
	}
	defer db.Close()

	var article Article
	fetchArticle(db, "2", &article)
	templates.ExecuteTemplate(w, "article.html", article)
}

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	fmt.Println("THIS WILL BE AN AWESOME BLOG!!!")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/testing", testingHandler)
	http.HandleFunc("/test_article", basicArticleHandler)
	http.ListenAndServe(":8080", nil)
}
