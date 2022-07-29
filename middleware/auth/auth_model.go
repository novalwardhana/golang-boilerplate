package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Role struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type JwtUserData struct {
	User
	Roles []Role `json:"roles"`
}

type JwtCustomClaims struct {
	Data JwtUserData `json:"data"`
	jwt.StandardClaims
}

type NewContext struct {
	User  User
	Roles []Role
	echo.Context
}
