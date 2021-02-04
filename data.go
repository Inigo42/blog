package main

import (
	"bytes"
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

var err error

type Article struct {
	Title   string   `json:"title"`
	Author  string   `json:author`
	Content []string `json:content`
}

func addArticle(db *leveldb.DB, article Article, key string) {
	testArticle := article
	superBytes := new(bytes.Buffer)
	json.NewEncoder(superBytes).Encode(testArticle)
	db.Put([]byte(key), superBytes.Bytes(), nil)
}

func fetchArticle(db *leveldb.DB, key string, article *Article) {
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		panic("OH FUCKAROOONIE DOOO, we were not able to fetch the article")
	}
	json.Unmarshal(data, &article)
}

//example of how to do stuff
/*
func main() {
	db, err := leveldb.OpenFile("farticles", nil)
	if err != nil {
		panic("fuck")
	}
	defer db.Close()

	addArticle(db, Article{"superduper title", "Jeremi Forman-Duranona",
		[]string{"Bravely Bold Sir Robin, rode forth from Camelot."}}, "lezfuckin_go")
	var article Article
	fetchArticle(db, "lezfuckin_go", &article)
	fmt.Println(article.Content[0])
}
*/
