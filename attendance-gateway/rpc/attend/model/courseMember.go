package model

import "attend/attendservice"

type Member struct {
	Code        string
	StudentName string
	University  string
	CourseId    string
	Week        uint
	Unpresent   bool
}

func GetCourseMem(Couseid, University string) (mem []*attendservice.CourseMember) {
	DB.Table("course_group").
		Select("character_msg.code as code,character_msg.name as student_name,course_group.course_id as course_id,character_msg.university as university").Joins("left join character_msg on course_group.student_id=character_msg.code where course_id=? and course_group.university=?", Couseid, University).
		Scan(&mem)
	return
}
