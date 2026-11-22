package service

import (
	"context"

	"haupcar/backend/internal/domain"
)

type CarRepository interface {
	List(context.Context) ([]domain.Car, error)
	Get(context.Context, uint64) (domain.Car, error)
	Create(context.Context, domain.CarInput) (domain.Car, error)
	Update(context.Context, uint64, domain.CarInput) (domain.Car, error)
	Delete(context.Context, uint64) error
}
