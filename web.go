package main

import (
	"fmt"
	"log"
	"net/http"

	nError "github.com/news-ai/web/errors"
)

func verifySMTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(VerifySMTPAccount("mail.privateemail.com:465", "hi@abhi.co", ""))

	switch r.Method {
	case "POST":
		fmt.Fprintf(w, "Hello")
		return
	}

	nError.ReturnError(w, http.StatusInternalServerError, "SMTP error", "method not implemented")
	return
}

func main() {
	http.HandleFunc("/", verifySMTP)         // set router
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
