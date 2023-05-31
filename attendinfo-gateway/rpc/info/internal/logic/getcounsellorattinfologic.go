package logic

import (
	"context"
	"fmt"
	"sync"

	"info/attendinfo"
	"info/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCounsellorAttInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCounsellorAttInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCounsellorAttInfoLogic {
	return &GetCounsellorAttInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCounsellorAttInfoLogic) GetCounsellorAttInfo(in *attendinfo.CounsellorAttInfoReq) (*attendinfo.CounsellorAttInfoRsp, error) {
	// todo: add your logic here and delete this line
	majorRate := make([]*attendinfo.MajorRate, 0)
	courseRate := make([]*attendinfo.CourseRate, 0)
	var courseRateMap sync.Map
	fmt.Println("------测试5-------------")
	for _, v := range in.GetMajorList() {
		rate1 := &attendinfo.MajorRate{}
		rate1.Major = v.GetMajor()
		fmt.Println("------测试6-------------")
		for _, u := range v.GetCouseMain() {
			_, ok := courseRateMap.Load(u)
			if ok {
				continue
			}
			cmr := &CounsellorModelReq{
				Req: l,
				Req1: &GetCourseAttInfoLogic{
					ctx:    l.ctx,
					svcCtx: l.svcCtx,
					Logger: l.Logger,
				},
				RateMap:    &courseRateMap,
				CourseMain: u,
				University: v.GetUniversity(),
			}
			wait.Add(1)
			ants_Pool.Invoke(cmr)
			//返回最近结束周的考勤率即可
			//考勤率存入map中
		}
		// fmt.Println("------测试7-------------")
		wait.Wait()
		rate := 0.0
		cnt := 0.0
		okcnt := 0.0
		judge := make(map[int]bool, 0)
		for _, v := range v.GetCouseMain() {

			value, ok := courseRateMap.Load(v)
			if ok {
				cnt++
				// fmt.Println("value:", value.(*RateModel).AverageRate)
				rate += value.(*RateModel).AverageRate
				if !judge[int(v)] { //去重
					rate2 := &attendinfo.CourseRate{}
					rate2.AttendRate = value.(*RateModel).LastWeekRate
					rate2.CourseMain = v
					rate2.Week = value.(*RateModel).Week
					courseRate = append(courseRate, rate2)
					judge[int(v)] = true
				}
			} else {
				//课程还未开始
				okcnt++
			}
		} //获取出勤率综合
		majorrate := rate / cnt
		rate1.AttendRate = majorrate
		majorRate = append(majorRate, rate1)
		// fmt.Println("rate:", rate)
		// fmt.Println("cnt:", cnt)
		// fmt.Println("majorrate:", majorrate)
		//rangeCoursemain 获取att的值 然后求平均值
	}
	// fmt.Println("------测试4-------------")
	return &attendinfo.CounsellorAttInfoRsp{
		Status:         200,
		Message:        "successful",
		MajorRateList:  majorRate,
		CourseRateList: courseRate,
		Total1:         uint32(len(majorRate)),
		Total2:         uint32(len(courseRate)),
	}, nil
}
