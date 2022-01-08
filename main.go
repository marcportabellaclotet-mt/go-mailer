package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type smtpConnectionResource struct {
	hostname    string
	port        int64
	username    string
	password    string
	encryption  string
	conTimeout  int64
	sendTimeout int64
}

type emailMessage struct {
	from       string
	to         string
	body       string
	subject    string
	attachment string
}

func mailEncryption(enc string) mail.Encryption {
	switch enc {
	case "TLS":
		return mail.EncryptionTLS
	case "STARTTLS":
		return mail.EncryptionSTARTTLS
	case "SSL":
		return mail.EncryptionSSL
	case "SSLTLS":
		return mail.EncryptionSSLTLS
	}
	return mail.EncryptionNone
}

func smtpConnect(con smtpConnectionResource) (*mail.SMTPClient, error) {
	client := mail.NewSMTPClient()
	client.Host = con.hostname
	client.Port = int(con.port)
	if con.username != "" {
		client.Username = con.username
		client.Password = con.password
	}
	client.Encryption = mailEncryption(con.encryption)
	client.ConnectTimeout = 10 * time.Second
	client.SendTimeout = 10 * time.Second
	smtpClient, err := client.Connect()
	return smtpClient, err
}

func sendEmail(mailMessage emailMessage, smtpConnection smtpConnectionResource) {
	zp := regexp.MustCompile(` *, *`)
	for _, recipient := range zp.Split(mailMessage.to, -1) {
		smtpClient, err := smtpConnect(smtpConnection)
		if err != nil {
			log.Fatal(err)
		}

		email := mail.NewMSG()
		email.SetFrom(mailMessage.from).
			AddTo(recipient).
			SetSubject(mailMessage.subject).
			SetBody(mail.TextHTML, mailMessage.body)

		if mailMessage.attachment != "" {
			email.AddAttachment(mailMessage.attachment)
		}
		if email.Error != nil {
			log.Fatal(email.Error)
		}

		// Call Send and pass the client
		err = email.Send(smtpClient)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Email Sent to %s", recipient)
		}
	}
}

func checkFlags(con smtpConnectionResource, msg emailMessage) {
	error := false
	if con.hostname == "" {
		log.Println("smtpHost parameter missing")
		error = true
	}
	if msg.to == "" {
		log.Println("mailTo parameter missing")
		error = true
	}
	if msg.from == "" {
		log.Println("mailFrom parameter missing")
		error = true
	}
	if error {
		fmt.Println("\nusage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	mailFrom := flag.String("mailFrom", "", "sender's address email (required)")
	mailTo := flag.String("mailTo", "", "recipient(s) email address (comma separated) (required)")
	mailBody := flag.String("mailBody", "", "htlm message body (optional) ")
	mailSubject := flag.String("mailSubject", "", "message subject (optional)")
	mailAttachment := flag.String("mailAttachment", "", "file to be attached in the email (optional)")
	smtpHost := flag.String("smtpHost", "", "smtp server host address (required)")
	smtpPort := flag.Int64("smtpPort", 587, "smtp server connection port (optional)")
	smtpUsername := flag.String("smtpUsername", "", "smtp auth username (optional)")
	smtpPassword := flag.String("smtpPassword", "", "smtp auth password (optional)")
	smtpEncryption := flag.String("smtpEncryption", "none", "SMTP Encryption or none (optional) ")
	smtpConTimeout := flag.Int64("smtpConTimeout", 10, "SMTP Connection timeout seconds (optional)")
	smtpSendTimeout := flag.Int64("smtpSendTimeout", 10, "SMTP Send timeout seconds (optional)")
	flag.Parse()

	smtpConnection := smtpConnectionResource{
		hostname:    *smtpHost,
		port:        *smtpPort,
		username:    *smtpUsername,
		password:    *smtpPassword,
		encryption:  *smtpEncryption,
		conTimeout:  *smtpConTimeout,
		sendTimeout: *smtpSendTimeout,
	}

	mailMessage := emailMessage{
		from:       *mailFrom,
		to:         *mailTo,
		body:       *mailBody,
		subject:    *mailSubject,
		attachment: *mailAttachment,
	}

	checkFlags(smtpConnection, mailMessage)
	sendEmail(mailMessage, smtpConnection)
}
