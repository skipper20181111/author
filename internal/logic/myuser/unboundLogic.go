package myuser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (l *UnboundLogic) Unbound(req *types.UnboundRes) (resp *types.UnboundResp, err error) {
	UserPhone := l.ctx.Value("phone").(string)
	wxmsg, err := l.code2Session(req.LoginCode)
	if err != nil || wxmsg.Openid == "" {
		return &types.UnboundResp{Code: "4004", Msg: wxmsg.Errmsg}, nil
	}
	infos, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, UserPhone)
	infos.Openid = wxmsg.Openid

	err = l.svcCtx.UserModel.UpdateByPhone(l.ctx, UserPhone, infos)
	if err != nil {
		print(err)
		return &types.UnboundResp{Code: "4004", Msg: "修改失败"}, nil
	}
	newinfos, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, UserPhone)
	if newinfos == nil {
		return &types.UnboundResp{Code: "4004", Msg: "修改失败"}, nil
	}
	info := respons(*newinfos)
	jwtToken, accessExpire, refreshAfter, _ := l.getToken(info.Openid, info.Phone)
	UserBehaviour(l.ctx, l.svcCtx, "绑定新的微信号", UserPhone)
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
func (l *UnboundLogic) code2Session(code string) (wxLogMsg code2sess, err error) {
	var res code2sess
	params := url.Values{}
	Url, err := url.Parse("https://api.weixin.qq.com/sns/jscode2session")
	if err != nil {
		return res, err
	}

	params.Set("appid", l.svcCtx.Config.WxConf.AppId)
	params.Set("secret", l.svcCtx.Config.WxConf.Secret)
	params.Set("js_code", code)
	params.Set("grant_type", l.svcCtx.Config.WxConf.Grant_type)
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath) // https://httpbin.org/get?age=23&name=zhaofan
	resp, err := http.Get(urlPath)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	json.Unmarshal(body, &res)
	fmt.Println(res)
	return res, nil
}
