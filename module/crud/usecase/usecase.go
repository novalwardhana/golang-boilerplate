package usecase

import (
	"math"

	"github.com/novalwardhana/golang-boilerplate/module/crud/model"
	"github.com/novalwardhana/golang-boilerplate/module/crud/repository"
)

type usecase struct {
	repo repository.Repository
}

type Usecase interface {
	Create(params *model.Person) <-chan model.Result
	GetData(page, limit int) <-chan model.Result
	Detail(id int) <-chan model.Result
	Update(params *model.Person) <-chan model.Result
	Delete(id int) <-chan model.Result
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

// Create:
func (uc *usecase) Create(params *model.Person) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Create process */
		process := <-uc.repo.Create(params)
		if process.Error != nil {
			result <- model.Result{Error: process.Error}
			return
		}
		result <- model.Result{Data: process.Data}

	}()
	return result
}

// GetData:
func (uc *usecase) GetData(page, limit int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Count data process */
		processCountData := <-uc.repo.CountData()
		if processCountData.Error != nil {
			result <- model.Result{Error: processCountData.Error}
			return
		}
		totalData := int(processCountData.Data.(int64))
		numberOfPage := int(math.Ceil(float64(totalData) / float64(limit)))

		/* Get data process */
		processGetData := <-uc.repo.GetData(page, limit)
		if processGetData.Error != nil {
			result <- model.Result{Error: processGetData.Error}
			return
		}
		result <- model.Result{Data: model.Pagination{
			Page:         page,
			Limit:        limit,
			TotalData:    totalData,
			NumberOfPage: numberOfPage,
			Data:         processGetData.Data.([]model.Person),
		}}

	}()
	return result
}

// Detail:
func (uc *usecase) Detail(id int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Detail process */
		process := <-uc.repo.Detail(id)
		if process.Error != nil {
			result <- model.Result{Error: process.Error}
			return
		}
		result <- model.Result{Data: process.Data}

	}()
	return result
}

// Update:
func (uc *usecase) Update(params *model.Person) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Update process */
		process := <-uc.repo.Update(params)
		if process.Error != nil {
			result <- model.Result{Error: process.Error}
			return
		}
		result <- model.Result{Data: process.Data}

	}()
	return result
}

// Delete:
func (uc *usecase) Delete(id int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Delete process */
		process := <-uc.repo.Delete(id)
		if process.Error != nil {
			result <- model.Result{Error: process.Error}
			return
		}
		result <- model.Result{}

	}()
	return result
}
