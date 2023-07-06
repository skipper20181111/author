package myuser

import (
	"context"

	"author/internal/svc"
	"author/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NophoneloginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	tul    *TokenUtilLogic
}

func NewNophoneloginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NophoneloginLogic {
	return &NophoneloginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		tul:    NewTokenUtilLogic(ctx, svcCtx),
	}
}

func (l *NophoneloginLogic) Nophonelogin(req *types.NoPhoneLoginRes) (resp *types.LoginResp, err error) {
	loginresp := &types.LoginResp{Code: "10000"}
	wxmsg, err := l.tul.code2Session(req.LoginCode)
	//wxmsg.Openid = req.LoginCode
	if err != nil || wxmsg.Openid == "" {
		return &types.LoginResp{Code: "4004", Msg: wxmsg.Errmsg}, nil
	}
	userinfos, _ := l.svcCtx.UserModel.FindOneByOpenid(l.ctx, wxmsg.Openid)
	if userinfos == nil {
		loginresp.Msg = "failed"
		loginresp.Data = &types.LoginRp{LoginSuccess: false, IsRebund: false, IsNew: false}
	} else {
		loginresp.Msg = "success"
		jwtToken, _, _, _ := l.tul.getToken(wxmsg.Openid, userinfos.Phone)
		loginresp.Data = &types.LoginRp{AccessToken: jwtToken, LoginSuccess: true, IsRebund: false, IsNew: false, Userinfo: db2info(userinfos)}
	}
	return loginresp, nil
}
