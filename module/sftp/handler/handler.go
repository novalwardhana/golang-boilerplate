package handler

import (
	"github.com/labstack/echo"
	"github.com/novalwardhana/golang-boilerplate/module/sftp/usecase"
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
}
