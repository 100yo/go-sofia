package main

import "log"
import "net/http"
	
func main() {

	log.Print("Hello, World")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal()
	}
}
