package logic

import (
	"context"

	"appeal/appeal"
	"appeal/internal/svc"
	"appeal/model"

	"appeal/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppealListBySidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppealListBySidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppealListBySidLogic {
	return &GetAppealListBySidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppealListBySidLogic) GetAppealListBySid(in *appeal.AppealListRequset) (*appeal.AppealListReply, error) {
	// todo: add your logic here and delete this line
	lea_list := make([]*appeal.AppealModel, 0)
	if in.GetTag() == 1 {
		err := l.svcCtx.MysqlDB.Model(&model.LeaveTable{}).Order("created_at desc").
			Where("student_id=?", in.GetStudentID()).Find(&lea_list).Error
		if err != nil {
			return &appeal.AppealListReply{
				Status:  errorx.FailTOPull,
				Message: errorx.GetERROR(errorx.FailTOPull),
				Error:   err.Error(),
			}, nil
		}
	} else {
		err := l.svcCtx.MysqlDB.Model(&model.LeaveTable{}).Order("created_at desc").
			Where("counsellor_id=?", in.GetCounsellorID()).Find(&lea_list).Error
		if err != nil {
			return &appeal.AppealListReply{
				Status:  errorx.FailTOPull,
				Message: errorx.GetERROR(errorx.FailTOPull),
				Error:   err.Error(),
			}, nil
		}
	}

	return &appeal.AppealListReply{
		Status:     errorx.SUCCESS,
		AppealList: lea_list,
		Total:      uint32(len(lea_list)),
		Message:    errorx.GetERROR(errorx.SUCCESS),
	}, nil
}
