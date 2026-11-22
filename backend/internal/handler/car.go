package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	"haupcar/backend/internal/domain"
	"haupcar/backend/internal/service"
	"haupcar/backend/internal/utils"
)

type CarHandler struct{ service *service.CarService }

func NewCar(service *service.CarService) *CarHandler {
	return &CarHandler{service: service}
}

func (h *CarHandler) List(c *echo.Context) error {
	cars, err := h.service.List(c.Request().Context())
	if err != nil {
		return utils.Failure(c, err)
	}
	response := make([]carResponse, 0, len(cars))
	for _, car := range cars {
		response = append(response, toCarResponse(car))
	}
	return c.JSON(http.StatusOK, response)
}
func (h *CarHandler) Get(c *echo.Context) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return utils.Failure(c, err)
	}
	result, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return utils.Failure(c, err)
	}
	return c.JSON(http.StatusOK, toCarResponse(result))
}
func (h *CarHandler) Create(c *echo.Context) error {
	var input carRequest
	if err := utils.DecodeJSON(c.Request().Body, &input); err != nil {
		return utils.Failure(c, err)
	}
	result, err := h.service.Create(c.Request().Context(), domain.CarInput(input))
	if err != nil {
		return utils.Failure(c, err)
	}
	c.Response().Header().Set("Location", "/api/cars/"+strconv.FormatUint(result.ID, 10))
	return c.JSON(http.StatusCreated, toCarResponse(result))
}
func (h *CarHandler) Update(c *echo.Context) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return utils.Failure(c, err)
	}
	var input carRequest
	if err := utils.DecodeJSON(c.Request().Body, &input); err != nil {
		return utils.Failure(c, err)
	}
	result, err := h.service.Update(c.Request().Context(), id, domain.CarInput(input))
	if err != nil {
		return utils.Failure(c, err)
	}
	return c.JSON(http.StatusOK, toCarResponse(result))
}
func (h *CarHandler) Delete(c *echo.Context) error {
	id, err := utils.ParseID(c)
	if err != nil {
		return utils.Failure(c, err)
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return utils.Failure(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}
