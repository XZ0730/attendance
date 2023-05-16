package logic

import (
	"appeal/appeal"
	"context"

	"appeal-gateway/internal/svc"
	"appeal-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PassAppealTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPassAppealTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PassAppealTableLogic {
	return &PassAppealTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PassAppealTableLogic) PassAppealTable(req *types.AppealPassRequest) (resp *types.AppealResponse, err error) {
	// todo: add your logic here and delete this line
	var apr *appeal.AppealPassRequest
	apr = &appeal.AppealPassRequest{
		Aid:        uint64(req.Aid),
		CourseMain: req.CourseMain,
	}
	if req.Pass == 1 {
		apr.Pass = true
	} else {
		apr.Pass = false
	}
	ar, err2 := l.svcCtx.Appealer.PassAppealTables(l.ctx, apr)
	if err2 != nil {
		return &types.AppealResponse{
			Status:  ar.Status,
			Message: ar.Message,
			Error:   err2.Error(),
		}, nil
	}
	return &types.AppealResponse{
		Status:  ar.Status,
		Message: ar.Message,
		Error:   ar.Error,
	}, nil
}
