package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"info/attendinfo"
	"info/internal/model"
	"info/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCourseAttInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCourseAttInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseAttInfoLogic {
	return &GetCourseAttInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCourseAttInfoLogic) GetCourseAttInfo(in *attendinfo.CourseInfoRequest) (*attendinfo.CourseInfoResponse, error) {
	// todo: add your logic here and delete this line
	//传入参数: University Courseid  Course_main
	//先查询课程信息
	course, err := model.GetCourseinfo(in.GetCourse_Main())
	if err != nil {
		return &attendinfo.CourseInfoResponse{
			Status:  36677,
			Message: "课程信息获取失败",
			Error:   err.Error(),
		}, nil
	}
	courseAttinfo := make([]*attendinfo.CourseAttInfo, 0)
	//根据课程起始周进行range
	week, err2 := l.svcCtx.RDB.HGet(l.ctx, in.GetUniversity(), strconv.Itoa(int(in.GetCourse_Main()))).Result()
	if err2 != nil {
		return &attendinfo.CourseInfoResponse{
			Status:  36676,
			Message: "课程redis出错",
			Error:   err2.Error(),
		}, nil
	}
	weeklast, err3 := strconv.Atoi(week)
	if err3 != nil {
		return &attendinfo.CourseInfoResponse{
			Status:  36675,
			Message: "strconv转换错误",
			Error:   err3.Error(),
		}, nil
	}
	for i := course.WeekStart; i <= uint32(weeklast); i++ {
		//查询redis中此周有无点名
		staus, err2 := l.svcCtx.RDB7.HGet(l.ctx, strconv.Itoa(int(course.Id)), strconv.Itoa(int(i))).Result()
		if staus == "0" || err2 != nil {
			cAinfo := &attendinfo.CourseAttInfo{
				Week:           i,
				AttendRate:     1,
				UnpresentCount: 0,
				UnpresentInfos: nil,
			}
			courseAttinfo = append(courseAttinfo, cAinfo)
			continue
		}
		//点名了就查询redis中的信息
		mres := make(map[string]*model.Result, 0)
		res, err4 := l.svcCtx.RDB6.HGet(l.ctx, strconv.Itoa(int(in.GetCourse_Main())), strconv.Itoa(int(i))).Result()
		if err4 != nil {
			//查询不到就去查询mysql中的信息

			//还查不到则未点名
			cAinfo := &attendinfo.CourseAttInfo{
				Week:           i,
				AttendRate:     1,
				UnpresentCount: 0,
				UnpresentInfos: nil,
			}
			courseAttinfo = append(courseAttinfo, cAinfo)
			continue
		}
		err5 := json.Unmarshal([]byte(res), &mres)
		if err5 != nil {
			fmt.Println("err:", err5.Error())
		}
		um := make([]*attendinfo.UnpresentModel, 0)
		cnt := 0.0
		for _, v := range mres {
			if v.MissAttend == 1 {
				cnt++
				aum := &attendinfo.UnpresentModel{
					Code:       v.Code,
					Name:       v.StudentName,
					MissAttend: uint32(v.MissAttend),
				}
				um = append(um, aum)
			}
		}
		//查到了就构造list
		//查询attendtable 获取信息
		cAinfo := &attendinfo.CourseAttInfo{
			Week:           i,
			AttendRate:     math.Round(float64(len(mres))-cnt) / float64(len(mres)),
			UnpresentCount: uint32(cnt),
			UnpresentInfos: um,
		}
		// fmt.Println("len(um):", len(um))
		// fmt.Println("cnt:", cnt)
		// fmt.Println("rate:", cAinfo.AttendRate)
		courseAttinfo = append(courseAttinfo, cAinfo)
	}

	return &attendinfo.CourseInfoResponse{
		Status:         200,
		Message:        "successful",
		CourseAttInfos: courseAttinfo,
		Total:          uint32(len(courseAttinfo)),
	}, nil
}
func ForGetCourseAttInfo() {

}
