package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"haupcar/backend/internal/domain"

	"gorm.io/gorm"
)

type carRecord struct {
	ID              uint64 `gorm:"primaryKey"`
	Registration    string
	RegistrationKey string `gorm:"uniqueIndex;not null"`
	Brand           string
	Model           string
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (carRecord) TableName() string { return "cars" }

func (r carRecord) domain() domain.Car {
	return domain.Car{
		ID:           r.ID,
		Registration: r.Registration,
		Brand:        r.Brand,
		Model:        r.Model,
		Notes:        r.Notes,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate() error {
	if err := r.db.AutoMigrate(&carRecord{}); err != nil {
		return fmt.Errorf("migrate cars: %w", err)
	}
	return nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return domain.ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return domain.ErrDuplicate
	default:
		return err
	}
}
func (r *Repository) List(ctx context.Context) ([]domain.Car, error) {
	var records []carRecord
	if err := r.db.WithContext(ctx).Order("id DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	cars := make([]domain.Car, 0, len(records))
	for _, record := range records {
		cars = append(cars, record.domain())
	}
	return cars, nil
}
func (r *Repository) Get(ctx context.Context, id uint64) (domain.Car, error) {
	var record carRecord
	err := r.db.WithContext(ctx).First(&record, id).Error
	return record.domain(), mapError(err)
}
func (r *Repository) Create(ctx context.Context, input domain.CarInput) (domain.Car, error) {
	record := carRecord{
		Registration:    input.Registration,
		RegistrationKey: domain.RegistrationKey(input.Registration),
		Brand:           input.Brand,
		Model:           input.Model,
		Notes:           input.Notes,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := checkLegacyRegistration(tx, input.Registration, 0); err != nil {
			return err
		}
		return mapError(tx.Create(&record).Error)
	})
	if err != nil {
		return domain.Car{}, err
	}
	return record.domain(), nil
}
func (r *Repository) Update(ctx context.Context, id uint64, input domain.CarInput) (domain.Car, error) {
	var result domain.Car
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"registration":     input.Registration,
			"registration_key": domain.RegistrationKey(input.Registration),
			"brand":            input.Brand,
			"model":            input.Model,
			"notes":            input.Notes,
		}
		query := tx.Model(&carRecord{}).Where("id = ?", id).Updates(updates)
		if query.Error != nil {
			return mapError(query.Error)
		}
		if query.RowsAffected == 0 {
			return domain.ErrNotFound
		}
		if err := checkLegacyRegistration(tx, input.Registration, id); err != nil {
			return err
		}
		var record carRecord
		if err := tx.First(&record, id).Error; err != nil {
			return mapError(err)
		}
		result = record.domain()
		return nil
	})
	return result, err
}
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	query := r.db.WithContext(ctx).Delete(&carRecord{}, id)
	if query.Error != nil {
		return mapError(query.Error)
	}
	if query.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func checkLegacyRegistration(tx *gorm.DB, registration string, excludeID uint64) error {
	var count int64
	err := tx.Model(&carRecord{}).
		Where("registration_key = ? AND id <> ?", registration, excludeID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count != 0 {
		return domain.ErrDuplicate
	}
	return nil
}
