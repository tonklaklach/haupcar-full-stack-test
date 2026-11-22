package handler

import (
	"time"

	"haupcar/backend/internal/domain"
)

type carRequest struct {
	Registration string `json:"registration"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Notes        string `json:"notes"`
}

type carResponse struct {
	ID           uint64    `json:"id"`
	Registration string    `json:"registration"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func toCarResponse(car domain.Car) carResponse {
	return carResponse{
		ID:           car.ID,
		Registration: car.Registration,
		Brand:        car.Brand,
		Model:        car.Model,
		Notes:        car.Notes,
		CreatedAt:    car.CreatedAt,
		UpdatedAt:    car.UpdatedAt,
	}
}
