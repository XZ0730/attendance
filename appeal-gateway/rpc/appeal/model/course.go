package model

type Course struct {
	Id           int64
	CourseId     string
	Name         string
	WeekStart    uint32
	WeekEnd      uint32
	Term         uint32
	SemesterYear uint32
	University   string
}
