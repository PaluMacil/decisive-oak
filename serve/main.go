package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("serving from:", dir)
	fs := http.FileServer(http.Dir("./serve/www"))
	fsTree := http.FileServer(http.Dir("./out"))
	http.HandleFunc("/api/list/files", func(w http.ResponseWriter, r *http.Request) {
		var treeItems []TreeItem
		files, _ := filepath.Glob("out/*.data.tree.json")
		for _, filename := range files {
			item := TreeItem{
				Filename: "tree/" + filepath.Base(filename),
			}
			treeItems = append(treeItems, item)
		}
		json.NewEncoder(w).Encode(treeItems)
	})
	http.Handle("/", fs)
	http.Handle("/tree/", http.StripPrefix("/tree/", fsTree))

	log.Println("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type TreeItem struct {
	Filename string
}
