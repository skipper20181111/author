package myuser

import (
	"author/cachemodel"
	"context"

	"author/internal/svc"
	"author/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateinfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateinfoLogic {
	return &UpdateinfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateinfoLogic) Updateinfo(req *types.UpdateUserInfoRes) (resp *types.UpdateUserInfoResp, err error) {
	if l.ctx.Value("openid") != req.Openid {
		return &types.UpdateUserInfoResp{
			Code: "4004",
			Msg:  "请勿使用其他用户的token",
		}, nil
	}
	var Userinfos cachemodel.Userinfos
	Userinfos = info2infos(req)
	err = l.svcCtx.UserModel.UpdateByOpenid(l.ctx, &Userinfos)
	if err != nil {
		print(err)
		return &types.UpdateUserInfoResp{Code: "4004", Msg: "修改失败"}, nil
	}
	infos, err := l.svcCtx.UserModel.FindOneByOpenid(l.ctx, req.Openid)
	info := respons(*infos)

	return &types.UpdateUserInfoResp{Code: "10000", Msg: "修改成功", Data: types.UpdateUserInfoRp{Userinfo: info}}, nil
}
func info2infos(userinfo *types.UpdateUserInfoRes) (Userinfos cachemodel.Userinfos) {
	Userinfos.Gender = int64(userinfo.Gender)
	Userinfos.NickName = userinfo.NickName
	Userinfos.Avatar = userinfo.Avatar
	Userinfos.Birthday = userinfo.Birthday
	Userinfos.Region = userinfo.Region
	Userinfos.Phone = userinfo.Phone
	Userinfos.Openid = userinfo.Openid
	Userinfos.Wechatbind = int64(userinfo.Wechatbind)
	Userinfos.Effective = int64(userinfo.Effective)
	return Userinfos
}
