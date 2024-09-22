package services

import "net/smtp"

func SendWarning(email string) error {
	auth := smtp.PlainAuth(
		"",
		"obmaza.aaa212@gmail.com",
		"dxkc mpdk srcq veew",
		"smtp.gmail.com",
	)

	msg := "Subject: Warning! Someone used your refresh token!\nIf it was you, don't reply to this message."

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"obmaza.aaa212@gmail.com",
		[]string{email},
		[]byte(msg),
	)
	if err != nil {
		return err
	}

	return nil
}
