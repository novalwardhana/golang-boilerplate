package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/novalwardhana/golang-boilerplate/middleware/auth"
	"github.com/novalwardhana/golang-boilerplate/module/user-management/model"
	"github.com/novalwardhana/golang-boilerplate/module/user-management/usecase"
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
	group.POST("/create", h.Create, auth.CheckAuth())
	group.GET("/get-data", h.GetData, auth.CheckAuth())
	group.GET("/detail/:id", h.Detail, auth.CheckAuth())
	group.PUT("/update/:id", h.Update, auth.CheckAuth())
	group.DELETE("/delete/:id", h.Delete, auth.CheckAuth())
}

// Create:
func (h *Handler) Create(c echo.Context) error {

	mc := c.(auth.NewContext)

	/* Role check */
	var grant bool
	for _, role := range mc.Roles {
		if role.Code == "root" || role.Code == "admin" {
			grant = true
			break
		}
	}
	if !grant {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusUnauthorized, Message: "User not have grant to make new user"})
	}

	/* Payload validation */
	payload := new(model.NewUser)
	if err := mc.Bind(payload); err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Process add new user */
	result := <-h.usecase.Create(payload)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success add new user", Data: result.Data})
}

// GetData:
func (h *Handler) GetData(c echo.Context) error {

	mc := c.(auth.NewContext)

	/* Role check */
	var grant bool
	for _, role := range mc.Roles {
		if role.Code == "root" || role.Code == "admin" {
			grant = true
			break
		}
	}
	if !grant {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusUnauthorized, Message: "User not have grant to access user list"})
	}

	/* Page parameter validation */
	paramPage := mc.QueryParam("page")
	page, err := strconv.Atoi(paramPage)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Limit parameter validation */
	paramLimit := mc.QueryParam("limit")
	limit, err := strconv.Atoi(paramLimit)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Process get data */
	result := <-h.usecase.GetData(page, limit)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success get User Data", Data: result.Data})
}

// Detail:
func (h *Handler) Detail(c echo.Context) error {

	mc := c.(auth.NewContext)

	/* ID parameter validation */
	idParam := mc.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Role check */
	var grant bool
	for _, role := range mc.Roles {
		if role.Code == "root" || role.Code == "admin" {
			grant = true
			break
		}
	}
	if mc.User.ID == id {
		grant = true
	}
	if !grant {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusUnauthorized, Message: "User not have grant to access user list"})
	}

	/* Detail process */
	result := <-h.usecase.Detail(id)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success get user data", Data: result.Data})
}

// Update:
func (h *Handler) Update(c echo.Context) error {

	mc := c.(auth.NewContext)

	/* ID paramater validation */
	paramID := mc.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Role check */
	var grant bool
	for _, role := range mc.Roles {
		if role.Code == "root" || role.Code == "admin" {
			grant = true
		}
	}
	if mc.User.ID == id {
		grant = true
	}
	if !grant {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusUnauthorized, Message: "User not have access to update data"})
	}

	/* Payload */
	payload := new(model.NewUser)
	if err := mc.Bind(payload); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}
	payload.ID = id

	/* Process update data */
	result := <-h.usecase.Update(payload)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success update user data", Data: result.Data})
}

// Delete:
func (h *Handler) Delete(c echo.Context) error {

	mc := c.(auth.NewContext)

	/* ID parameter validation */
	idParam := mc.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusBadRequest, Message: err.Error()})
	}

	/* Role check */
	var grant bool
	for _, role := range mc.Roles {
		if role.Code == "root" || role.Code == "admin" {
			grant = true
			break
		}
	}
	if !grant {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusUnauthorized, Message: "User not have access to delete data"})
	}

	/* Process delete data */
	result := <-h.usecase.Delete(id)
	if result.Error != nil {
		return c.JSON(http.StatusOK, model.Response{Status: http.StatusNotFound, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{Status: http.StatusOK, Message: "Success delete user"})
}
