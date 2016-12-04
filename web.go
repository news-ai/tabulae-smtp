package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/news-ai/web/emails"
	"github.com/news-ai/web/encrypt"
	nError "github.com/news-ai/web/errors"

	"github.com/news-ai/tabulae/models"
)

type SMTPResonse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func verifySMTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		buf, _ := ioutil.ReadAll(r.Body)

		decoder := ffjson.NewDecoder()
		var emailSettings models.SMTPSettings
		err := decoder.Decode(buf, &emailSettings)
		if err != nil {
			nError.ReturnError(w, http.StatusInternalServerError, "SMTP error", err.Error())
			return
		}

		userPassword, err := encrypt.DecryptString(emailSettings.EmailPassword)
		if err != nil {
			nError.ReturnError(w, http.StatusInternalServerError, "SMTP error", err.Error())
			return
		}

		response := SMTPResonse{}

		smtpError := emails.VerifySMTP(emailSettings.Servername, emailSettings.EmailUser, userPassword)

		if smtpError == nil {
			response.Status = true
		} else {
			response.Status = false
			response.Error = smtpError.Error()
		}

		if err == nil {
			err = ffjson.NewEncoder(w).Encode(response)
		}

		if err != nil {
			nError.ReturnError(w, http.StatusInternalServerError, "Publication handling error", err.Error())
		}

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
