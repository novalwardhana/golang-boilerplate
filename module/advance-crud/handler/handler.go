package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/module/advance-crud/model"
	"github.com/novalwardhana/golang-boilerplate/module/advance-crud/usecase"
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
	group.POST("/bulk-insert", h.bulkInsert)
	group.GET("/export-csv", h.exportCSV)
}

// BulkInsert:
func (h *Handler) bulkInsert(c echo.Context) error {

	/* File validation */
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Process */
	result := <-h.usecase.BulkInsert(file)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success bulk insert"})
}

// ExportCSV:
func (h *Handler) exportCSV(c echo.Context) error {

	/* Process */
	result := <-h.usecase.ExportCSV()
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	filename := result.Data.(string)

	return c.Attachment(filepath.Join(os.Getenv(env.EnvAdvanceCrudDirectory), filename), filename)
}
