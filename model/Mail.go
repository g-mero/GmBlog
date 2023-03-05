package model

import (
	"gmeroblog/utils/errmsg"
	"log"
	"strconv"

	mail "github.com/xhit/go-simple-mail/v2"
)

type smtpInfo struct {
	host     string
	port     int
	username string
	password string
}

func getSmtpInfo() (smtp smtpInfo) {
	smtp.host = SITE_SETTING["mail_host"]
	smtp.username = SITE_SETTING["mail_username"]
	smtp.password = SITE_SETTING["mail_password"]
	smtp.port, _ = strconv.Atoi(SITE_SETTING["mail_port"])

	return smtp
}

func Enabled() bool {
	return SITE_SETTING["mail_host"] != ""
}

func newSMTPClient() (*mail.SMTPClient, error) {
	// 使用网易的邮箱
	server := mail.NewSMTPClient()
	server.Host = getSmtpInfo().host
	server.Port = getSmtpInfo().port
	server.Username = getSmtpInfo().username
	server.Password = getSmtpInfo().password
	server.Encryption = mail.EncryptionTLS

	// 外部调用记得调用Close()
	return server.Connect()
}

func SendEmail(subject, content, to string) int {
	if !Enabled() {
		return errmsg.ERROR_MAIL_NOT_ENABLIED
	}

	smtpClient, err := newSMTPClient()
	if err != nil {
		log.Println("[MAIL]", err)
		return errmsg.ERROR_MAIL_SERVER
	}
	defer smtpClient.Close()

	// Create email
	email := mail.NewMSG()
	email.SetFrom("From gmero.com <bot_gmero@163.com>")
	email.AddTo("2431775986@qq.com")
	// email.AddCc("another_you@example.com")
	email.SetSubject(subject)

	email.SetBody(mail.TextHTML, content) //发送html信息
	// email.AddAttachment("super_cool_file.png") // 附件

	// Send email
	if err = email.Send(smtpClient); err != nil {
		log.Println("[MAIL]", err)
		return errmsg.ERROR_MAIL_SEND
	}

	return errmsg.SUCCES
}
