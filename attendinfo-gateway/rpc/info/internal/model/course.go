package model

type Course struct {
	Id        uint32
	WeekStart uint32
	WeekEnd   uint32
}

func GetCourseinfo(id uint32) (c *Course, err error) {
	err = DB.Table("course").Where("id=?", id).First(&c).Error
	return
}
