package service

import (
	"context"
	"fmt"

	domain "example.com/golang-test-task-api-hitalent/internal/domain/department"
)

// Объявляем сервис. Тут будем проводить валидацию

type Service struct {
	repo domain.DepartmentRepository
}

func NewService(r domain.DepartmentRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(ctx context.Context, input InputCreateDepartment) (*domain.Department, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}

	model := &domain.Department{
		Name:     input.Name,
		ParentId: input.ParentId,
	}
	res, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) CreateEmployee(ctx context.Context, input InputCreateEmployee) (*domain.Employee, error) {
	if input.FullName == "" {
		return nil, fmt.Errorf("%w: full_name is required", ErrInvalidInput)
	}
	if input.Position == "" {
		return nil, fmt.Errorf("%w: position is required", ErrInvalidInput)
	}

	model := &domain.Employee{
		DepartmentId: input.DepartmentId,
		FullName:     input.FullName,
		Position:     input.Position,
		HiredAt:      input.HiredAt,
	}
	res, err := s.repo.CreateEmployee(ctx, model)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetById(ctx context.Context, input InputGetDepartment) (*domain.OutputDepartmentDetails, error) {
	model := &domain.Department{
		Id: input.Id,
	}
	res, err := s.repo.GetById(ctx, model, input.Depth, input.IncludeEmployees)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Moved(ctx context.Context, input InputMovedDepartment) (*domain.Department, error) {
	model := &domain.Department{
		Id:       input.Id,
		Name:     input.Name,
		ParentId: input.ParentId,
	}
	res, err := s.repo.Moved(ctx, model)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, input InputDeleteDepartment) error {
	if input.Mode != CASCADE && input.Mode != REASSING {
		return fmt.Errorf("%w: mode must be cascade or reassing", ErrInvalidInput)
	}
	if input.Mode == REASSING && input.ReassignToDepartmentId < 1 {
		return fmt.Errorf("%w: reassign_to_department_id is required", ErrInvalidInput)
	}

	err := s.repo.Delete(ctx, input.Id, input.ReassignToDepartmentId)
	if err != nil {
		return err
	}
	return nil
}
