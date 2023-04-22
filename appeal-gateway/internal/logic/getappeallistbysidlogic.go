package logic

import (
	"context"

	"appeal-gateway/internal/svc"
	"appeal-gateway/internal/types"
	"appeal-gateway/rpc/appeal/appeal"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppealListBySidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAppealListBySidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppealListBySidLogic {
	return &GetAppealListBySidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppealListBySidLogic) GetAppealListBySid(req *types.AppealListRequest) (resp *types.ListReply, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.Appealer.GetAppealListBySid(l.ctx, &appeal.AppealListRequset{
		StudentID: req.StudentID,
	})
	if err != nil {
		return nil, err
	}
	return &types.ListReply{
		Status:  res.Status,
		Message: res.Message,
		Error:   res.Error,
		Data:    res.AppealList,
		Total:   res.Total,
	}, nil
}
