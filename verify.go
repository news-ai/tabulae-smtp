package main

import (
	"github.com/news-ai/web/emails"
)

func VerifySMTPAccount(servername string, email string, password string) error {
	return emails.VerifySMTP(servername, email, password)
}
