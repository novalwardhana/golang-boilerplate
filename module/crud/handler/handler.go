package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/novalwardhana/golang-boilerplate/module/crud/model"
	"github.com/novalwardhana/golang-boilerplate/module/crud/usecase"
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
	group.POST("/create", h.create)
	group.GET("/get-data", h.getData)
	group.GET("/detail", h.detail)
	group.PUT("/update/:id", h.update)
	group.DELETE("/delete/:id", h.delete)
}

// Create:
func (h *Handler) create(c echo.Context) error {

	/* Payload validation */
	params := new(model.Person)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Create process */
	result := <-h.uc.Create(params)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success create new data", Data: params})
}

// GetData:
func (h *Handler) getData(c echo.Context) error {

	/* Page parameter validation */
	paramPage := c.QueryParam("page")
	page, err := strconv.Atoi(paramPage)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Page parameter not valid"})
	}
	if page <= 0 {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Page parameter not valid"})
	}

	/* Limit parameter validation */
	paramLimit := c.QueryParam("limit")
	limit, err := strconv.Atoi(paramLimit)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Limit parameter not valid"})
	}
	if limit <= 0 {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: "Limit parameter not valid"})
	}

	/* Get data process */
	result := <-h.uc.GetData(page, limit)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success get data", Data: result.Data})
}

// Detail:
func (h *Handler) detail(c echo.Context) error {

	/* ID paramater validation */
	idParam := c.QueryParam("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Detail process */
	result := <-h.uc.Detail(id)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success get data", Data: result.Data})
}

// Update:
func (h *Handler) update(c echo.Context) error {

	/* ID parameter validation */
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Payload validation */
	params := new(model.Person)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}
	params.ID = id

	/* Update process */
	result := <-h.uc.Update(params)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success update data", Data: result.Data})
}

// Delete
func (h *Handler) delete(c echo.Context) error {

	/* ID paramater validation */
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: err.Error()})
	}

	/* Delete process */
	result := <-h.uc.Delete(id)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success delete data"})
}
