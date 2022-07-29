package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Role struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Result struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var JWTSignatureKey = []byte("novalwardhana")

type JWTData struct {
	Data JWTUserData `json:"data"`
	jwt.StandardClaims
}

type JWTUserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Roles []Role `json:"roles"`
}
