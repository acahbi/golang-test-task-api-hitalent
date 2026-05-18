package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	domain "example.com/golang-test-task-api-hitalent/internal/domain/department"
	service "example.com/golang-test-task-api-hitalent/internal/service"
)

// Наши хендлеры для управления запросами

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func decodeJSON(r *http.Request, dest any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return err
	}

	return nil
}

func getIDFromRequest(r *http.Request) (uint64, error) {
	rawID := r.PathValue("id")
	if rawID == "" {
		return 0, errors.New("missing task id")
	}

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid task id")
	}

	if id <= 0 {
		return 0, errors.New("invalid task id")
	}

	return id, nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var department *domain.Department
	if err := decodeJSON(r, &department); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.Create(r.Context(), service.InputCreateDepartment{Name: department.Name, ParentId: department.ParentId})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	id, errId := getIDFromRequest(r)
	if errId != nil {
		writeError(w, http.StatusBadRequest, errId)
		return
	}

	var employee *domain.Employee
	if err := decodeJSON(r, &employee); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.CreateEmployee(r.Context(), service.InputCreateEmployee{DepartmentId: id, FullName: employee.FullName, Position: employee.Position, HiredAt: employee.HiredAt})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	id, errId := getIDFromRequest(r)
	if errId != nil {
		writeError(w, http.StatusBadRequest, errId)
		return
	}

	query := r.URL.Query()

	depth, errD := strconv.Atoi(query.Get("depth"))
	if errD != nil || depth < 1 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}

	includeEmployees, errI := strconv.ParseBool(query.Get("include_employees"))
	if errI != nil {
		includeEmployees = true
	}

	res, err := h.service.GetById(r.Context(), service.InputGetDepartment{Id: id, Depth: uint8(depth), IncludeEmployees: includeEmployees})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) MovedById(w http.ResponseWriter, r *http.Request) {
	id, errId := getIDFromRequest(r)
	if errId != nil {
		writeError(w, http.StatusBadRequest, errId)
		return
	}

	var department *domain.Department
	if err := decodeJSON(r, &department); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.Moved(r.Context(), service.InputMovedDepartment{Id: id, Name: department.Name, ParentId: department.ParentId})
	if err != nil {
		if err.Error() == "A cycle has been detected" {
			writeError(w, http.StatusConflict, err)
		} else {
			writeError(w, http.StatusBadRequest, err)
		}
		return
	}
	writeJSON(w, http.StatusNoContent, res)
}

func (h *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id, errId := getIDFromRequest(r)
	if errId != nil {
		writeError(w, http.StatusBadRequest, errId)
		return
	}

	query := r.URL.Query()

	mode := query.Get("mode")
	reassign_to_department_id, errReassign := strconv.Atoi(query.Get("reassign_to_department_id"))
	if errReassign != nil {
		reassign_to_department_id = 0
	}

	err := h.service.Delete(r.Context(), service.InputDeleteDepartment{Id: id, Mode: mode, ReassignToDepartmentId: uint64(reassign_to_department_id)})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)

}
