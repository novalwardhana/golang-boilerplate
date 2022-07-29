package repository

import (
	"sync"

	"github.com/novalwardhana/golang-boilerplate/module/advance-crud/model"
	"gorm.io/gorm"
)

type repository struct {
	dbMaster *gorm.DB
}

type Repository interface {
	Insert(payload *[]*model.Person) <-chan model.Result
	GetData(persons *[]*model.Person) <-chan model.Result
}

func NewRepository(dbMaster *gorm.DB) Repository {
	return &repository{
		dbMaster: dbMaster,
	}
}

// Insert:
func (r *repository) Insert(payload *[]*model.Person) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process insert data to database */
		wg := sync.WaitGroup{}
		for _, data := range *payload {
			wg.Add(1)
			person := data
			go func() {
				defer wg.Done()
				tx := r.dbMaster.Begin()
				if err := tx.Create(person).Error; err != nil {
					tx.Rollback()
				}
				tx.Commit()
			}()
		}
		wg.Wait()

		result <- model.Result{}
	}()
	return result
}

// ExportCSV:
func (r *repository) GetData(persons *[]*model.Person) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		sql := `select * from persons order by id`
		rows, err := r.dbMaster.Raw(sql).Rows()
		if err != nil {
			result <- model.Result{Error: err}
			return
		}
		for rows.Next() {
			var person model.Person
			if err := rows.Scan(
				&person.ID,
				&person.Name,
				&person.Age,
				&person.Address,
			); err != nil {
				continue
			}
			*persons = append(*persons, &person)
		}

	}()
	return result
}
