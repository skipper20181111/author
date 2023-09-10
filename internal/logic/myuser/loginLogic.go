package myuser

import (
	"author/internal/svc"
	"author/internal/types"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	tul    *TokenUtilLogic
}

// code2sess结构体用于parse小程序login接口返回的参数，code2Session函数中使用了

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		tul:    NewTokenUtilLogic(ctx, svcCtx),
	}
}

func (l *LoginLogic) Login(req *types.LoginRes) (resp *types.LoginResp, err error) {
	loginresp := &types.LoginResp{}
	/*
	   调用小程序code2Session接口获取openid等
	   测试状态下，logincode=openid
	*/
	wxmsg, err := l.tul.code2Session(req.LoginCode)
	//wxmsg.Openid = req.LoginCode
	if err != nil || wxmsg.Openid == "" {
		return &types.LoginResp{Code: "4004", Msg: wxmsg.Errmsg}, nil
	}
	infos, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, req.Phone)
	if infos == nil {
		/*
		   此时为新用户，需要新建用户，并指示新用户
		*/
		l.tul.newuser(wxmsg.Openid, req.Phone)
		loginresp.Msg = "返回用户信息"
		loginresp.Code = "10000"
		newinfos, _ := l.svcCtx.UserModel.FindOneByPhone(l.ctx, req.Phone)
		if newinfos == nil {
			loginresp.Msg = "数据库断联"
			loginresp.Code = "4004"
			return loginresp, nil
		}
		userinfo := db2info(newinfos)
		jwtToken, _, _, _ := l.tul.getToken(wxmsg.Openid, req.Phone)
		loginresp.Data = &types.LoginRp{Userinfo: userinfo, IsRebund: false, LoginSuccess: true, IsNew: true, AccessToken: jwtToken}
		UserBehaviour(l.ctx, l.svcCtx, "新建用户", userinfo.Phone)
		return loginresp, nil
	} else {
		/*
			此时确定为老用户，但是应当告知是否解绑
		*/
		if wxmsg.Openid != infos.Openid {
			l.svcCtx.UserModel.UnBoundOpenId(l.ctx, wxmsg.Openid)
			infos.Openid = wxmsg.Openid
			l.svcCtx.UserModel.Update(l.ctx, infos)
			jwtToken, _, _, _ := l.tul.getToken(wxmsg.Openid, req.Phone)
			return &types.LoginResp{Code: "10000", Msg: "解除绑定", Data: &types.LoginRp{Userinfo: db2info(infos), AccessToken: jwtToken, LoginSuccess: true, IsNew: false, IsRebund: true}}, nil
		} else {
			userinfo := db2info(infos)
			loginresp.Msg = "登录成功，返回用户信息"
			loginresp.Code = "10000"
			jwtToken, _, _, _ := l.tul.getToken(wxmsg.Openid, req.Phone)
			loginresp.Data = &types.LoginRp{Userinfo: userinfo, IsNew: false, AccessToken: jwtToken}
			UserBehaviour(l.ctx, l.svcCtx, "用户登录", userinfo.Phone)
			return loginresp, nil
		}
	}
}
