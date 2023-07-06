package myuser

import (
	"author/cachemodel"
	"author/internal/svc"
	"author/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

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
type TokenUtilLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTokenUtilLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenUtilLogic {
	return &TokenUtilLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *TokenUtilLogic) getJwtToken(secretKey string, iat, seconds int64, openid, phone string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["openid"] = openid
	claims["phone"] = phone
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
func (l *TokenUtilLogic) getToken(openid, phone string) (jwtToken string, accessExpire int64, refreshAfter int64, err error) {
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
func (l *TokenUtilLogic) code2Session(code string) (code2sess, error) {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
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
func (l *TokenUtilLogic) newuser(openid, phone string) error {
	_, err := l.svcCtx.UserModel.Insert(l.ctx, &cachemodel.Userinfos{Openid: openid, Phone: phone, NickName: "蟹谢", Avatar: "https://img.waterflowfit.top/avatar/picture/crabavatar.jpg", Gender: 0, Birthday: "2008-08-07", Region: "静安区"})
	return err
}
func UserBehaviour(ctx context.Context, svcCtx *svc.ServiceContext, behaviour, phone string) error {
	nowtime := time.Now()
	_, err := svcCtx.UserBehaviourModel.Insert(ctx, &cachemodel.UserBehaviourLog{Behaviour: behaviour, Phone: phone, Date: nowtime})
	return err
}

func db2info(Userinfos *cachemodel.Userinfos) (userinfo *types.UserInfo) {
	userinfo = &types.UserInfo{}
	userinfo.Gender = Userinfos.Gender
	userinfo.NickName = Userinfos.NickName
	userinfo.Avatar = Userinfos.Avatar
	userinfo.Birthday = Userinfos.Birthday
	userinfo.Region = Userinfos.Region
	userinfo.Phone = Userinfos.Phone
	userinfo.Openid = Userinfos.Openid
	return userinfo
}
func info2db(olddb *cachemodel.Userinfos, userinfo *types.UpdateUserInfoRes) *cachemodel.Userinfos {
	olddb.Gender = userinfo.Gender
	olddb.NickName = userinfo.NickName
	olddb.Avatar = userinfo.Avatar
	olddb.Birthday = userinfo.Birthday
	olddb.Region = userinfo.Region
	return olddb
}
