package usecase

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/module/advance-crud/model"
	"github.com/novalwardhana/golang-boilerplate/module/advance-crud/repository"
)

type usecase struct {
	repo repository.Repository
}

type Usecase interface {
	BulkInsert(file *multipart.FileHeader) <-chan model.Result
	ExportCSV() <-chan model.Result
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

// BulkInsert:
func (u *usecase) BulkInsert(file *multipart.FileHeader) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Check filedir */
		filedir := os.Getenv(env.EnvAdvanceCrudDirectory)
		if err := os.MkdirAll(filedir, os.ModePerm); err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Create file target */
		filename := time.Now().Format("20060102_150405") + "_" + file.Filename
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
		}

		/* Copy file source to file target */
		_, err = io.Copy(fileTarget, fileSource)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Read file after save to directory */
		var payload []*model.Person
		fileCSV, err := os.Open(filepath.Join(filedir, filename))
		defer fileCSV.Close()
		if err != nil {
			result <- model.Result{Error: err}
			return
		}
		fileScan := bufio.NewScanner(fileCSV)
		for fileScan.Scan() {
			rawData := fileScan.Text()
			arrData := strings.Split(rawData, ",")
			if len(arrData) != 3 {
				continue
			}
			age, err := strconv.Atoi(arrData[1])
			if err != nil {
				continue
			}
			person := &model.Person{
				Name:    arrData[0],
				Age:     age,
				Address: arrData[2],
			}
			payload = append(payload, person)
		}

		/* Insert into database */
		processInsert := <-u.repo.Insert(&payload)
		if processInsert.Error != nil {
			result <- model.Result{Error: processInsert.Error}
			return
		}

		result <- model.Result{}
	}()
	return result
}

// ExportCSV
func (u *usecase) ExportCSV() <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get data */
		var persons []*model.Person
		processGetData := <-u.repo.GetData(&persons)
		if processGetData.Error != nil {
			result <- model.Result{Error: processGetData.Error}
			return
		}

		/* Prepare csv data */
		var csvData string
		csvData += "ID,NAME,AGE,ADDRESS\n"
		for _, data := range persons {
			csvData += fmt.Sprintf("%d,%s,%d,%s\n", data.ID, data.Name, data.Age, data.Address)
		}

		/* Create csv */
		filedir := os.Getenv(env.EnvAdvanceCrudDirectory)
		filename := time.Now().Format("20060102_150405") + "_" + "Download_Data.csv"
		if err := ioutil.WriteFile(filepath.Join(filedir, filename), []byte(csvData), os.ModePerm); err != nil {
			result <- model.Result{Error: err}
			return
		}

		result <- model.Result{Data: filename}
	}()
	return result
}
