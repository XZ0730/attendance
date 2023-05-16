package model

type Course struct {
	Id           int64
	University   string
	PitchStart   uint
	PitchEnd     uint
	WeekStart    uint
	WeekEnd      uint
	WeekDay      uint
	WeekType     uint
	SemesterYear uint
	Term         uint
}
