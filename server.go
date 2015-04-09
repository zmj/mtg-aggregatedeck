package main

import "fmt"
import "net/http"
//import "net/url"

func handle(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		http.ServeFile(w, r, "static/submit.html")
	} else {
		fmt.Println("post")
	}
}

func main() {
	fmt.Printf("Hello world\n")
	appPrefix := "metadeck"
	http.HandleFunc(fmt.Sprintf("/%s/", appPrefix), handle)	
	staticPath := fmt.Sprintf("/%s/static/", appPrefix)
	http.Handle(staticPath, http.StripPrefix(staticPath, http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8981", nil)
}