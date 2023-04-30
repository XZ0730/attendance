package model

import (
	"mq_server/internal/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

type ComplainTable struct {
	gorm.Model
}

func Init(c *config.Config) *gorm.DB {

	ormLogger := logger.Default
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       c.Mysql.DataSource,
		DefaultStringSize:         256,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {

		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)                 //设置链接数
	sqlDB.SetMaxOpenConns(1000)                //打开链接数
	sqlDB.SetConnMaxLifetime(time.Second * 30) //链接生命周期

	_ = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(c.Mysql.DataSource)},
		Replicas: []gorm.Dialector{mysql.Open(c.Mysql.DataSource)},
		Policy:   dbresolver.RandomPolicy{},
	}))
	DB = db
	return db
}

type LeaveTable struct {
	gorm.Model
	//学生信息
	StudentID    string `gorm:"not null"`
	StudentName  string
	StudentMajor string
	StudentClass uint
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
	CounsellorID   string `gorm:"not null"`
	//申诉-请假理由
	LeaveReason string
	//申诉-请假课程
	CourseName string `gorm:"not null"`
	CourseID   string `gorm:"not null"`
	//课时起止
	LeaveCourseFrom int `gorm:"not null"`
	LeaveCourseTo   int `gorm:"not null"`
	//申诉表-请假条区分
	TagAs   uint `gorm:"default:1"` //默认为请假条
	IsAudit uint `gorm:"default:1"` //1为未审核 2为审核已通过
}
