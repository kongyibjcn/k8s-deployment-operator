package main

import (
	"fmt"
	"log"
	"net/http"
)

func SayHello(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the Green Docker!")
	//fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", SayHello)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}