package main

import "fmt"
import "net/http"
import "io"

func respond(w http.ResponseWriter, decks []*Deck) {
	deck,err := aggregate(decks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Aggregation error: %s\n", err.Error())
		return
	}
	io.WriteString(w, deck.String())
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
	appPrefix := "metadeck"
	http.HandleFunc(fmt.Sprintf("/%s/", appPrefix), handle)	
	staticPath := fmt.Sprintf("/%s/static/", appPrefix)
	http.Handle(staticPath, http.StripPrefix(staticPath, http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8981", nil)
}