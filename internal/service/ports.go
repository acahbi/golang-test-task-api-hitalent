package service

import (
	"time"
)

const (
	CASCADE  string = "cascade"
	REASSING string = "reassing"
)

type InputCreateDepartment struct {
	Name     string
	ParentId *uint64
}

type InputCreateEmployee struct {
	Id           uint64
	FullName     string
	Position     string
	HiredAt      *time.Time
	DepartmentId uint64
}

type InputGetDepartment struct {
	Id               uint64
	Depth            uint8
	IncludeEmployees bool
}

type InputMovedDepartment struct {
	Id       uint64
	Name     string
	ParentId *uint64
}

type InputDeleteDepartment struct {
	Id                     uint64
	Mode                   string
	ReassignToDepartmentId uint64
}
