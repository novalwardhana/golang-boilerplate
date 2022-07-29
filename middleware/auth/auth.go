package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// CheckAuth:
func CheckAuth() echo.MiddlewareFunc {

	/* Return handler function */
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		/* Return http */
		return func(c echo.Context) error {

			/* Get Auth header */
			headers := c.Request().Header
			authorization := headers.Get("authorization")
			if len(authorization) == 0 {
				return c.JSON(http.StatusOK, Response{Status: http.StatusUnauthorized, Message: "Auth bearer must provided"})
			}

			/* Bearer validation */
			authorizationaArray := strings.Split(authorization, " ")
			if len(authorizationaArray) < 2 {
				return c.JSON(http.StatusOK, Response{Status: http.StatusUnauthorized, Message: "Bearer not valid"})
			}
			if authorizationaArray[0] != "Bearer" {
				return c.JSON(http.StatusOK, Response{Status: http.StatusUnauthorized, Message: "Bearer not valid"})
			}
			if len(authorizationaArray[1]) < 1 {
				return c.JSON(http.StatusOK, Response{Status: http.StatusUnauthorized, Message: "Bearer not valid"})
			}

			/* Token validation */
			token, err := jwt.ParseWithClaims(authorizationaArray[1], &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte("novalwardhana"), nil
			})
			if err != nil {
				return c.JSON(http.StatusOK, Response{Status: http.StatusUnauthorized, Message: err.Error()})
			}
			if !token.Valid {
				return c.JSON(http.StatusOK, Response{Status: http.StatusUnauthorized, Message: "Token invalid"})
			}

			/* Decode token */
			decodeToken, status := token.Claims.(*JwtCustomClaims)
			if !status {
				return c.JSON(http.StatusOK, Response{Status: http.StatusUnauthorized, Message: "Failed decode token"})
			}

			/* Expired token validation */
			if decodeToken.ExpiresAt < time.Now().Local().Unix() {
				return c.JSON(http.StatusOK, Response{Status: http.StatusUnauthorized, Message: "Token expired"})
			}

			return next(NewContext{User: decodeToken.Data.User, Roles: decodeToken.Data.Roles, Context: c})
		}

	}
}
