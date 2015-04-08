package main

import "fmt"
import "net/http"

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("req");
}

func main() {
	fmt.Printf("Hello world\n")
	http.HandleFunc("/aggregatedeck", handle)
	http.ListenAndServe(":8981", nil)
}