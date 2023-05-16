package model

import "gorm.io/gorm"

type AttendTable struct {
	gorm.Model

	CourseID      string
	CourseName    string
	Week          uint
	Teacher       string
	University    string
	Unpresent     uint
	Unpresenter   string
	UnpresenterID string
}
