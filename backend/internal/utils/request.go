package utils

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/labstack/echo/v5"

	"haupcar/backend/internal/domain"
)

func ParseID(c *echo.Context) (uint64, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 63)
	if err != nil || id == 0 {
		return 0, &domain.ValidationError{Message: "รหัสรถไม่ถูกต้อง"}
	}
	return id, nil
}

func DecodeJSON(body io.Reader, dst any) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, echo.ErrStatusRequestEntityTooLarge) {
			return err
		}
		return &domain.ValidationError{Message: "รูปแบบ JSON ไม่ถูกต้อง"}
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if errors.Is(err, echo.ErrStatusRequestEntityTooLarge) {
			return err
		}
		return &domain.ValidationError{Message: "ต้องส่ง JSON เพียงหนึ่งรายการ"}
	}
	return nil
}
