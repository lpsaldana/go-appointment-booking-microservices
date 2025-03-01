package config

import (
	"log"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPConfig() *SMTPConfig {
	// En producción, usa variables de entorno (ej: os.Getenv("SMTP_HOST"))
	return &SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: "tu.email@gmail.com",   // Cambia esto
		Password: "tu-contraseña-de-app", // Genera una contraseña de app en Gmail
		From:     "tu.email@gmail.com",   // Debe coincidir con Username
	}
}

func (c *SMTPConfig) SendMail(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)
	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(c.Host+":"+c.Port, auth, c.From, to, msg)
	if err != nil {
		log.Printf("Error al enviar correo: %v", err)
		return err
	}
	return nil
}
