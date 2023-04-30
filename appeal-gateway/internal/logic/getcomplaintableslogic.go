package logic

import (
	"context"

	"appeal-gateway/internal/svc"
	"appeal-gateway/internal/types"
	"appeal/appeal"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetComplainTablesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetComplainTablesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetComplainTablesLogic {
	return &GetComplainTablesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetComplainTablesLogic) GetComplainTables(req *types.ComplainGetRequest) (resp *types.ListReply, err error) {
	// todo: add your logic here and delete this line
	req.SupervisorName = "%" + req.SupervisorName + "%"
	req.College = "%" + req.College + "%"
	req.StudentName = "%" + req.StudentName + "%"
	req.Major = "%" + req.Major + "%"

	res, err := l.svcCtx.Appealer.GetComplainTables(l.ctx, &appeal.ComplainGetRequest{
		SupervisorId:   req.SupervisorID,
		SupervisorName: req.SupervisorName,
		College:        req.College,
		StudentName:    req.StudentName,
		ConsellorID:    req.CounsellorID,
		Major:          req.Major,
	})
	return &types.ListReply{
		Status:  res.Status,
		Data:    res.ComplainList,
		Total:   res.Total,
		Error:   res.Error,
		Message: res.Message,
	}, err
}
