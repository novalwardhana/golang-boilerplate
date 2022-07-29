package usecase

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math"

	"github.com/novalwardhana/golang-boilerplate/module/user-management/model"
	"github.com/novalwardhana/golang-boilerplate/module/user-management/repository"
)

type usecase struct {
	repo repository.Repository
}

type Usecase interface {
	Create(user *model.NewUser) <-chan model.Result
	GetData(page, limit int) <-chan model.Result
	Detail(id int) <-chan model.Result
	Update(user *model.NewUser) <-chan model.Result
	Delete(id int) <-chan model.Result
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

// Create:
func (u *usecase) Create(payload *model.NewUser) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Create new user data process */
		user := payload.User
		passwordByte := []byte(user.Password)
		passwordHex := md5.Sum(passwordByte)
		passwordHash := hex.EncodeToString(passwordHex[:])
		user.Password = passwordHash
		roleIDs := payload.Roles
		processCreateUser := <-u.repo.Create(&user, roleIDs)
		if processCreateUser.Error != nil {
			result <- model.Result{Error: processCreateUser.Error}
			return
		}
		user = *(processCreateUser.Data.(*model.User))
		user.Password = ""

		/* Get roles process */
		processGetRoles := <-u.repo.GetRoles(user.ID)
		if processGetRoles.Error != nil {
			result <- model.Result{Error: processGetRoles.Error}
			return
		}
		roles := processGetRoles.Data.([]model.Role)

		result <- model.Result{Data: model.UserWithRoles{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			JsonRoles: roles,
		}}

	}()
	return result
}

// GetData:
func (u *usecase) GetData(page, limit int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Count data process */
		processCountData := <-u.repo.CountData()
		if processCountData.Error != nil {
			result <- model.Result{Error: processCountData.Error}
			return
		}
		totalData := int(processCountData.Data.(int64))
		numberOfPage := int(math.Ceil(float64(totalData) / float64(limit)))

		/* Get data process */
		processGetData := <-u.repo.GetData(page, limit)
		if processGetData.Error != nil {
			result <- model.Result{Error: processGetData.Error}
			return
		}
		result <- model.Result{Data: model.Pagination{
			Page:         page,
			Limit:        limit,
			TotalData:    totalData,
			NumberOfPage: numberOfPage,
			Data:         processGetData.Data.([]model.UserWithRoles),
		}}

	}()
	return result
}

// Detail:
func (u *usecase) Detail(id int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Get user process */
		processGetUser := <-u.repo.GetUser(id)
		if processGetUser.Error != nil {
			result <- model.Result{Error: processGetUser.Error}
			return
		}
		user := processGetUser.Data.(model.User)
		user.Password = ""

		/* Get roles process */
		processGetRoles := <-u.repo.GetRoles(user.ID)
		if processGetRoles.Error != nil {
			result <- model.Result{Error: processGetRoles.Error}
			return
		}
		roles := processGetRoles.Data.([]model.Role)

		result <- model.Result{Data: model.UserWithRoles{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			JsonRoles: roles,
		}}
	}()
	return result
}

// Update:
func (u *usecase) Update(user *model.NewUser) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Hash password */
		passwordHex := md5.Sum([]byte(user.Password))
		passwordHash := hex.EncodeToString(passwordHex[:])
		user.Password = passwordHash

		/* Update data process */
		processUpdateData := <-u.repo.Update(user)
		if processUpdateData.Error != nil {
			result <- model.Result{Error: processUpdateData.Error}
			return
		}
		user := processUpdateData.Data.(model.User)

		/* Get roles process */
		processGetRoles := <-u.repo.GetRoles(user.ID)
		if processGetRoles.Error != nil {
			result <- model.Result{Error: processGetRoles.Error}
			return
		}
		roles := processGetRoles.Data.([]model.Role)

		result <- model.Result{Data: model.UserWithRoles{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			JsonRoles: roles,
		}}

	}()
	return result
}

// Delete:
func (u *usecase) Delete(id int) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Get roles process */
		processGetRoles := <-u.repo.GetRoles(id)
		if processGetRoles.Error != nil {
			result <- model.Result{Error: processGetRoles.Error}
			return
		}
		roles := processGetRoles.Data.([]model.Role)
		var isRootAdmin bool
		for _, role := range roles {
			if role.Code == "root" || role.Code == "admin" {
				isRootAdmin = true
				break
			}
		}

		/* Check role */
		if isRootAdmin {
			result <- model.Result{Error: errors.New("Cannot delete user with role root or admin")}
			return
		}

		/* Delete user process */
		processDeleteUser := <-u.repo.Delete(id)
		if processDeleteUser.Error != nil {
			result <- model.Result{Error: processDeleteUser.Error}
			return
		}

		result <- model.Result{}

	}()
	return result
}
