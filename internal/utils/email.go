package utils

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strconv"

	// Configuration
	config "github.com/BIQ-Cat/easyserver/config/base"
)

type headers struct {
	From    string
	To      string
	Subject string
}

func (h headers) String() string {
	msg := fmt.Sprintf("From: %s\r\n", h.From)
	msg += fmt.Sprintf("To: %s\r\n", h.To)
	msg += fmt.Sprintf("Subject: %s\r\n", h.Subject)
	return msg
}

func SendEmail(to string, subject string, data interface{}, emailTemp string) error {
	from := config.EnvConfig.EmailFrom
	smtpPass := config.EnvConfig.SMTPPass
	smtpUser := config.EnvConfig.SMTPUser
	smtpHost := config.EnvConfig.SMTPHost
	smtpPort := config.EnvConfig.SMTPPort

	smtpAddr := smtpHost + ":" + strconv.Itoa(smtpPort)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	var body bytes.Buffer

	templ, err := ParseTemplateDir("mail")
	if err != nil {
		return err
	}

	err = templ.ExecuteTemplate(&body, emailTemp, &data)
	if err != nil {
		return err
	}

	msg := buildMessage(headers{
		From:    from,
		To:      to,
		Subject: subject,
	}, body.Bytes())

	return smtp.SendMail(smtpAddr, auth, from, []string{to}, msg)
}

func buildMessage(hdrs headers, body []byte) []byte {
	mimeHeaders := make(map[string]string)

	mimeHeaders["MIME-Version"] = "1.0"
	mimeHeaders["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", "text/html")
	mimeHeaders["Content-Disposition"] = "inline"
	mimeHeaders["Content-Transfer-Encoding"] = "quoted-printable"

	headerMsg := bytes.NewBufferString(hdrs.String())
	for key, value := range mimeHeaders {
		headerMsg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	var encodedBody bytes.Buffer
	w := quotedprintable.NewWriter(&encodedBody)
	_, err := w.Write(body)
	if err != nil {
		panic(err) // Bad argument
	}
	w.Close()

	msg := headerMsg.Bytes()
	msg = append(msg, "\r\n"...)
	msg = append(msg, encodedBody.Bytes()...)
	return msg
}
