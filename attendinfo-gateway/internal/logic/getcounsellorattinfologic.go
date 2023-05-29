package logic

import (
	"context"
	"info/attendinfo"

	"attendinfo-gateway/internal/svc"
	"attendinfo-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCounsellorAttInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCounsellorAttInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCounsellorAttInfoLogic {
	return &GetCounsellorAttInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCounsellorAttInfoLogic) GetCounsellorAttInfo(req *types.CounsellorInfoRequest) (resp *types.CounsellorInfoResponse, err error) {
	// todo: add your logic here and delete this line
	// fmt.Println("测试1------------------------------")
	mrl := make([]*attendinfo.MajorReqList, 0)
	for _, v := range req.MajorInfoList {
		aimr := &attendinfo.MajorReqList{
			Major:      v.Major,
			University: v.University,
			CouseMain:  v.CourseMain,
		}
		mrl = append(mrl, aimr)
	}
	// fmt.Println("测试2------------------------------")
	rsp, _ := l.svcCtx.InfoCli.GetCounsellorAttInfo(l.ctx, &attendinfo.CounsellorAttInfoReq{
		MajorList: mrl,
	})
	// fmt.Println("rsp:", rsp)
	// fmt.Println("测试3------------------------------")
	return &types.CounsellorInfoResponse{
		Status:     rsp.Status,
		MajorData:  rsp.MajorRateList,
		Total1:     rsp.Total1,
		CourseData: rsp.CourseRateList,
		Total2:     rsp.Total2,
		Message:    rsp.Message,
	}, nil
}
