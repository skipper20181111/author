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
	if l.ctx.Value("openid") != req.Openid || l.ctx.Value("phone") != req.Phone {
		return &types.GetUserInfoResp{
			Code: "4004",
			Msg:  "请勿使用其他用户的token",
		}, nil
	}
	newinfos, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, req.Phone)
	if err != nil {
		return &types.GetUserInfoResp{Code: "4004", Msg: "未查询到用户信息"}, nil
	}
	userinfo := respons(*newinfos)

	return &types.GetUserInfoResp{
		Code: "10000",
		Msg:  "成功查询用户信息",
		Data: types.GetUserInfoRp{Userinfo: userinfo},
	}, nil
}
