package emailService

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"

	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
)

// SendEmail ...
func SendEmail(emailType string, user *map[string]*models.User, multiple bool) error {
	switch emailType {
	case "approval":
		mail, err := CreateHtmlMessageApproval(user)
		if err == nil {
			go SendBatchMail(mail, multiple, "approval")
			return nil
		}
	default:
		return nil
	}
	return nil
}

// SendBatchMail ...
func SendBatchMail(emails *map[string]string, multiple bool, emailType string) {
	if multiple == false {
		subject := ""
		if emailType == "approval" {
			subject = "Accept the Invitation"
		}

		go func(subject string, emails map[string]string) {
			for user, mail := range emails {
				e := &email.Email{
					To:      []string{user},
					From:    "Pugcha Backend <pugcha.backend@gmail.com>",
					Subject: subject,
					HTML:    []byte(mail),
					Headers: textproto.MIMEHeader{},
				}
				servername := "smtp.gmail.com:465"
				host, _, _ := net.SplitHostPort(servername)

				tsconfig := &tls.Config{InsecureSkipVerify: true, ServerName: host}

				auth := smtp.PlainAuth("Pugcha Backend", config.Configuration.GmailUser, config.Configuration.GmailPassword, host)

				err := e.SendWithTLS(servername, auth, tsconfig)

				if err != nil {
					fmt.Println(err, "here")
				}
			}
		}(subject, *emails)
	}
}
