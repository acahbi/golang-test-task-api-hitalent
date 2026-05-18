package domain

import "context"

// Объявляем домены
// Описываем интерфейсы взаимодействия

const DepartmentTableName = "department"
const EmployeeTableName = "employee"

type DepartmentRepository interface {
	Create(ctx context.Context, department *Department) (*Department, error)
	CreateEmployee(ctx context.Context, department *Employee) (*Employee, error)
	GetById(ctx context.Context, department *Department, depth uint8, includeEmployees bool) (*OutputDepartmentDetails, error)
	Moved(ctx context.Context, department *Department) (*Department, error)
	Delete(ctx context.Context, id uint64, reassignToDepartmentId uint64) error
}
