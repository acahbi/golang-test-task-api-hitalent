package transport

import (
	"net/http"

	department "example.com/golang-test-task-api-hitalent/internal/transport/handlers"
	swagger "example.com/golang-test-task-api-hitalent/internal/transport/swagger"
)

// Описываем точки входа для наших хендлеров

func NewRouter(department *department.Handler, swaggerHandler *swagger.Handler) *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc("GET /swagger/openapi.json", swaggerHandler.ServeSpec)
	router.HandleFunc("GET /swagger/", swaggerHandler.ServeUI)
	router.HandleFunc("GET /swagger", swaggerHandler.RedirectToUI)

	router.HandleFunc("POST /departments/", department.Create)
	router.HandleFunc("POST /departments/{id}/employees/", department.CreateEmployee)
	router.HandleFunc("GET /departments/{id}", department.GetById)
	router.HandleFunc("PATCH /departments/{id}", department.MovedById)
	router.HandleFunc("DELETE /departments/{id}", department.DeleteById)

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	return router
}
