package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
)

var (
	smtpHost       = ""
	smtpPort int64 = 0
	smtpUser       = ""
	smtpPass       = ""
)

type EmailConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

type Email struct {
	Subject string
	From    string
	To      []string
	Body    string
}
type Body struct {
	Purchases []Puchase44
}

func SendEmail(e Email) error {

	emailConf := &EmailConfig{smtpUser, smtpPass, smtpHost, int(smtpPort)}

	emailauth := smtp.PlainAuth("", emailConf.Username, emailConf.Password, emailConf.Host)
	adr := mail.Address{Name: "ООО Какая-то организация", Address: e.From}
	headerFrom := adr.String()
	sender := e.From

	receivers := e.To
	e.Body = "From: " + headerFrom + "\r\n" +
		"To: " + e.To[0] + "\r\n" +
		"MIME-Version: 1.0" + "\r\n" +
		"Content-type: text/html ; charset=\"utf-8\"" + "\r\n" +
		"Subject: " + e.Subject + "\r\n\r\n" +
		e.Body + "\r\n"
	message := []byte(e.Body)

	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(emailConf.Port),
		emailauth,
		sender,
		receivers,
		message,
	)

	if err != nil {
		return err
	}
	return nil
}

func SendEmailNew(e Email) error {

	from := mail.Address{"ООО Какая-то организация", e.From}
	to := mail.Address{"", e.To[0]}
	subj := e.Subject
	body := e.Body
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj
	headers["Content-type"] = "text/html;charset=\"utf-8\""
	headers["MIME-Version"] = "1.0"
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	servername := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", smtpUser, smtpPass, host)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}
func SendPurchaseInfo(UserEmail string, po []Puchase44) {
	var body = Body{po}
	b, _ := template.New("emailT").Parse(BodyEmailNew)
	var tpl bytes.Buffer
	_ = b.Execute(&tpl, body)
	BodyString := tpl.String()
	to := []string{UserEmail}
	var e = Email{`Тенедеры с неистекшим сроком подачи по запросам`, smtpUser, to, BodyString}
	err := SendEmailNew(e)
	if err != nil {
		Logging(err)

	}
}

var BodyEmailNew = `<p><strong>Добрый день!</strong><br />По Вашему запросу предоставлен список тендеров, срок подачи заявок на которые еще не истек. </p>
{{range .Purchases}}
<div><p><strong>Закупка:</strong> {{.PurNum}}<br /><strong>Наименование:</strong> {{.PurName}}<br /><strong>Дата публикации:</strong> {{.PubDate}}<br /><strong>НМЦК:</strong> {{.Nmck}}<br /><strong>Ссылка:</strong> {{.Href}}<hr /><p><p><p></div>
{{end}}
`
