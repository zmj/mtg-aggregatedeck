package main

import (
	"fmt"
	"go/build"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var static string

func init() {
	var err error
	if static, err = staticDir(); err != nil {
		log.Fatal("could not get static dir: %v", err)
	}
}

func staticDir() (string, error) {
	if fi, err := os.Stat("static"); err == nil && fi.IsDir() {
		return "static", nil
	}
	pkg, err := build.Import("github.com/malthrin/mtg-aggregatedeck", "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return filepath.Join(pkg.Dir, "static"), nil
}

func respond(w http.ResponseWriter, decks []*Deck) {
	deck, err := aggregate(decks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Aggregation error: %s\n", err.Error())
		return
	}
	io.WriteString(w, deck.String())
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, filepath.Join(static, "submit.html"))
		return
	}

	if r.Method == "POST" {
		files, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		decks := make([]*Deck, 0)
		for {
			file, err := files.NextPart()
			if err == io.EOF {
				break
			}
			if file.FileName() == "" {
				continue
			}
			deck, err := NewDeck(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Printf("Deck parse error: %s\n", err.Error())
				return
			}
			decks = append(decks, deck)
		}
		respond(w, decks)
	}
}

func main() {
	http.HandleFunc("/", handle)
	staticPrefix := "/static/"
	http.Handle(staticPrefix, http.StripPrefix(staticPrefix, http.FileServer(http.Dir(static))))
	fmt.Println("Starting server at http://localhost:8981")
	http.ListenAndServe(":8981", nil)
}
