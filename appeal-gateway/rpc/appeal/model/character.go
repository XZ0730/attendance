package model

type Character struct {
	Id         uint64
	Name       string
	University string
	Code       string
	College    string
	Major      string `gorm:"column:major"`
	ClassNum   uint
	Grade      int
}
