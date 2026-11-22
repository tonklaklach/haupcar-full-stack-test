package utils

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"

	"haupcar/backend/internal/domain"
)

func Failure(c *echo.Context, err error) error {
	status, message := http.StatusInternalServerError, "เกิดข้อผิดพลาดในระบบ กรุณาลองใหม่"
	var validation *domain.ValidationError
	switch {
	case errors.Is(err, echo.ErrStatusRequestEntityTooLarge):
		status, message = http.StatusRequestEntityTooLarge, "ขนาดข้อมูลเกินกำหนด"
	case errors.As(err, &validation):
		status, message = http.StatusBadRequest, validation.Message
	case errors.Is(err, domain.ErrNotFound):
		status, message = http.StatusNotFound, "ไม่พบข้อมูลรถยนต์"
	case errors.Is(err, domain.ErrDuplicate):
		status, message = http.StatusConflict, "ทะเบียนรถนี้มีอยู่ในระบบแล้ว"
	default:
		slog.Error("car request failed", "error", err)
	}
	return c.JSON(status, map[string]string{"message": message})
}
