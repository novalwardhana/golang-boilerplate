package usecase

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/module/http-client/model"
	"github.com/novalwardhana/golang-boilerplate/module/http-client/repository"
)

type usecase struct {
	repo repository.Repository
}

type Usecase interface {
	Create(payload *model.Person) <-chan model.Result
	GetData(page, limit int) <-chan model.Result
	BulkInsert(file *multipart.FileHeader) <-chan model.Result
	DownloadCSV() <-chan model.Result
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

// CrudCreate:
func (u *usecase) Create(payload *model.Person) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process create */
		processCrudCreate := <-u.repo.Create(payload)
		if processCrudCreate.Error != nil {
			result <- model.Result{Error: processCrudCreate.Error}
			return
		}

		result <- model.Result{}
	}()
	return result
}

// GetData
func (u *usecase) GetData(page, limit int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get data */
		processGetData := <-u.repo.GetData(page, limit)
		if processGetData.Error != nil {
			result <- model.Result{Error: processGetData.Error}
			return
		}

		result <- model.Result{Data: processGetData.Data}
	}()
	return result
}

// BulkInsert:
func (u *usecase) BulkInsert(file *multipart.FileHeader) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Prepare file directory */
		filedir := os.Getenv(env.EnvHTTPClientDirectory)
		if err := os.MkdirAll(filedir, os.ModePerm); err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Prepare file target */
		filename := time.Now().Format("20060102_150405") + "_" + file.Filename
		fileTarget, err := os.Create(filepath.Join(filedir, filename))
		defer fileTarget.Close()
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Prepare file source */
		fileSource, err := file.Open()
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Copy file target into file source */
		_, err = io.Copy(fileTarget, fileSource)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Process bulk insert */
		processBulkInsert := <-u.repo.BulkInsert(filedir, filename)
		if processBulkInsert.Error != nil {
			result <- model.Result{Error: processBulkInsert.Error}
			return
		}

		result <- model.Result{}
	}()
	return result
}

// ExportCSV:
func (u *usecase) DownloadCSV() <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process export csv */
		processDownloadCSV := <-u.repo.DownloadCSV()
		if processDownloadCSV.Error != nil {
			result <- model.Result{Error: processDownloadCSV.Error}
			return
		}
		csvContent := processDownloadCSV.Data.([]byte)

		/* Save csv into file */
		filedir := os.Getenv(env.EnvHTTPClientDirectory)
		filename := time.Now().Format("20060102_150405") + "_" + "download_csv.csv"
		if err := ioutil.WriteFile(filepath.Join(filedir, filename), csvContent, os.ModePerm); err != nil {
			result <- model.Result{Error: err}
			return
		}

		result <- model.Result{Data: filename}
	}()
	return result
}
