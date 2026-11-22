package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	registrationPattern = regexp.MustCompile(`^([0-9]?[ก-ฮ]{2})-?([0-9]{1,4})$`)
	ErrNotFound         = errors.New("car not found")
	ErrDuplicate        = errors.New("registration already exists")
)

type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }

type Car struct {
	ID           uint64
	Registration string
	Brand        string
	Model        string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CarInput struct {
	Registration string
	Brand        string
	Model        string
	Notes        string
}

func RegistrationKey(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || r == '-' {
			return -1
		}
		return r
	}, value)
}

func (in CarInput) Normalize() (CarInput, error) {
	in.Registration = strings.TrimSpace(in.Registration)
	in.Brand = strings.TrimSpace(in.Brand)
	in.Model = strings.TrimSpace(in.Model)
	in.Notes = strings.TrimSpace(in.Notes)
	fields := []struct {
		name, value string
		limit       int
	}{
		{"ทะเบียนรถ", in.Registration, 30}, {"ยี่ห้อ", in.Brand, 60}, {"รุ่น", in.Model, 80},
	}
	for _, field := range fields {
		if field.value == "" {
			return in, &ValidationError{Message: "กรุณาระบุ" + field.name}
		}
		if utf8.RuneCountInString(field.value) > field.limit {
			return in, &ValidationError{Message: field.name + "ยาวเกินกำหนด"}
		}
	}
	if utf8.RuneCountInString(in.Notes) > 1000 {
		return in, &ValidationError{Message: "หมายเหตุต้องไม่เกิน 1,000 ตัวอักษร"}
	}
	compact := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, in.Registration)
	parts := registrationPattern.FindStringSubmatch(compact)
	if parts == nil {
		return in, &ValidationError{Message: "รูปแบบทะเบียนรถไม่ถูกต้อง เช่น 4ขห-798 หรือ ขก-1234"}
	}
	in.Registration = parts[1] + "-" + parts[2]
	return in, nil
}
