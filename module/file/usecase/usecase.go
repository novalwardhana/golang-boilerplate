package usecase

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/module/file/model"
	"github.com/novalwardhana/golang-boilerplate/module/file/repository"
)

type usecase struct {
	repo repository.Repository
}

type Usecase interface {
	Upload(file *multipart.FileHeader) <-chan model.Result
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

// Upload:
func (u *usecase) Upload(file *multipart.FileHeader) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Check filedir */
		filedir := os.Getenv(env.EnvFileDirectory)
		if err := os.MkdirAll(filedir, os.ModePerm); err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Create file target */
		filename := time.Now().Format("2006_01_02_15_04_05") + "_" + file.Filename
		fileTarget, err := os.Create(filepath.Join(filedir, filename))
		defer fileTarget.Close()
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Create file source */
		fileSource, err := file.Open()
		defer fileSource.Close()
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Copy file source to file target */
		_, err = io.Copy(fileTarget, fileSource)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		result <- model.Result{}
	}()
	return result
}
