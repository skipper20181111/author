package myuser

import (
	"author/cachemodel"
	"author/internal/svc"
	"author/internal/types"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// code2sess结构体用于parse小程序login接口返回的参数，code2Session函数中使用了
type code2sess struct {
	Openid      string
	Session_key string
	Unionid     string
	Errcode     int
	Errmsg      string
}
type NullTime struct {
	Time  time.Time
	Valid bool
}

func ToNullTime(t time.Time) NullTime {
	return NullTime{Time: t, Valid: !t.IsZero()}
}
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds int64, openid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["openid"] = openid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func (l *LoginLogic) Login(req *types.LoginRes) (resp *types.LoginResp, err error) {
	var loginresp types.LoginResp
	/*
	   调用小程序code2Session接口获取openid等
	   测试状态下，logincode=openid
	*/
	wxmsg, err := l.code2Session(req.LoginCode)
	//wxmsg.Openid = req.LoginCode
	if err != nil || wxmsg.Openid == "" {
		return &types.LoginResp{Code: "4004", Msg: wxmsg.Errmsg}, nil
	}
	infos, err := l.svcCtx.UserModel.FindOneByOpenid(l.ctx, wxmsg.Openid)
	if err != nil {
		/*
		   此时为新用户，需要新建用户，并指示新用户
		*/
		l.newuser(wxmsg.Openid)
		loginresp.Msg = "返回用户信息"
		loginresp.Code = "10000"
		newinfos, _ := l.svcCtx.UserModel.FindOneByOpenid(l.ctx, wxmsg.Openid)
		userinfo := respons(*newinfos)
		jwtToken, accessExpire, refreshAfter, _ := l.getToken(wxmsg.Openid)
		loginresp.Data = types.LoginRp{Userinfo: userinfo, IsNew: 1, AccessToken: jwtToken, AccessExpire: strconv.Itoa(int(accessExpire)), RefreshAfter: strconv.Itoa(int(refreshAfter))}
		return &loginresp, nil
	}
	userinfo := respons(*infos)
	loginresp.Msg = "登录成功，返回用户信息"
	loginresp.Code = "10000"
	jwtToken, accessExpire, refreshAfter, err := l.getToken(wxmsg.Openid)
	loginresp.Data = types.LoginRp{Userinfo: userinfo, IsNew: 0, AccessToken: jwtToken, AccessExpire: strconv.Itoa(int(accessExpire)), RefreshAfter: strconv.Itoa(int(refreshAfter))}
	return &loginresp, nil
}
func (l *LoginLogic) getToken(openid string) (jwtToken string, accessExpire int64, refreshAfter int64, err error) {
	// ---start---
	now := time.Now().Unix()
	fmt.Println(now)
	accessExpire = l.svcCtx.Config.Auth.AccessExpire
	refreshAfter = now + l.svcCtx.Config.Auth.AccessExpire/2
	jwtToken, err = l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, openid)
	if err != nil {
		return "", 0, 0, err
	}
	fmt.Println(accessExpire, refreshAfter)
	return jwtToken, accessExpire, refreshAfter, nil
	// ---end---
}
func (l *LoginLogic) newuser(openid string) error {
	_, err := l.svcCtx.UserModel.Insert(l.ctx, &cachemodel.Userinfos{Openid: openid, Phone: openid, NickName: "蟹蟹", Avatar: "'https://img.waterflowfit.top/img/微信图片_20220925111848.jpg'", Gender: 0, Birthday: "2004-09-01", Region: "静安区", Wechatbind: 0, Effective: 1, CreatedAt: sql.NullTime(ToNullTime(time.Now())), UpdatedAt: sql.NullTime(ToNullTime(time.Now()))})
	return err
}
func respons(Userinfos cachemodel.Userinfos) (userinfo types.UserInfo) {
	userinfo.Gender = int(Userinfos.Gender)
	userinfo.NickName = Userinfos.NickName
	userinfo.Avatar = Userinfos.Avatar
	userinfo.Birthday = Userinfos.Birthday
	userinfo.Region = Userinfos.Region
	userinfo.Phone = Userinfos.Phone
	userinfo.Openid = Userinfos.Openid
	userinfo.Wechatbind = int(Userinfos.Wechatbind)
	userinfo.Effective = int(Userinfos.Effective)
	return userinfo
}
func (l *LoginLogic) code2Session(code string) (wxLogMsg code2sess, err error) {
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
