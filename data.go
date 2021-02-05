package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

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

func generate_Article() {
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
	fmt.Print(keys)
	var greatestIndex int64 = 0
	for _, val := range keys {
		valrev, _ := strconv.ParseInt(val, 10, 64)
		if valrev > greatestIndex {
			greatestIndex = valrev
		}
	}

	title := "Learning Go by Writing the Code for this Blog"
	author := "Jeremi Forman-Duranona"
	content := []string{`Well, this has been a bit of a journey. It started with me wanting to both learn 
	a language more powerful than python but not quite as low level as C++ or beurocratic as Java, and make a website from the bottom up. 
	I decided on Go, and I figured I'd make a tech blog while I'm at it. Not to mention I can add features any time I want.`,
		`I spent a couple hours absorbing as much Golang content on youtube to familiarize myself with the syntax and customs. Next, I 
	started messing around with the language itself as I discovered interesting things -- e.g. in Go there are only for loops, no such thing
	as a while loop. Though this surprised me, I came to appreciate it`}
	newIndex := int(greatestIndex) + 1
	addArticle(db, Article{title, newIndex, author, content}, strconv.Itoa(newIndex))
}

func main() {
	generate_Article()
}

//example of how to do stuff
/*
func main() {
	db, err := leveldb.OpenFile("farticles", nil)
	if err != nil {
		panic("fuck")
	}
	defer db.Close()

	addArticle(db, Article{"superduper title", 1, "Jeremi Forman-Duranona",
		[]string{"Bravely Bold Sir Robin, rode forth from Camelot."}}, "1")

	deleteArticle(db, "2")
	var article Article
	fetchArticle(db, "1", &article)
	fmt.Println(article.Title)

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		fmt.Println(string(key))
	}
}
*/
