package repository

import (
	"github.com/novalwardhana/golang-boilerplate/module/crud/model"
	"gorm.io/gorm"
)

type repository struct {
	dbMaster *gorm.DB
}

type Repository interface {
	Create(params *model.Person) <-chan model.Result
	CountData() <-chan model.Result
	GetData(page, limit int) <-chan model.Result
	Detail(id int) <-chan model.Result
	Update(params *model.Person) <-chan model.Result
	Delete(id int) <-chan model.Result
}

func NewRepository(dbMaster *gorm.DB) Repository {
	return &repository{
		dbMaster: dbMaster,
	}
}

// Create:
func (repo *repository) Create(params *model.Person) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process add data to database */
		tx := repo.dbMaster.Begin()
		if err := tx.Create(params).Error; err != nil {
			tx.Rollback()
			result <- model.Result{Error: err}
			return
		}
		tx.Commit()

		result <- model.Result{}
	}()
	return result
}

// CountData:
func (repo *repository) CountData() <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process count data */
		var count int64
		sql := `select count(id) from persons`
		if err := repo.dbMaster.Raw(sql).Count(&count).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}
		result <- model.Result{Data: count}

	}()
	return result
}

// GetData:
func (repo *repository) GetData(page, limit int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get data */
		var persons []model.Person
		offset := (page - 1) * limit
		sql := `select * from persons order by id desc offset ? limit ? `
		if err := repo.dbMaster.Raw(sql, offset, limit).Find(&persons).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}
		result <- model.Result{Data: persons}

	}()
	return result
}

// Detail:
func (repo *repository) Detail(id int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get data from database */
		var person model.Person
		sql := `select * from persons where id = ?`
		if err := repo.dbMaster.Raw(sql, id).First(&person).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}
		result <- model.Result{Data: person}

	}()
	return result
}

// Update:
func (repo *repository) Update(params *model.Person) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get data */
		tx := repo.dbMaster.Begin()
		var person model.Person
		sql := `select * from persons where id = ?`
		if err := tx.Raw(sql, params.ID).First(&person).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Process update data */
		person.Name = params.Name
		person.Age = params.Age
		person.Address = params.Address
		if err := tx.Save(&person).Error; err != nil {
			tx.Rollback()
			result <- model.Result{Error: err}
			return
		}
		tx.Commit()

		result <- model.Result{Data: person}
	}()
	return result
}

// Delete:
func (repo *repository) Delete(id int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process delete */
		tx := repo.dbMaster.Begin()
		if err := tx.Delete(&model.Person{}, id).Error; err != nil {
			tx.Rollback()
			result <- model.Result{Error: err}
			return
		}
		tx.Commit()
		result <- model.Result{}
	}()
	return result
}
