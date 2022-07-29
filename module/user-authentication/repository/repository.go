package repository

import (
	"github.com/novalwardhana/golang-boilerplate/module/user-authentication/model"
	"gorm.io/gorm"
)

type repository struct {
	dbMaster *gorm.DB
}

type Repository interface {
	GetUser(email string) <-chan model.Result
	GetRole(id int) <-chan model.Result
}

func NewRepository(dbMaster *gorm.DB) Repository {
	return &repository{
		dbMaster: dbMaster,
	}
}

// GetUser:
func (repo *repository) GetUser(email string) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get user */
		var user model.User
		sql := `select * from users where email = ?`
		if err := repo.dbMaster.Raw(sql, email).First(&user).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}
		result <- model.Result{Data: user}

	}()
	return result
}

// GetRole:
func (repo *repository) GetRole(id int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get roles */
		var roles []model.Role
		sql := `select
				roles.id,
				roles.code,
				roles.name
			from user_has_roles 
			inner join roles on roles.id = user_has_roles.role_id
			where user_has_roles.user_id = ?
			order by roles.id asc
		`
		if err := repo.dbMaster.Raw(sql, id).Find(&roles).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}
		result <- model.Result{Data: roles}

	}()
	return result
}
