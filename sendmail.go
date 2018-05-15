package memogo

/*
  Sending emails with attachments realization
*/

import (
	"log"
	"net/mail"
	"net/smtp"
)

//EmailCredentials description of server and account description
type EmailCredentials struct {
	Username, Password, Server, From, FromName string
	Port                                       int
	UseTLS                                     bool
}

func SendEmailMsg(authCreds EmailCredentials, msg *Message) error {

	//sender information
	msg.From = mail.Address{Name: authCreds.FromName, Address: authCreds.From}

	//do sendmail
	auth := smtp.PlainAuth("", authCreds.Username, authCreds.Password, authCreds.Server)

	//sending without TLS
	if !authCreds.UseTLS {
		if err := SendMail(authCreds.Server, uint(authCreds.Port), auth, msg); err != nil {
			log.Println("SendEmailMsg error:", err)
			return err
		}
		return nil
	}
	//sending with TLS
	if err := SendMailSSL(authCreds.Server, uint(authCreds.Port), auth, msg); err != nil {
		log.Println("SendEmailMsgSSL error:", err)
		return err
	}
	return nil
}
