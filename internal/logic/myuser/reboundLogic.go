package myuser

import (
	"author/internal/svc"
	"author/internal/types"
	"context"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReboundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	tul    *TokenUtilLogic
}

func NewReboundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReboundLogic {
	return &ReboundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		tul:    NewTokenUtilLogic(ctx, svcCtx),
	}
}

func (l *ReboundLogic) Rebound(req *types.ReboundRes) (resp *types.ReboundResp, err error) {
	UserPhone := l.ctx.Value("phone").(string)
	wxmsg, err := l.tul.code2Session(req.LoginCode)
	if err != nil || wxmsg.Openid == "" {
		return &types.ReboundResp{Code: "4004", Msg: wxmsg.Errmsg}, nil
	}
	infos, _ := l.svcCtx.UserModel.FindOneByPhone(l.ctx, UserPhone)
	if infos != nil {
		infos.Openid = wxmsg.Openid
		l.svcCtx.UserModel.Update(l.ctx, infos)
	}
	newinfos, _ := l.svcCtx.UserModel.FindOneByPhone(l.ctx, UserPhone)
	if newinfos == nil || newinfos.Openid != wxmsg.Openid {
		return &types.ReboundResp{Code: "4004", Msg: "failed"}, nil
	}
	info := db2info(newinfos)
	jwtToken, accessExpire, refreshAfter, _ := l.tul.getToken(info.Openid, info.Phone)
	UserBehaviour(l.ctx, l.svcCtx, "绑定新的微信号", UserPhone)
	return &types.ReboundResp{Code: "10000", Msg: "修改成功", Data: &types.ReboundRp{Userinfo: info, AccessToken: jwtToken, AccessExpire: strconv.Itoa(int(accessExpire)), RefreshAfter: strconv.Itoa(int(refreshAfter))}}, nil

}
