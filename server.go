package main

import "fmt"
import "net/http"
import "io"
import "io/ioutil"
import "strings"

func respond(w http.ResponseWriter, decklists []string) {
	decks := make([]*Deck, len(decklists))
	for i,decklist := range decklists {
		deck, err := NewDeck(strings.Split(decklist, "\n"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Printf("Deck parse error: %s\n", err.Error())
			return
		}
		decks[i] = deck
	}
	result := aggregate(decks)
	io.WriteString(w, result.String())
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/submit.html")
	} else if r.Method == "POST" {
		files, err := r.MultipartReader()		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		decklists := make([]string, 0)
		for {
			file, err := files.NextPart()
			if err == io.EOF {
				break
			}
			if file.FileName() == "" {
				continue
			}			
			content, err := ioutil.ReadAll(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			decklists = append(decklists, string(content))
		}		
		respond(w, decklists)
	} else {
		fmt.Println(r.Method)
	}
}

func main() {
	appPrefix := "metadeck"
	http.HandleFunc(fmt.Sprintf("/%s/", appPrefix), handle)	
	staticPath := fmt.Sprintf("/%s/static/", appPrefix)
	http.Handle(staticPath, http.StripPrefix(staticPath, http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8981", nil)
}