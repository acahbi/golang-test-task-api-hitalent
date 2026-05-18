package postgresql

import (
	"context"
	"fmt"

	domain "example.com/golang-test-task-api-hitalent/internal/domain/department"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Функции, необходимые для работы с БД

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	var result *gorm.DB

	if department.ParentId != nil {
		result = r.db.Raw(`
WITH RECURSIVE cte AS (
    SELECT id, parent_id, ARRAY[id] AS visited
      FROM department 
     WHERE id = ?
	   AND parent_id IS NOT NULL
    UNION ALL
    SELECT c.id, c.parent_id, a.visited || c.id AS visited
      FROM department c
      JOIN cte a 
        ON c.parent_id = a.id
	 WHERE NOT (c.id = ANY(a.visited))
)
INSERT INTO department("name", parent_id)
SELECT ?, ?
 WHERE NOT EXISTS(SELECT 1 FROM cte WHERE id = ?)
RETURNING id, "name", parent_id, created_at, updated_at
	`, department.ParentId, department.Name, department.ParentId, department.ParentId).Scan(&department)
	} else {
		result = r.db.Create(&department)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return department, nil
}

func (r *Repository) CreateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
	result := r.db.Create(&employee)
	if result.Error != nil {
		return nil, result.Error
	}

	return employee, nil
}
func (r *Repository) GetById(ctx context.Context, department *domain.Department, depth uint8, includeEmployees bool) (*domain.OutputDepartmentDetails, error) {
	resultDepartment := r.db.First(&department, department.Id)
	if resultDepartment.Error != nil {
		return nil, resultDepartment.Error
	}

	employees := []domain.Employee{}
	if includeEmployees {
		resultEmployee := r.db.Where("department_id=?", department.Id).Order("created_at DESC").Find(&employees)
		if resultEmployee.Error != nil {
			return nil, resultEmployee.Error
		}
	}

	children := []domain.Department{}
	errChildren := r.db.Raw(`
WITH RECURSIVE cte AS (
	SELECT d.*, 0 AS depth
		FROM department d
		WHERE id = ?
	UNION ALL
	SELECT d.*, c.depth + 1
	FROM department d
	JOIN cte c
		ON c.id = d.parent_id
	WHERE c.depth < ?
)
SELECT *
	FROM cte
	WHERE depth > 0;
`, department.Id, depth).Scan(&children)
	if errChildren.Error != nil {
		return nil, errChildren.Error
	}

	return &domain.OutputDepartmentDetails{
		Department: *department,
		Employees:  employees,
		Children:   children,
	}, nil
}

func (r *Repository) Moved(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	var flgCycle bool
	resultCycle := r.db.Raw(`
WITH RECURSIVE cte AS (
    SELECT id, parent_id, ARRAY[id] AS visited
      FROM department 
     WHERE id = ?
	   AND parent_id IS NOT NULL
    UNION ALL
    SELECT c.id, c.parent_id, a.visited || c.id AS visited
      FROM department c
      JOIN cte a 
        ON c.parent_id = a.id
	 WHERE NOT (c.id = ANY(a.visited))
)
SELECT 1 FROM cte WHERE id = ?
	`, department.ParentId, department.ParentId).Scan(&flgCycle)
	if resultCycle.Error != nil {
		return nil, resultCycle.Error
	}

	if flgCycle {
		return nil, fmt.Errorf("A cycle has been detected")
	}

	fields := make(map[string]any)
	fields["name"] = &department.Name
	fields["parent_id"] = &department.ParentId

	var returning domain.Department
	result := r.db.
		Model(&returning).
		Where("id=?", department.Id).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}, {Name: "name"}, {Name: "parent_id"}, {Name: "created_at"}, {Name: "updated_at"}}}).
		Updates(fields)

	if result.Error != nil {
		return nil, result.Error
	}

	return &returning, nil
}
func (r *Repository) Delete(ctx context.Context, id uint64, reassignToDepartmentId uint64) error {
	if reassignToDepartmentId > 0 {
		fields := make(map[string]uint64)
		fields["department_id"] = reassignToDepartmentId

		resultEmployee := r.db.
			Model(&domain.Employee{}).
			Where("department_id=?", id).
			Updates(fields)

		if resultEmployee.Error != nil {
			return resultEmployee.Error
		}
	}

	resultDelete := r.db.Delete(&domain.Department{}, id)
	if resultDelete.Error != nil {
		return resultDelete.Error
	}
	return nil
}
