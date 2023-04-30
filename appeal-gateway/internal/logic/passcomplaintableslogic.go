package logic

import (
	"context"

	"appeal-gateway/internal/svc"
	"appeal-gateway/internal/types"
	"appeal/appeal"

	"github.com/zeromicro/go-zero/core/logx"
)

type PassComplainTablesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPassComplainTablesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PassComplainTablesLogic {
	return &PassComplainTablesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PassComplainTablesLogic) PassComplainTables(req *types.ComplainPassReqest) (resp *types.AppealResponse, err error) {
	// todo: add your logic here and delete this line
	resq := &appeal.ComplainPassRequest{
		Cid:          req.Cid,
		ConsellorID:  req.CounsellorID, //辅导员id应该为当前用户id
		SupervisorID: req.SupervisorID,
		SchoolName:   req.University,
		Pass:         false,
	}
	if req.Pass == 1 {
		resq.Pass = true
	}
	rsp, err := l.svcCtx.Appealer.PassComplainTables(l.ctx, resq)
	return &types.AppealResponse{
		Status:  rsp.Status,
		Message: rsp.Message,
		Error:   rsp.Error,
	}, err
}
