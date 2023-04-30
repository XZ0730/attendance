package logic

import (
	"context"

	"appeal-gateway/internal/svc"
	"appeal-gateway/internal/types"
	"appeal/appeal"

	"github.com/zeromicro/go-zero/core/logx"
)

type ComplainToSupervisorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewComplainToSupervisorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ComplainToSupervisorLogic {
	return &ComplainToSupervisorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ComplainToSupervisorLogic) ComplainToSupervisor(req *types.ComplainRequest) (resp *types.ComplainResponse, err error) {
	// todo: add your logic here and delete this line
	rsp, err := l.svcCtx.Appealer.ComplainToSupervisor(l.ctx, &appeal.ComplainRequest{
		SupervisorID:   req.SupervisorID,
		SupervisorName: req.SupervisorName,
		SchoolName:     req.SchoolName,
		Reason:         req.Reason,
		CounsellorName: req.CounsellorName,
		CounsellorID:   req.CounsellorID,
		StudentID:      req.StudentID,
	})
	if err != nil {
		return &types.ComplainResponse{
			Status:  rsp.Status,
			Message: rsp.Message,
			Error:   rsp.Error,
		}, err
	}
	return &types.ComplainResponse{
		Status:  rsp.Status,
		Message: rsp.Message,
		Error:   rsp.Error,
	}, nil
}
