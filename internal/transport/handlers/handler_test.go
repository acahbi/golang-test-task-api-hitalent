package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"testing"

	domain "example.com/golang-test-task-api-hitalent/internal/domain/department"
	connection "example.com/golang-test-task-api-hitalent/internal/infrastructure/postgres"
	repo "example.com/golang-test-task-api-hitalent/internal/repository/postgres"
	service "example.com/golang-test-task-api-hitalent/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type config struct {
	dbConnect string
}

func loadConfig() config {
	cfg := config{
		dbConnect: os.Getenv("DB_CONNECT"),
	}

	if cfg.dbConnect == "" {
		panic(fmt.Errorf("DB_CONNECT not received"))
	}

	return cfg
}

func getHandler() (*Handler, *gorm.DB, error) {
	// Получаем конфигурацию подключений
	cfg := loadConfig()
	/*
		cfg := config{
			//httpAddr:  ":8080",
			dbConnect: "postgres://postgres:postgres1@localhost:5432/postgres?sslmode=disable",
		}
		//*/

	// Подписываемся на сигналы syscall.SIGINT (завершение процесса CTRL+C), syscall.SIGTERM (Завершение, например, по kill pid)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Получаем открытое подключение к БД
	conn, err := connection.NewConnection(ctx, cfg.dbConnect)
	if err != nil {
		return nil, nil, fmt.Errorf("There is a connection error")
	}

	departmentRepository := repo.NewRepository(conn)
	service := service.NewService(departmentRepository)
	handler := NewHandler(service)

	return handler, conn, nil
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantCode     int
		wantResponse *domain.Department
	}{
		{name: "Correct", body: `{"name": "TestDepartment", "parent_id": null}`, wantCode: http.StatusCreated, wantResponse: &domain.Department{}},
		{name: "NoCorrect", body: `{"name": "TestNoValidDepartment", "parent_id": 99999}`, wantCode: http.StatusBadRequest, wantResponse: nil},
	}

	handler, conn, err := getHandler()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Закрываем подключение в отложенной функции
	defer func() {
		db, err := conn.DB()
		if err != nil {
			return
		}

		_ = db.Close()
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/departments/", bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler.Create(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestCreateEmployee(t *testing.T) {
	tests := []struct {
		name         string
		parentId     string // но по факту у uint64
		body         string
		wantCode     int
		wantResponse *domain.Employee
	}{
		{name: "Correct", parentId: "1", body: `{"full_name": "TestEmployee", "position": "manager", "hired_at": "2026-03-24T00:00:00Z"}`, wantCode: http.StatusCreated, wantResponse: &domain.Employee{}},
		{name: "NoCorrect", parentId: "999999999", body: `{"full_name": "TestEmployee", "position": "manager", "hired_at": "2026-03-24T00:00:00Z"}`, wantCode: http.StatusBadRequest, wantResponse: nil},
	}

	handler, conn, err := getHandler()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Закрываем подключение в отложенной функции
	defer func() {
		db, err := conn.DB()
		if err != nil {
			return
		}

		_ = db.Close()
	}()

	router := http.NewServeMux()
	router.HandleFunc("/departments/{id}/employees/", handler.CreateEmployee)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/departments/" + tt.parentId + "/employees/"
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestGetById(t *testing.T) {
	tests := []struct {
		name         string
		departmentId string // но по факту у uint64
		query        url.Values
		wantCode     int
		wantResponse *domain.OutputDepartmentDetails
	}{
		{name: "Correct", departmentId: "1", query: url.Values{"depth": {"3"}, "include_employees": {"true"}}, wantCode: http.StatusOK, wantResponse: &domain.OutputDepartmentDetails{}},
		{name: "NoCorrect", departmentId: "-999999999", query: url.Values{"depth": {"1"}, "include_employees": {"false"}}, wantCode: http.StatusBadRequest, wantResponse: nil},
	}

	handler, conn, err := getHandler()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Закрываем подключение в отложенной функции
	defer func() {
		db, err := conn.DB()
		if err != nil {
			return
		}

		_ = db.Close()
	}()

	router := http.NewServeMux()
	router.HandleFunc("/departments/{id}", handler.GetById)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/departments/" + tt.departmentId
			u := url.URL{
				Path:     path,
				RawQuery: tt.query.Encode(),
			}
			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestMovedById(t *testing.T) {
	tests := []struct {
		name         string
		departmentId string // но по факту у uint64
		body         string
		wantCode     int
		wantResponse *domain.Department
	}{
		{name: "Correct", departmentId: "1", body: `{"name": "TestDepartmentUpdated", "parent_id": null}`, wantCode: http.StatusNoContent, wantResponse: &domain.Department{}},
		{name: "NoCorrect", departmentId: "999999999", body: `{"full_name": "TestDepartmentFailUpdated", "parent_id": "7410"}`, wantCode: http.StatusBadRequest, wantResponse: nil},
	}

	handler, conn, err := getHandler()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Закрываем подключение в отложенной функции
	defer func() {
		db, err := conn.DB()
		if err != nil {
			return
		}

		_ = db.Close()
	}()

	router := http.NewServeMux()
	router.HandleFunc("/departments/{id}", handler.MovedById)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/departments/" + tt.departmentId
			req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestDeleteById(t *testing.T) {
	tests := []struct {
		name         string
		departmentId string // но по факту у uint64
		query        url.Values
		wantCode     int
		wantResponse *domain.Department
	}{
		{name: "Correct", departmentId: "5", query: url.Values{"mode": {"cascade"}}, wantCode: http.StatusNoContent, wantResponse: &domain.Department{}},
	}

	handler, conn, err := getHandler()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Закрываем подключение в отложенной функции
	defer func() {
		db, err := conn.DB()
		if err != nil {
			return
		}

		_ = db.Close()
	}()

	router := http.NewServeMux()
	router.HandleFunc("/departments/{id}", handler.DeleteById)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/departments/" + tt.departmentId
			u := url.URL{
				Path:     path,
				RawQuery: tt.query.Encode(),
			}
			req := httptest.NewRequest(http.MethodDelete, u.String(), nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
