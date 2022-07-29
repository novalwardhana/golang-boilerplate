package repository

import (
	"encoding/json"

	"github.com/novalwardhana/golang-boilerplate/module/user-management/model"
	"gorm.io/gorm"
)

type repository struct {
	dbMaster *gorm.DB
}

type Repository interface {
	Create(user *model.User, roles []int) <-chan model.Result
	CountData() <-chan model.Result
	GetData(page, limit int) <-chan model.Result
	GetRoles(userID int) <-chan model.Result
	GetUser(id int) <-chan model.Result
	Update(payload *model.NewUser) <-chan model.Result
	Delete(id int) <-chan model.Result
}

func NewRepository(dbMaster *gorm.DB) Repository {
	return &repository{
		dbMaster: dbMaster,
	}
}

// Create:
func (r *repository) Create(user *model.User, roles []int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process create new user */
		tx := r.dbMaster.Begin()
		if err := tx.Create(user).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Process create user has roles */
		var userHasRoles []model.UserHasRole
		for _, roleID := range roles {
			userHasRole := model.UserHasRole{
				UserID: user.ID,
				RoleID: roleID,
			}
			userHasRoles = append(userHasRoles, userHasRole)
		}
		if err := tx.Create(&userHasRoles).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}

		tx.Commit()
		result <- model.Result{Data: user}

	}()
	return result
}

// CountData:
func (r *repository) CountData() <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process count data */
		var count int64
		sql := `select count(id) from users`
		if err := r.dbMaster.Raw(sql).Count(&count).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}
		result <- model.Result{Data: count}

	}()
	return result
}

// GetData:
func (r *repository) GetData(page, limit int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get data */
		var list []model.UserWithRoles
		offset := (page - 1) * limit
		sql := `select 
					u.id,
					u.name,
					u.email,
					jsonb_agg(concat('{', 
						'"id"', ':', r.id , ',',
						'"code"', ':', '"', r.code , '",',
						'"name"', ':', '"', r.name , '"',
					'}')::json) as roles
				from users as u
				inner join user_has_roles uhr on u.id = uhr.user_id 
				inner join roles r on uhr.role_id = r.id
				group by u.id, u.name, u.email
				order by u.id desc 
				offset ? limit ?`
		if err := r.dbMaster.Raw(sql, offset, limit).Find(&list).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}
		for index := range list {
			var roles []model.Role
			if err := json.Unmarshal(list[index].Roles, &roles); err != nil {
				result <- model.Result{Error: err}
				return
			}
			list[index].JsonRoles = roles
		}
		result <- model.Result{Data: list}

	}()
	return result
}

// GetRoles:
func (r *repository) GetRoles(userID int) <-chan model.Result {
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
		if err := r.dbMaster.Raw(sql, userID).Find(&roles).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}
		result <- model.Result{Data: roles}

	}()
	return result
}

// GetUser:
func (r *repository) GetUser(id int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get user */
		var user model.User
		sql := "select * from users where id = ?"
		if err := r.dbMaster.Raw(sql, id).Find(&user).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}
		result <- model.Result{Data: user}

	}()
	return result
}

// Update:
func (r *repository) Update(payload *model.NewUser) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get data */
		tx := r.dbMaster.Begin()
		var user model.User
		sql := `select * from users where id = ?`
		if err := tx.Raw(sql, payload.ID).First(&user).Error; err != nil {
			tx.Rollback()
			result <- model.Result{Error: err}
			return
		}

		/* Process update data */
		user.Name = payload.Name
		user.Password = payload.Password
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			result <- model.Result{Error: err}
			return
		}

		/* Delete current user has roles */
		sql = `delete from user_has_roles where user_id = ?`
		if err := tx.Exec(sql, user.ID).Error; err != nil {
			tx.Rollback()
			result <- model.Result{Error: err}
			return
		}

		/* Insert new user has roles */
		var userHasRoles []model.UserHasRole
		for _, roleID := range payload.Roles {
			userHasRole := model.UserHasRole{
				UserID: user.ID,
				RoleID: roleID,
			}
			userHasRoles = append(userHasRoles, userHasRole)
		}
		if err := tx.Create(&userHasRoles).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}

		tx.Commit()
		result <- model.Result{Data: user}

	}()
	return result
}

// Delete:
func (r *repository) Delete(id int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Delete user has roles */
		tx := r.dbMaster.Begin()
		sql := `delete from user_has_roles where user_id = ?`
		if err := tx.Exec(sql, id).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}

		/* Delete user */
		if err := tx.Delete(&model.User{}, id).Error; err != nil {
			result <- model.Result{Error: err}
			return
		}

		tx.Commit()
		result <- model.Result{}
	}()
	return result
}
