package model

import (
	"appeal/internal/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

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
	Migrate(db)
	return db
}
func Migrate(DB *gorm.DB) {
	err := DB.Set(`gorm:table_options`, "charset=utf8mb4").
		AutoMigrate(&LeaveTable{}, &ComplainTable{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
