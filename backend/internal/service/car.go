package service

import (
	"context"

	"haupcar/backend/internal/domain"
)

type CarService struct{ repository CarRepository }

func NewCar(repository CarRepository) *CarService {
	return &CarService{repository: repository}
}

func (s *CarService) List(ctx context.Context) ([]domain.Car, error) {
	return s.repository.List(ctx)
}

func (s *CarService) Get(ctx context.Context, id uint64) (domain.Car, error) {
	return s.repository.Get(ctx, id)
}

func (s *CarService) Create(ctx context.Context, input domain.CarInput) (domain.Car, error) {
	input, err := input.Normalize()
	if err != nil {
		return domain.Car{}, err
	}
	return s.repository.Create(ctx, input)
}

func (s *CarService) Update(ctx context.Context, id uint64, input domain.CarInput) (domain.Car, error) {
	input, err := input.Normalize()
	if err != nil {
		return domain.Car{}, err
	}
	return s.repository.Update(ctx, id, input)
}

func (s *CarService) Delete(ctx context.Context, id uint64) error {
	return s.repository.Delete(ctx, id)
}
