package model

import "gorm.io/gorm"

type LeaveTable struct {
	gorm.Model
	//学生信息
	StudentID    uint64 `gorm:"not null"`
	StudentName  string
	StudentMajor string
	StudentClass string
	College      string
	//手机号联系方式
	ContactPhone string `gorm:"not null"`
	//紧急联系人 如果为请假条该项不能为空
	EmergencyName  string
	EmergencyPhone string
	//学校名称
	SchoolName string `gorm:"not null"`
	//辅导员信息
	CounsellorName string `gorm:"not null"`
	CounsellorID   uint64 `gorm:"not null"`
	//申诉-请假理由
	LeaveReason string
	//申诉-请假课程
	CourseName string `gorm:"not null"`
	CourseID   uint64 `gorm:"not null"`
	//课时起止
	LeaveCourseFrom int `gorm:"not null"`
	LeaveCourseTo   int `gorm:"not null"`
	//申诉表-请假条区分
	TagAs   uint `gorm:"default:1"` //默认为请假条
	IsAudit uint `gorm:"default:1"` //1为未审核 2为审核已通过
}
