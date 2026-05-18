package domain

import "time"

// Объявляем домены
// Описываем модели

type Department struct {
	Id        uint64    `json:"id"`
	Name      string    `json:"name"`
	ParentId  *uint64   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Employee struct {
	Id           uint64     `json:"id"`
	DepartmentId uint64     `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Department) TableName() string {
	return DepartmentTableName
}

func (Employee) TableName() string {
	return EmployeeTableName
}

type OutputDepartmentDetails struct {
	Department Department
	Employees  []Employee
	Children   []Department
}
