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

type SMTPResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func sendSMTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		buf, _ := ioutil.ReadAll(r.Body)

		decoder := ffjson.NewDecoder()
		var emailSettings models.SMTPEmailSettings
		err := decoder.Decode(buf, &emailSettings)
		if err != nil {
			nError.ReturnError(w, http.StatusInternalServerError, "SMTP verify error", err.Error())
			return
		}

		userPassword, err := encrypt.DecryptString(emailSettings.EmailPassword)
		if err != nil {
			nError.ReturnError(w, http.StatusInternalServerError, "SMTP verify error", err.Error())
			return
		}

		smtpError := emails.SendSMTPEmail(emailSettings.Servername, emailSettings.EmailUser, userPassword, emailSettings.To, emailSettings.Subject, emailSettings.Body)

		response := SMTPResponse{}
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
			nError.ReturnError(w, http.StatusInternalServerError, "SMTP verify error", err.Error())
		}

		return
	}

	nError.ReturnError(w, http.StatusInternalServerError, "SMTP verify error", "method not implemented")
	return
}

func verifySMTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		buf, _ := ioutil.ReadAll(r.Body)

		decoder := ffjson.NewDecoder()
		var emailSettings models.SMTPSettings
		err := decoder.Decode(buf, &emailSettings)
		if err != nil {
			nError.ReturnError(w, http.StatusInternalServerError, "SMTP send error", err.Error())
			return
		}

		userPassword, err := encrypt.DecryptString(emailSettings.EmailPassword)
		if err != nil {
			nError.ReturnError(w, http.StatusInternalServerError, "SMTP send error", err.Error())
			return
		}

		smtpError := emails.VerifySMTP(emailSettings.Servername, emailSettings.EmailUser, userPassword)

		response := SMTPResponse{}
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
			nError.ReturnError(w, http.StatusInternalServerError, "SMTP send error", err.Error())
		}

		return
	}
	nError.ReturnError(w, http.StatusInternalServerError, "SMTP send error", "method not implemented")
	return
}

func main() {
	http.HandleFunc("/send", sendSMTP)
	http.HandleFunc("/verify", verifySMTP)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
