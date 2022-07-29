package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/novalwardhana/golang-boilerplate/middleware/auth"
	"github.com/novalwardhana/golang-boilerplate/module/email/model"
	"github.com/novalwardhana/golang-boilerplate/module/email/usecase"
)

type Handler struct {
	usecase usecase.Usecase
}

func NewHandler(usecase usecase.Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Mount(group *echo.Group) {
	group.POST("/send-mail-default", h.SendMailDefault, auth.CheckAuth())
	group.POST("/send-mail-gomail", h.SendMailGomail, auth.CheckAuth())
}

// SendMailDefault:
func (h *Handler) SendMailDefault(c echo.Context) error {

	/* Get payload */
	mc := c.(auth.NewContext)
	request := new(model.Request)
	if err := mc.Bind(request); err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Result */
	result := <-h.usecase.SendMailDefault(request.Email, request.Subject, request.Text)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success send email"})
}

// SendMailGomail:
func (h *Handler) SendMailGomail(c echo.Context) error {

	/* Get payload */
	mc := c.(auth.NewContext)
	email := mc.FormValue("email")
	subject := mc.FormValue("subject")
	text := mc.FormValue("text")
	if len(email) == 0 || len(subject) == 0 || len(text) == 0 {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Email, subject, and text must be filled"})
	}

	/* Get file */
	file, err := mc.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Result */
	result := <-h.usecase.SendMailGomail(email, subject, text, file)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success send email"})
}
