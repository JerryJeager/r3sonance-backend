
	package emails

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(receiverEmail, subject, body string) error {
	email := os.Getenv("EMAIL")
	emailUsername := os.Getenv("EMAILUSERNAME")
	emailPassword := os.Getenv("EMAILPASSWORD")
	m := gomail.NewMessage()
	m.SetAddressHeader("From", email, "r3sonance-backend")
	m.SetHeader("To", receiverEmail)
	m.SetAddressHeader("Cc", receiverEmail, "r3sonance-backend")
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 465, emailUsername, emailPassword)

	// Send the email to user
	if err := d.DialAndSend(m); err != nil {
		// log.Print(err)
		return err
	}
	return nil
}

	