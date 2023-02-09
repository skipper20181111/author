package myuser

import (
	"author/internal/svc"
	"author/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetinfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetinfoLogic {
	return &GetinfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetinfoLogic) Getinfo(req *types.GetUserInfoRes) (resp *types.GetUserInfoResp, err error) {
	UserPhone := l.ctx.Value("phone").(string)
	newinfos, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, UserPhone)
	if newinfos == nil {
		return &types.GetUserInfoResp{Code: "4004", Msg: "未查询到用户信息"}, nil
	}
	userinfo := respons(*newinfos)

	return &types.GetUserInfoResp{
		Code: "10000",
		Msg:  "成功查询用户信息",
		Data: types.GetUserInfoRp{Userinfo: userinfo},
	}, nil
}
