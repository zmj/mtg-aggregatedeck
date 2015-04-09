package main

import "fmt"
import "net/http"
import "io"
import "io/ioutil"
//import "net/url"
//import "mime/multipart"

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/submit.html")
	} else if r.Method == "POST" {
		files, err := r.MultipartReader()		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
				fmt.Println("err reading file %s", err.Error())
				continue
			}
			fmt.Println(file.FileName())
			fmt.Println(string(content))
		}
	} else {
		fmt.Println(r.Method)
	}
}

func main() {
	fmt.Println("Hello world")
	appPrefix := "metadeck"
	http.HandleFunc(fmt.Sprintf("/%s/", appPrefix), handle)	
	staticPath := fmt.Sprintf("/%s/static/", appPrefix)
	http.Handle(staticPath, http.StripPrefix(staticPath, http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8981", nil)
}