package myuser

import (
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
	phone, _ := l.svcCtx.UserModel.FindOneByPhone(l.ctx, UserPhone)
	if phone != nil {
		phone = info2db(phone, req)
		l.svcCtx.UserModel.Update(l.ctx, phone)
	} else {
		return &types.UpdateUserInfoResp{Code: "4004", Msg: "修改失败"}, nil
	}
	infos, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, UserPhone)
	if infos == nil {
		return &types.UpdateUserInfoResp{Code: "4004", Msg: "修改失败"}, nil
	}
	info := db2info(infos)
	UserBehaviour(l.ctx, l.svcCtx, "更新用户信息", UserPhone)
	return &types.UpdateUserInfoResp{Code: "10000", Msg: "修改成功", Data: &types.UpdateUserInfoRp{Userinfo: info}}, nil
}
