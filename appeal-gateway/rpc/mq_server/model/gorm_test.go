package model

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)



func TestCosume(t *testing.T) {
	Results := make([]*Result, 0)
	initto()
	DB.Table("course_group").
		Select("character_msg.code as code,character_msg.name as student_name,course_group.course_id as course_id").Joins("left join character_msg on course_group.student_id=character_msg.code where course_id=?", "10000").Scan(&Results)
	for _, v := range Results {
		fmt.Println("result:", *v)
	}
}
func initto() {
	ormLogger := logger.Default
	datasource := "attendanceSystem:147258@tcp(47.113.216.236:3306)/attendancesystem?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       datasource,
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
		Sources:  []gorm.Dialector{mysql.Open(datasource)},
		Replicas: []gorm.Dialector{mysql.Open(datasource)},
		Policy:   dbresolver.RandomPolicy{},
	}))
	DB = db
}
