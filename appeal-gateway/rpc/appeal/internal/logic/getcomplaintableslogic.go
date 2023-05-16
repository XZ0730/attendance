package logic

import (
	"context"

	"appeal/appeal"
	"appeal/common/errorx"
	"appeal/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetComplainTablesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetComplainTablesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetComplainTablesLogic {
	return &GetComplainTablesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetComplainTablesLogic) GetComplainTables(in *appeal.ComplainGetRequest) (*appeal.ComplainResponse, error) {
	// todo: 辅导员获取投诉表
	//获取特定辅导员管理的投诉表，模糊查询：学院 督导员名字 督导员学号 督导员专业 投诉学生名字
	cts := make([]*appeal.ComplainModel, 0)
	err := l.svcCtx.MysqlDB.Table("complain_table").
		Where("counsellor_id=? AND college LIKE ? AND supervisor_major LIKE ? AND student_name LIKE ? AND supervisor_name LIKE ? OR supervisor_id=? AND deleted_at IS NULL",
			in.GetConsellorID(), in.GetCollege(), in.GetMajor(), in.GetStudentName(), in.GetSupervisorName(), in.GetSupervisorId()).
		Find(&cts).Error
	if err != nil {
		return &appeal.ComplainResponse{
			Status:  errorx.RecordGetError,
			Message: errorx.GetERROR(errorx.RecordGetError),
			Error:   err.Error(),
		}, err
	}
	return &appeal.ComplainResponse{
		Status:       errorx.SUCCESS,
		ComplainList: cts,
		Total:        uint32(len(cts)),
		Message:      errorx.GetERROR(errorx.SUCCESS),
	}, nil
}
