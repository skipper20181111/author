package myuser

import (
	"author/cachemodel"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"

	"author/internal/svc"
	"author/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnboundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnboundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnboundLogic {
	return &UnboundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnboundLogic) Unbound(req *types.UpdateUserInfoRes) (resp *types.UnboundResp, err error) {
	if l.ctx.Value("phone") != req.Phone {
		return &types.UnboundResp{
			Code: "4004",
			Msg:  "请勿使用其他用户的token",
		}, nil
	}
	var Userinfos cachemodel.Userinfos
	Userinfos = info2infos(req)
	err = l.svcCtx.UserModel.UpdateByPhone(l.ctx, req.Phone, &Userinfos)
	if err != nil {
		print(err)
		return &types.UnboundResp{Code: "4004", Msg: "修改失败"}, nil
	}
	infos, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, req.Phone)
	info := respons(*infos)
	jwtToken, accessExpire, refreshAfter, _ := l.getToken(info.Openid, info.Phone)
	UserBehaviour(l.ctx, l.svcCtx, "绑定新的微信号", req.Phone)
	return &types.UnboundResp{Code: "10000", Msg: "修改成功", Data: types.UnboundRp{Userinfo: info, AccessToken: jwtToken, AccessExpire: strconv.Itoa(int(accessExpire)), RefreshAfter: strconv.Itoa(int(refreshAfter))}}, nil

}
func (l *UnboundLogic) getToken(openid, phone string) (jwtToken string, accessExpire int64, refreshAfter int64, err error) {
	// ---start---
	now := time.Now().Unix()
	fmt.Println(now)
	accessExpire = l.svcCtx.Config.Auth.AccessExpire
	refreshAfter = now + l.svcCtx.Config.Auth.AccessExpire/2
	jwtToken, err = l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, openid, phone)
	if err != nil {
		return "", 0, 0, err
	}
	fmt.Println(accessExpire, refreshAfter)
	return jwtToken, accessExpire, refreshAfter, nil
	// ---end---
}
func (l *UnboundLogic) getJwtToken(secretKey string, iat, seconds int64, openid, phone string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["openid"] = openid
	claims["phone"] = phone
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
