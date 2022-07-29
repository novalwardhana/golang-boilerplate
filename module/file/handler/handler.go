package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/module/file/model"
	"github.com/novalwardhana/golang-boilerplate/module/file/usecase"
)

type Handler struct {
	usecase usecase.Usecase
}

func NewHandler(usecase usecase.Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Mount(g *echo.Group) {
	g.POST("/upload", h.upload)
	g.GET("/download", h.download)
}

func (h *Handler) upload(c echo.Context) error {

	/* File validation */
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Process */
	result := <-h.usecase.Upload(file)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success upload new file"})
}

// Download:
func (h *Handler) download(c echo.Context) error {

	/* filename parameter validation */
	filename := c.QueryParam("filename")
	if len(filename) == 0 {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Filename must filled"})
	}

	/* Download file */
	filedir := os.Getenv(env.EnvFileDirectory)
	return c.Attachment(filepath.Join(filedir, filename), filename)
}
