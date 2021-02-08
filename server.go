package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	//"github.com/gorilla/mux"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/syndtr/goleveldb/leveldb"
)

//----------------database and whatnot------------------
type Payload struct {
	Farticles    []farticle
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

func reverse(a []string) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

type farticle struct {
	Title string
	Key   int
}

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
	reverse(titles)
	reverse(keys)
	var k []int
	var x int
	for _, key := range keys {
		x, _ = strconv.Atoi(key)
		k = append(k, x-1)
	}
	var farticles []farticle
	for i, title := range titles {
		farticles = append(farticles, farticle{title, k[i]})
	}
	payload := Payload{Farticles: farticles, Link: link2xkcd, ArticleLinks: keys}
	templates.ExecuteTemplate(w, "index.html", payload)
}

func getKeys() []string {
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
	return keys
}

func main() {
	r := chi.NewRouter()

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "id")
		idNumber, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "that's not a number, dummy: "+idString, http.StatusBadRequest)
			return
		}
		db, err := leveldb.OpenFile("farticles", nil)
		if err != nil {
			panic("fuck")
		}
		defer db.Close()
		idString = strconv.Itoa(idNumber + 1)
		var article Article
		fetchArticle(db, idString, &article)
		templates.ExecuteTemplate(w, "article.html", article)
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	fmt.Println("THIS WILL BE AN AWESOME BLOG!!!")

	r.Get("/", indexHandler)
	http.ListenAndServe(":8080", r)
}
