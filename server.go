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
	Stuff        []string
	Link         string
	ArticleLinks []string
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
	//words := []string{"first title", "second title", "third title"}
	link2xkcd := "https://imgs.xkcd.com/comics/hug_count.png"
	db, err := leveldb.OpenFile("farticles", nil)
	if err != nil {
		panic("fuck")
	}
	defer db.Close()
	var keys []string
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		keys = append(keys, string(iter.Key()))
	}
	var titles []string
	var article Article
	for _, key := range keys {
		fetchArticle(db, key, &article)
		title := article.Title
		titles = append(titles, title)
	}
	payload := Payload{Stuff: titles, Link: link2xkcd, ArticleLinks: keys}
	templates.ExecuteTemplate(w, "index.html", payload)
}

/*
func testingHandler(w http.ResponseWriter, r *http.Request) {
	stuff := []string{"blah1", "blah2", "blah3"}
	link2xkcd := "https://imgs.xkcd.com/comics/hug_count.png"

	payload := Payload{Stuff: stuff, Link: link2xkcd}
	templates.ExecuteTemplate(w, "testing.html", payload)
}
*/
func basicArticleHandler(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Path[1:]

	db, err := leveldb.OpenFile("farticles", nil)
	if err != nil {
		panic("fuck")
	}
	defer db.Close()

	var article Article
	fetchArticle(db, k, &article)
	templates.ExecuteTemplate(w, "article.html", article)
}

func main() {
	/*
		db, err := leveldb.OpenFile("farticles", nil)
		if err != nil {
			panic("fuck")
		}
		defer db.Close()
		var keys []string
		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			keys = append(keys, string(iter.Key()))
		}
		for _, key := range keys {
			http.HandleFunc(key, basicArticleHandler)
		}
	*/
	templates = template.Must(template.ParseGlob("templates/*.html"))
	fmt.Println("THIS WILL BE AN AWESOME BLOG!!!")

	http.HandleFunc("/", indexHandler)
	//http.HandleFunc("/testing", testingHandler)
	http.HandleFunc("/test_article", basicArticleHandler)
	http.ListenAndServe(":8080", nil)
}
