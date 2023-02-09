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
	UserPhone := l.ctx.Value("phone").(string)
	var Userinfos cachemodel.Userinfos
	Userinfos = l.info2infos(req)
	err = l.svcCtx.UserModel.UpdateByPhone(l.ctx, UserPhone, &Userinfos)
	if err != nil {
		print(err)
		return &types.UpdateUserInfoResp{Code: "4004", Msg: "修改失败"}, nil
	}
	infos, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, UserPhone)
	if infos == nil {
		return &types.UpdateUserInfoResp{Code: "4004", Msg: "修改失败"}, nil
	}
	info := respons(*infos)
	UserBehaviour(l.ctx, l.svcCtx, "更新用户信息", UserPhone)
	return &types.UpdateUserInfoResp{Code: "10000", Msg: "修改成功", Data: types.UpdateUserInfoRp{Userinfo: info}}, nil
}
func (l *UpdateinfoLogic) info2infos(userinfo *types.UpdateUserInfoRes) (Userinfos cachemodel.Userinfos) {
	Userinfos.Gender = userinfo.Gender
	Userinfos.NickName = userinfo.NickName
	Userinfos.Avatar = userinfo.Avatar
	Userinfos.Birthday = userinfo.Birthday
	Userinfos.Region = userinfo.Region
	Userinfos.Phone = l.ctx.Value("phone").(string)
	Userinfos.Openid = l.ctx.Value("openid").(string)
	return Userinfos
}
