package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/novalwardhana/golang-boilerplate/middleware/auth"
	"github.com/novalwardhana/golang-boilerplate/module/user-authentication/model"
	"github.com/novalwardhana/golang-boilerplate/module/user-authentication/usecase"
)

type Handler struct {
	uc usecase.Usecase
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) Mount(group *echo.Group) {
	group.POST("/login", h.login)
	group.POST("/login-test", h.loginTest, auth.CheckAuth())
}

// Login:
func (h *Handler) login(c echo.Context) error {

	/* payload validation */
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}
	if len(user.Email) == 0 || len(user.Password) == 0 {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Email and password must be filled"})
	}

	/* Login process */
	result := <-h.uc.UserLogin(user.Email, user.Password)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success login and generate JSON Web Token", Data: result.Data})
}

// LoginTest:
func (h *Handler) loginTest(c echo.Context) error {
	mc := c.(auth.NewContext)

	/* Payload validation */
	type TestPayload struct {
		Text   string `json:"text"`
		Number int    `json:"number"`
	}
	var testPayload = new(TestPayload)
	if err := mc.Bind(testPayload); err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success login test", Data: map[string]interface{}{
		"user":    mc.User,
		"roles":   mc.Roles,
		"payload": testPayload,
	}})
}
