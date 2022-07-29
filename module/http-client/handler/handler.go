package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo"
	"github.com/novalwardhana/golang-boilerplate/config/env"
	"github.com/novalwardhana/golang-boilerplate/module/http-client/model"
	"github.com/novalwardhana/golang-boilerplate/module/http-client/usecase"
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
	group.POST("/crud/create", h.crudCreate)
	group.GET("/crud/get-data", h.crudGetData)
	group.POST("/advance-crud/bulk-insert", h.advanceCrudBulkInsert)
	group.GET("/advance-crud/download-csv", h.advanceDownloadCsv)
}

// CrudCreate:
func (h *Handler) crudCreate(c echo.Context) error {

	/* Payload verify */
	var payload = new(model.Person)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Process */
	result := <-h.usecase.Create(payload)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success insert new data"})
}

// CrudGetData:
func (h *Handler) crudGetData(c echo.Context) error {

	/* Page parameter validation */
	paramPage := c.QueryParam("page")
	page, err := strconv.Atoi(paramPage)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Page parameter not valid"})
	}
	if page <= 0 {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Page parameter not valid"})
	}

	/* Limit parameter not valid */
	paramLimit := c.QueryParam("limit")
	limit, err := strconv.Atoi(paramLimit)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Limit parameter not valid"})
	}
	if limit <= 0 {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Limit parameter not valid"})
	}

	/* Process */
	result := <-h.usecase.GetData(page, limit)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success get data", Data: result.Data})
}

// AdvanceCrudInsertBulk
func (h *Handler) advanceCrudBulkInsert(c echo.Context) error {

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

// AdvanceCrudExportCsv:
func (h *Handler) advanceDownloadCsv(c echo.Context) error {

	/* process */
	result := <-h.usecase.DownloadCSV()
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	filename := result.Data.(string)

	return c.Attachment(filepath.Join(os.Getenv(env.EnvHTTPClientDirectory), filename), filename)
}
