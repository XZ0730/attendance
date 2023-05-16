package pkg

import (
	"attend/model"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type CourseTime struct {
	start string
	end   string
}

func NewCouseTime(start, end string) *CourseTime {
	return &CourseTime{
		start: start,
		end:   end,
	}
}

var CourseMap = map[uint]CourseTime{
	1:  *NewCouseTime("8:10", "8:50"),
	2:  *NewCouseTime("9:5", "9:45"),
	3:  *NewCouseTime("10:10:", "10:50"),
	4:  *NewCouseTime("11:5", "11:45"),
	5:  *NewCouseTime("13:50", "14:30"),
	6:  *NewCouseTime("14:45", "15:25"),
	7:  *NewCouseTime("15:40", "16:20"),
	8:  *NewCouseTime("16:35", "17:15"),
	9:  *NewCouseTime("18:50", "19:30"),
	10: *NewCouseTime("19:45", "20:25"),
	11: *NewCouseTime("20:50", "21:30"),
}
var weekDay = map[uint]string{
	1: "Monday",
	2: "Tuesday",
	3: "Wednesday",
	4: "Thursday",
	5: "Friday",
	6: "Saturday",
	7: "Sunday",
}

func JudgeTime(cr *model.Course) bool {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	Weekday := now.Weekday()
	fmt.Println("=========================")
	fmt.Println("cr:", cr)
	if model.RDB5 == nil {
		fmt.Println("nilnilnil")
		return false
	}
	timestamp, err := model.RDB5.Get(model.RDB5.Context(), cr.University).Result()
	fmt.Println("timestamp", timestamp)
	if err != nil {
		fmt.Println("err", err)
		return false
	}
	fmt.Println("====1=====================")
	//当前学期开学时间
	timestamp2, _ := strconv.Atoi(timestamp)
	tbegin := time.Unix(int64(timestamp2), 0)
	dist := now.Sub(tbegin)
	//判断日期--在当前学期内
	fmt.Println("====2=====================")
	weeknow := math.Ceil(dist.Hours() / 24 / 7)
	fmt.Println("crweektype:", cr.WeekType)
	if cr.WeekType == 1 && (int(weeknow)%2) == 0 {
		return false
	} else if cr.WeekDay == 2 && (int(weeknow)%2) == 1 {
		return false
	}

	if weeknow < float64(cr.WeekStart) ||
		weeknow > float64(cr.WeekEnd) ||
		Weekday.String() != weekDay[cr.WeekDay] {
		fmt.Println("weekday:", Weekday.String())
		fmt.Println("weekday1:", weekDay[cr.WeekDay])
		fmt.Println("Cceecec1")
		return false
	}
	//判断当天时间
	//上课前10分钟
	ct := CourseMap[cr.PitchStart]
	coursetime := strings.Split(ct.start, ":")
	hourbegin, _ := strconv.Atoi(coursetime[0])
	minutebegin, err := strconv.Atoi(coursetime[1])
	if err != nil {
		fmt.Println("err:", err)
		fmt.Println("Cceecec2")
		return false
	}
	coursebegin := time.Date(year, month, now.Day(), hourbegin, minutebegin, 0, 0, time.Local)
	d := now.Sub(coursebegin)
	if d.Milliseconds() < 0 {
		return false
	}
	//下课前十五分钟
	ct1 := CourseMap[cr.PitchEnd]
	coursetime = strings.Split(ct1.end, ":")
	hourend, _ := strconv.Atoi(coursetime[0])
	minuteend, err := strconv.Atoi(coursetime[1])
	if err != nil {
		fmt.Println("err:", err)
		fmt.Println("Cceecec3")
		return false
	}

	courseend := time.Date(year, month, now.Day(), hourend, minuteend, 0, 0, time.Local)
	d2 := courseend.Sub(now)
	if d2.Milliseconds() < 0 {
		fmt.Println("Cceecec4")
		return false
	}
	return true
}
