package logic

import (
	"context"

	"appeal-gateway/rpc/appeal/appeal"
	"appeal-gateway/rpc/appeal/internal/svc"
	"appeal-gateway/rpc/appeal/model"

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
	err := l.svcCtx.MysqlDB.Model(&model.LeaveTable{}).Order("created_at desc").
		Where("student_id=? AND tag_as=1", in.GetStudentID()).Find(&lea_list).Error
	if err != nil {
		return &appeal.AppealListReply{
			Status:  30021,
			Message: "拉取失败",
			Error:   err.Error(),
		}, err
	}

	return &appeal.AppealListReply{
		Status:     200,
		AppealList: lea_list,
		Total:      uint32(len(lea_list)),
		Message:    "successful",
	}, nil
}
