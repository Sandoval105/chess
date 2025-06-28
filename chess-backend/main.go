package main

import (
	"log"
	"net/http"
)

func main() {

	InitDB()

	//endpoints

	http.Handle("/", http.FileServer(http.Dir("../frontend/public")))

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/ws", wsHandler)

	log.Println("servidor rodando em  http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
