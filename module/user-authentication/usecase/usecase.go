package usecase

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/novalwardhana/golang-boilerplate/module/user-authentication/model"
	"github.com/novalwardhana/golang-boilerplate/module/user-authentication/repository"
)

type usecase struct {
	repo repository.Repository
}

type Usecase interface {
	UserLogin(username, password string) <-chan model.Result
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

// UserLogin:
func (uc *usecase) UserLogin(email, password string) <-chan model.Result {
	result := make(chan model.Result)
	go func() {
		defer close(result)

		/* Process get user */
		processGetUser := <-uc.repo.GetUser(email)
		if processGetUser.Error != nil {
			result <- model.Result{Error: processGetUser.Error}
			return
		}
		user := processGetUser.Data.(model.User)

		/* Password check */
		passwordByte := []byte(password)
		passwordHex := md5.Sum(passwordByte)
		passwordHash := hex.EncodeToString(passwordHex[:])
		if passwordHash != user.Password {
			result <- model.Result{Error: errors.New("Password not match")}
			return
		}

		/* Process get roles */
		processGetRoles := <-uc.repo.GetRole(user.ID)
		if processGetRoles.Error != nil {
			result <- model.Result{Error: errors.New("User role not found")}
			return
		}
		roles := processGetRoles.Data.([]model.Role)
		if len(roles) == 0 {
			result <- model.Result{Error: errors.New("User role not found")}
			return
		}

		/* Generate JWT */
		jwtData := model.JWTData{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "Golang Boilerplate",
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
			},
			Data: model.JWTUserData{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Roles: roles,
			},
		}
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
		jwtString, err := jwtToken.SignedString(model.JWTSignatureKey)
		if err != nil {
			result <- model.Result{Error: err}
			return
		}

		result <- model.Result{Data: jwtString}
	}()
	return result
}
