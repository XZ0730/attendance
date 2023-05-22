package errorx

const (
	SUCCESS = 200

	FailTOPull         = 30021
	RecordGetError     = 40013
	BusySysError       = 39908
	DeleteError        = 39901
	RepeteAppealError  = 30008
	CourseRepeteError  = 30022
	LastWeekError      = 30329
	RedisInitError     = 30321
	CourseSelectError1 = 30001
	CourseSelectError2 = 30003
	CourseISNoExist    = 30004
	COURSE_ING         = 30014
	AUDIT_ING          = 30034
	JSON_MARSHAL_ERROR = 30015
	MQ_RETURN_ERROR    = 30080
	ComplainPostError  = 40033
	RoleISNOEXIST      = 40003
	AppealIsNoExist    = 40011
	CourseMatchERROR   = 40012
)

var errorMap = map[int]string{
	SUCCESS:            "successful",
	RecordGetError:     "记录查询失败",
	BusySysError:       "系统繁忙稍后重试",
	DeleteError:        "删除错误",
	RepeteAppealError:  "重复申诉",
	CourseRepeteError:  "课时已经结束，或者课时选择超过课时范围",
	LastWeekError:      "上周未点名",
	RedisInitError:     "redis中课程未初始化完成",
	CourseSelectError1: "申诉表:课时选择错误",
	CourseSelectError2: "申诉课时已经结束或者未开始，只能选择最近上完的课时捏",
	CourseISNoExist:    "课程不存在",
	COURSE_ING:         "课程正在上课",
	AUDIT_ING:          "您选择的课时已经审核通过或者还在审核中",
	ComplainPostError:  "投诉发送失败",
	RoleISNOEXIST:      "人员不存在",
	AppealIsNoExist:    "假条不存在",
	CourseMatchERROR:   "课程不匹配",
}

func GetERROR(code int) string {
	return errorMap[code]
}
