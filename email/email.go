package email

import (
	"fmt"
	"bytes"
	"text/template"
	"errors"
	"net/smtp"
        "github.com/gagiopapinni/price-tracker-api-exercise/helper"
)


var(
	from string
	password string
	smtpHost string = "smtp.gmail.com"
	smtpPort string = "587"
)

type ConfirmationPool struct {
	Keys map[string]string
} 

func (cp *ConfirmationPool) GenerateKey(email string) string {
	if cp.Keys==nil {
		cp.Keys = make(map[string]string)
	}
	key := helper.RandSeq(10)
	cp.Keys[email] = key
	return key
}

func (cp *ConfirmationPool) Confirm(email, key string) bool {
	if cp.Keys == nil { return false }
	if k, ok := cp.Keys[email]; ok && k==key {
		delete(cp.Keys, email)
		return true
	}
	return false
}

func Configure(From, Password string) {
	from = From
	password = Password
}

func Send(msg string, to []string) error {
	if from=="" || password=="" {
		return errors.New("not configured")
	}

  	auth := smtp.PlainAuth("", from, password, smtpHost)
	t, _ := template.ParseFiles("templates/Email.html")
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Avito Price Change \n%s\n\n", mimeHeaders)))
	t.Execute(&body, struct {
		Msg    string
	}{
		Msg:   msg,
	})
  	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
  	if err != nil {
   		 return err
  	}

	return nil

}

























