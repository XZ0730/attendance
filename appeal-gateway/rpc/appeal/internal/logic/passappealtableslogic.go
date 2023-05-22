package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"appeal/appeal"
	"appeal/common/errorx"
	"appeal/internal/svc"
	"appeal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type PassAppealTablesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPassAppealTablesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PassAppealTablesLogic {
	return &PassAppealTablesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PassAppealTablesLogic) PassAppealTables(in *appeal.AppealPassRequest) (*appeal.AppealResponse, error) {
	// todo: add your logic here and delete this line

	//审核请假条  申诉
	//成功-->更新假条状态
	//请求参数:id pass
	//请假条--->查看假条周是否点名-->没点名就直接改状态-->
	//点名了和申诉表一起处理-->将redis中的考勤信息进行映射出来
	//更改当前假条或者申诉表的学号的考勤信息，更改mysql中的缺勤信息
	// l.Logger.Info()
	lt := &model.LeaveTable{}
	err2 := l.svcCtx.MysqlDB.Table("leave_table").Where("id=? AND is_audit!=3", in.GetAid()).First(&lt).Error
	if err2 != nil {
		return &appeal.AppealResponse{
			Status:  errorx.AppealIsNoExist,
			Message: errorx.GetERROR(errorx.AppealIsNoExist),
			Error:   err2.Error(),
		}, nil
	}
	course := &model.Course{}
	err1 := l.svcCtx.MysqlDB.Table("course").Where("course_id=? AND university=?", lt.CourseID, lt.SchoolName).First(&course).Error
	if err1 != nil || course.Id != in.GetCourseMain() {
		return &appeal.AppealResponse{
			Status:  errorx.CourseMatchERROR,
			Message: errorx.GetERROR(errorx.CourseMatchERROR),
			Error:   "",
		}, nil
	}

	lastweek, err := l.svcCtx.RDB.HGet(l.ctx, lt.SchoolName, strconv.Itoa(int(in.GetCourseMain()))).Result() //最近结束周
	if err != nil {
		return &appeal.AppealResponse{
			Status:  errorx.RedisInitError,
			Message: errorx.GetERROR(errorx.RedisInitError),
			Error:   err.Error(),
		}, nil
	}
	lweek, _ := strconv.Atoi(lastweek)
	attendMap := make(map[string]*model.Result, 0)
	// if lt.TagAs == 1 { //请假条
	for i := lt.LeaveCourseFrom; i <= lt.LeaveCourseTo; i++ {
		if i > lweek {
			//直接更改假条状态然后退出
			err2 := l.svcCtx.MysqlDB.Table("leave_table").Where("id=?", in.GetAid()).UpdateColumn("is_audit", 2).Error
			if err2 != nil {
				return &appeal.AppealResponse{
					Status:  38989,
					Message: "假条审核失败",
					Error:   err2.Error(),
				}, nil
			}
			break
		}
		status, err := l.svcCtx.RDB7.HGet(l.ctx, strconv.Itoa(int(in.GetCourseMain())), strconv.Itoa(int(i))).Result()
		//课程该周未点名
		if err != nil || status == "0" {
			//没点名就直接下一周
			continue
		}
		if in.GetPass() {
			fmt.Println("测试")
			info, err2 := l.svcCtx.RDB6.HGet(l.ctx, strconv.Itoa(int(in.GetCourseMain())), strconv.Itoa(int(i))).Result()
			if err2 != nil {
				fmt.Println("err2:", err2.Error())
				continue
			}
			err3 := json.Unmarshal([]byte(info), &attendMap)
			if err3 != nil {
				fmt.Println("err3:", err3)
			}
			fmt.Println("attendMap:", attendMap)
			_, ok := attendMap[lt.StudentID]
			if !ok {
				_ = l.svcCtx.MysqlDB.Table("leave_table").Where("id=?", in.GetAid()).UpdateColumn("is_audit", 3).Error
				fmt.Println("err:该学生不属于该课程")
				return &appeal.AppealResponse{
					Status:  38999,
					Message: "该学生不属于该课程,假条自动回绝",
				}, nil
			}
			if attendMap[lt.StudentID].MissAttend == 2 {
				fmt.Println("测试3")
				continue
			} else {
				//修改mysql中的值
				fmt.Println("测试4")
				att := &model.AttendTable{}
				fmt.Println("测试5")
				tx := l.svcCtx.MysqlDB.Begin()
				fmt.Println("测试1")
				err := tx.Table("attend_table").
					Where("university=? AND course_id=? AND week=?", lt.SchoolName, lt.CourseID, i).First(&att).Error
				if err != nil {
					tx.Rollback()
					return &appeal.AppealResponse{
						Status:  38989,
						Message: "假条审核失败",
						Error:   err2.Error(),
					}, nil
				}
				fmt.Println("测试2")
				unpresenters := strings.Split(att.Unpresenter, ",")
				unpresenterID := strings.Split(att.UnpresenterID, ",")
				for j := 0; j < len(unpresenterID); j++ {
					if unpresenterID[j] == lt.StudentID {
						unpresenterID = append(unpresenterID[:j], unpresenterID[j+1:]...)
						unpresenters = append(unpresenters[:j], unpresenters[j+1:]...)
						break
					}
				}
				id := ""
				name := ""
				for j := 0; j < len(unpresenterID); j++ {
					if j == 0 {
						id += unpresenterID[j]
						name += unpresenters[j]
					} else {
						id = id + "," + unpresenterID[j]
						name = name + "," + unpresenters[j]
					}
				}
				err = tx.Table("attend_table").
					Where("university=? AND course_id=? AND week=?", lt.SchoolName, lt.CourseID, i).
					UpdateColumns(model.AttendTable{Unpresenter: name, UnpresenterID: id, Unpresent: att.Unpresent - 1}).Error
				if err != nil {
					tx.Rollback()
					return &appeal.AppealResponse{
						Status:  38989,
						Message: "假条审核失败",
						Error:   err2.Error(),
					}, nil
				}
				tx.Commit()
				//
				//
				attendMap[lt.StudentID].MissAttend = 2
				attInfo, _ := json.Marshal(attendMap)
				l.svcCtx.RDB6.HSet(l.ctx, strconv.Itoa(int(in.GetCourseMain())), strconv.Itoa(int(i)), string(attInfo))
			}
		} else {
			//假条未通过----更新状态
			err2 := l.svcCtx.MysqlDB.Table("leave_table").Where("id=?", in.GetAid()).UpdateColumn("is_audit", 3).Error
			if err2 != nil {
				return &appeal.AppealResponse{
					Status:  38989,
					Message: "假条审核失败",
					Error:   err2.Error(),
				}, nil
			}
			return &appeal.AppealResponse{
				Status:  errorx.SUCCESS,
				Message: errorx.GetERROR(errorx.SUCCESS),
			}, nil
		}

	}
	err2 = l.svcCtx.MysqlDB.Table("leave_table").Where("id=?", in.GetAid()).UpdateColumn("is_audit", 2).Error
	if err2 != nil {
		return &appeal.AppealResponse{
			Status:  38989,
			Message: "假条审核失败",
			Error:   err2.Error(),
		}, nil
	}
	return &appeal.AppealResponse{
		Status:  errorx.SUCCESS,
		Message: errorx.GetERROR(errorx.SUCCESS),
	}, nil
}
