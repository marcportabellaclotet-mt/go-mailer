# gocli-mailer
Go CLI Email Sender

# usage
```
Usage of mailer:

  -mailAttachment string
    	file to be attached in the email (optional)
  -mailBody string
    	htlm message body (optional)
  -mailFrom string
    	sender's address email (required)
  -mailSubject string
    	message subject (optional)
  -mailTo string
    	recipient(s) email address (comma separated) (required)
  -smtpConTimeout int
    	SMTP Connection timeout seconds (optional) (default 10)
  -smtpEncryption string
    	SMTP Encryption or none (optional)  (default "none")
  -smtpHost string
    	smtp server host address (required)
  -smtpPassword string
    	smtp auth password (optional)
  -smtpPort int
    	smtp server connection port (optional) (default 587)
  -smtpSendTimeout int
    	SMTP Send timeout seconds (optional) (default 10)
  -smtpUsername string
    	smtp auth username (optional)
```