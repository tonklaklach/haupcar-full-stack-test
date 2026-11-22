package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"

	"haupcar/backend/internal/database"
)

const carsPath = "/api/cars"

type carPayload struct {
	ID           uint64    `json:"id"`
	Registration string    `json:"registration"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func newTestApp(t *testing.T) (http.Handler, *gorm.DB) {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "cars.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := database.Close(db); err != nil {
			t.Error(err)
		}
	})
	app, err := newApp(db)
	if err != nil {
		t.Fatal(err)
	}
	return app, db
}

func request(t *testing.T, router http.Handler, method, path, jsonBody string, want int) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	if response.Code != want {
		t.Fatalf("%s %s: status %d, want %d. Body: %s", method, path, response.Code, want, response.Body)
	}
	return response
}

func TestCarLifecycle(t *testing.T) {
	router, _ := newTestApp(t)
	initial := request(t, router, http.MethodGet, carsPath, "", http.StatusOK)
	if strings.TrimSpace(initial.Body.String()) != "[]" {
		t.Fatalf("empty list: %s", initial.Body)
	}

	created := request(t, router, http.MethodPost, carsPath,
		`{"registration":" 1กข 2345 ","brand":"Toyota","model":"Yaris","notes":"รถสำรอง"}`, http.StatusCreated)
	location := created.Header().Get("Location")
	if location == "" {
		t.Fatal("create response has no Location")
	}

	list := request(t, router, http.MethodGet, carsPath, "", http.StatusOK)
	var cars []carPayload
	if err := json.Unmarshal(list.Body.Bytes(), &cars); err != nil {
		t.Fatal(err)
	}
	if len(cars) != 1 || cars[0].ID == 0 || cars[0].Registration != "1กข-2345" ||
		cars[0].Brand != "Toyota" || cars[0].Model != "Yaris" || cars[0].Notes != "รถสำรอง" ||
		cars[0].CreatedAt.IsZero() || cars[0].UpdatedAt.IsZero() {
		t.Fatalf("unexpected cars: %+v", cars)
	}
	if want := carsPath + "/" + strconv.FormatUint(cars[0].ID, 10); location != want {
		t.Fatalf("Location = %q, want %q", location, want)
	}

	request(t, router, http.MethodPut, location,
		`{"registration":"1กข 2345","brand":"Toyota","model":"Yaris Ativ","notes":""}`, http.StatusOK)
	found := request(t, router, http.MethodGet, location, "", http.StatusOK)
	var updated carPayload
	if err := json.Unmarshal(found.Body.Bytes(), &updated); err != nil {
		t.Fatal(err)
	}
	if updated.Model != "Yaris Ativ" || updated.Notes != "" {
		t.Fatalf("unexpected updated car: %+v", updated)
	}

	request(t, router, http.MethodDelete, location, "", http.StatusNoContent)
	request(t, router, http.MethodGet, location, "", http.StatusNotFound)
}

func TestCarInputErrors(t *testing.T) {
	router, _ := newTestApp(t)
	request(t, router, http.MethodPost, carsPath,
		`{"registration":" ขก 1234 ","brand":"Honda","model":"City"}`, http.StatusCreated)
	request(t, router, http.MethodPost, carsPath,
		`{"registration":"ขก1234","brand":"Toyota","model":"Yaris"}`, http.StatusConflict)
	second := request(t, router, http.MethodPost, carsPath,
		`{"registration":"4ขห-798","brand":"Toyota","model":"Yaris"}`, http.StatusCreated)
	request(t, router, http.MethodPut, second.Header().Get("Location"),
		`{"registration":"ขก1234","brand":"Toyota","model":"Yaris"}`, http.StatusConflict)
	request(t, router, http.MethodPost, carsPath,
		`{"registration":"  ","brand":"Toyota","model":"Yaris"}`, http.StatusBadRequest)
	request(t, router, http.MethodPost, carsPath, `{}{}`, http.StatusBadRequest)
	for _, invalid := range []string{"4ข-798", "ขก-12345", "AB-123", "ข-ก1234", "ขก--1234"} {
		response := request(t, router, http.MethodPost, carsPath,
			`{"registration":"`+invalid+`","brand":"Toyota","model":"Yaris"}`, http.StatusBadRequest)
		if !strings.Contains(response.Body.String(), "รูปแบบทะเบียนรถไม่ถูกต้อง") {
			t.Fatalf("invalid plate %q: %s", invalid, response.Body)
		}
	}
	request(t, router, http.MethodGet, carsPath+"/invalid", "", http.StatusBadRequest)
	request(t, router, http.MethodDelete, carsPath+"/999", "", http.StatusNotFound)
}

func TestLegacyRegistrationKey(t *testing.T) {
	router, db := newTestApp(t)
	if err := db.Exec("INSERT INTO cars (registration, registration_key, brand, model) VALUES (?, ?, ?, ?)",
		"ขก-1234", "ขก-1234", "Honda", "City").Error; err != nil {
		t.Fatal(err)
	}
	request(t, router, http.MethodPost, carsPath,
		`{"registration":"ขก 1234","brand":"Toyota","model":"Yaris"}`, http.StatusConflict)

	other := request(t, router, http.MethodPost, carsPath,
		`{"registration":"4ขห798","brand":"Toyota","model":"Yaris"}`, http.StatusCreated)
	request(t, router, http.MethodPut, other.Header().Get("Location"),
		`{"registration":"ขก1234","brand":"Toyota","model":"Yaris"}`, http.StatusConflict)
	current := request(t, router, http.MethodGet, other.Header().Get("Location"), "", http.StatusOK)
	var car carPayload
	if err := json.Unmarshal(current.Body.Bytes(), &car); err != nil || car.Registration != "4ขห-798" {
		t.Fatalf("failed update changed the car: %+v, %v", car, err)
	}
	request(t, router, http.MethodPut, carsPath+"/1",
		`{"registration":"ขก1234","brand":"Honda","model":"City"}`, http.StatusOK)
}
