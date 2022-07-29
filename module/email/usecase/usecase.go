package usecase

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/module/email/model"
	"github.com/novalwardhana/golang-boilerplate/module/email/repository"
)

type usecase struct {
	repo repository.Repository
}

type Usecase interface {
	SendMailDefault(email, subject, text string) <-chan model.Result
	SendMailGomail(email, subject, text string, file *multipart.FileHeader) <-chan model.Result
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

// SendMailDefault:
func (u *usecase) SendMailDefault(email, subject, text string) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process */
		process := <-u.repo.SendMailDefault(email, subject, text)
		if process.Error != nil {
			result <- model.Result{Error: process.Error}
			return
		}
		result <- model.Result{}
	}()
	return result
}

// SendMailGomail:
func (u *usecase) SendMailGomail(email, subject, text string, file *multipart.FileHeader) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* File source */
		fileSource, err := file.Open()
		defer fileSource.Close()
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Prepare file directory */
		filedir := os.Getenv(env.EnvEmailAttachmentDirectory)
		if err := os.MkdirAll(filedir, os.ModePerm); err != nil {
			result <- model.Result{Error: err}
			return
		}
		filename := time.Now().Format("20060102_150405") + "_" + file.Filename

		/* File target */
		fileTarget, err := os.Create(filepath.Join(filedir, filename))
		defer fileTarget.Close()
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Save file attachment to directory */
		_, err = io.Copy(fileTarget, fileSource)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Process */
		process := <-u.repo.SendMailGomail(email, subject, text, filedir, filename)
		if process.Error != nil {
			result <- model.Result{Error: process.Error}
			return
		}
		result <- model.Result{}
	}()
	return result
}
