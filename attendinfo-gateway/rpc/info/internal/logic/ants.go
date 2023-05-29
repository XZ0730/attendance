package logic

import (
	"fmt"
	"info/attendinfo"
	"strconv"
	"sync"

	"github.com/panjf2000/ants"
)

type CounsellorModelReq struct {
	Req        *GetCounsellorAttInfoLogic
	Req1       *GetCourseAttInfoLogic
	RateMap    *sync.Map
	CourseMain uint64
	University string
}
type RateModel struct {
	Week         int32
	LastWeekRate float64
	AverageRate  float64
}

var (
	ants_Pool *ants.PoolWithFunc
	wait      sync.WaitGroup
)

func Init() {

	ants_Pool, _ = ants.NewPoolWithFunc(500, func(payload interface{}) {
		defer wait.Done()
		cmr := payload.(*CounsellorModelReq)
		//查询最近结束周的缺勤人数
		week, err := cmr.Req.svcCtx.RDB.HGet(cmr.Req.ctx, cmr.University, strconv.Itoa(int(cmr.CourseMain))).Result()
		if err != nil {
			return
		}
		//得到最近结束周出勤率
		//得到截止目前的课程平均出勤率
		lastweek, _ := strconv.Atoi(week)
		rsp, _ := cmr.Req1.GetCourseAttInfo(&attendinfo.CourseInfoRequest{
			University:  cmr.University,
			Course_Main: uint32(cmr.CourseMain),
		})
		rateSum := 0.0
		lastweekRate := 0.0
		cnt := 0.0
		if rsp.Status == 200 {
			for _, v := range rsp.CourseAttInfos {
				if lastweek == int(v.Week) {
					lastweekRate = v.AttendRate
				}
				rateSum += v.AttendRate //获取和
				cnt++
			}
		} else {
			return
		}
		if len(rsp.CourseAttInfos) == 0 {
			return
		}
		fmt.Println("--------------------:", rateSum, "---cnt:", cnt, "---:", cmr.CourseMain)
		rate := &RateModel{
			LastWeekRate: lastweekRate,
			AverageRate:  rateSum / cnt,
			Week:         int32(lastweek),
		}

		cmr.RateMap.Store(cmr.CourseMain, rate)
	})

}
