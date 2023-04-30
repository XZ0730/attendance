package model

import "gorm.io/gorm"

type ComplainTable struct {
	gorm.Model

	SupervisorID       string `gorm:"not null"`
	SupervisorName     string
	Supervisor_Major   string
	Supervisor_College string
	SchoolName         string

	Reason string

	CounsellorName string
	CounsellorID   string

	StudentID     string
	StudentName   string
	Student_Major string
	College       string
}
