package repository

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strconv"

	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/module/email/model"
	"gopkg.in/gomail.v2"
)

type repository struct {
}

type Repository interface {
	SendMailDefault(email, subject, text string) <-chan model.Result
	SendMailGomail(email, subject, text, filedir, filename string) <-chan model.Result
}

func NewRepository() Repository {
	return &repository{}
}

type loginAuth struct {
	username string
	password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (l *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(l.username), nil
}

func (l *loginAuth) Next(fromServer []byte, more bool) (toServer []byte, err error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(l.username), nil
		case "Password:":
			return []byte(l.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

// SendMailDefault:
func (r *repository) SendMailDefault(email, subject, text string) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Net dial */
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv(env.EnvEmailHost), os.Getenv(env.EnvEmailPort)))
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Create new client */
		client, err := smtp.NewClient(conn, os.Getenv(env.EnvEmailHost))
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Start tls config */
		tlsConfig := &tls.Config{
			ServerName: os.Getenv(env.EnvEmailHost),
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Create auth */
		auth := LoginAuth(os.Getenv(env.EnvEmailUser), os.Getenv(env.EnvEmailPassword))
		if err := client.Auth(auth); err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Generate email body */
		var messageIndex []string
		var messageValue []string
		var messageContent string
		messageIndex = append(messageIndex, "From")
		messageValue = append(messageValue, os.Getenv(env.EnvEmailUser))
		messageIndex = append(messageIndex, "To")
		messageValue = append(messageValue, email)
		messageIndex = append(messageIndex, "Cc")
		messageValue = append(messageValue, "")
		messageIndex = append(messageIndex, "Subject")
		messageValue = append(messageValue, subject)
		messageIndex = append(messageIndex, "MIME-Version")
		messageValue = append(messageValue, "1.0")
		messageIndex = append(messageIndex, "Content-Type")
		messageValue = append(messageValue, "text/html; charset=\"UTF-8\";\n")
		messageIndex = append(messageIndex, "Text")
		messageValue = append(messageValue, fmt.Sprintf("<html><body><p>%s</p></body></html>", text))
		for i := 0; i < len(messageIndex); i++ {
			if messageIndex[i] != "Text" {
				messageContent += fmt.Sprintf("%s: %s\n", messageIndex[i], messageValue[i])
			} else {
				messageContent += messageValue[i]
			}
		}

		/* Send email */
		address := fmt.Sprintf("%s:%s", os.Getenv(env.EnvEmailHost), os.Getenv(env.EnvEmailPort))
		from := os.Getenv(env.EnvEmailUser)
		to := []string{
			email,
		}
		cc := []string{}
		if err := smtp.SendMail(address, auth, from, append(to, cc...), []byte(messageContent)); err != nil {
			result <- model.Result{Error: err}
			return
		}

		result <- model.Result{}
	}()
	return result
}

// SendMailGomail:
func (r *repository) SendMailGomail(email, subject, text, filedir, filename string) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Create dialer */
		port, err := strconv.Atoi(os.Getenv(env.EnvEmailPort))
		if err != nil {
			result <- model.Result{Error: err}
			return
		}
		dial := gomail.NewDialer(os.Getenv(env.EnvEmailHost), port, os.Getenv(env.EnvEmailUser), os.Getenv(env.EnvEmailPassword))

		/* Compose messages */
		message := gomail.NewMessage()
		message.SetHeader("From", os.Getenv(env.EnvEmailUser))
		message.SetHeader("To", email)
		message.SetHeader("Subject", subject)
		message.SetBody("text/html", fmt.Sprintf("<html><body><b>Halo bro,...</b></br>%s</body></html>", text))
		//message.Attach(filepath.Join(filedir, filename))

		/* Send email */
		if err := dial.DialAndSend(message); err != nil {
			result <- model.Result{Error: err}
			return
		}
		result <- model.Result{}
	}()
	return result
}
