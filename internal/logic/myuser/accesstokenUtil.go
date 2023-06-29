package myuser

import (
	"author/cachemodel"
	"author/internal/svc"
	"author/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io/ioutil"
	"net/http"
	"time"
)

type AccessTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}
type Acyes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}
type GetPhone struct {
	Code string `json:"code"`
}

func NewAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessTokenLogic {
	return &AccessTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *AccessTokenLogic) GetPhone(code string) (bool, *types.PhoneStruct) {
	var getphone types.PhoneStruct
	token := ""
	ok := false
	one, _ := l.svcCtx.AccessTokenModel.FindOne(l.ctx, svc.Tockenid)
	if one == nil {
		if ok, token = l.ReturnToken(); !ok {
			return false, &getphone
		}
	} else {
		token = one.Token
	}

	urlPath := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", token)
	PhoneResp, err := httpc.Do(context.Background(), http.MethodPost, urlPath, GetPhone{Code: code})
	if err != nil || PhoneResp == nil {
		return false, &getphone
	}
	body, err := ioutil.ReadAll(PhoneResp.Body)
	json.Unmarshal(body, &getphone)
	PhoneResp.Body.Close()
	if getphone.Errcode != 0 {
		ok, token = l.ReturnToken()
		urlPath = fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", token)
		PhoneResp, err = httpc.Do(context.Background(), http.MethodPost, urlPath, GetPhone{Code: code})
		if err != nil || PhoneResp == nil {
			return false, &getphone
		}
		body, _ = ioutil.ReadAll(PhoneResp.Body)
		json.Unmarshal(body, &getphone)
		PhoneResp.Body.Close()
	}
	return true, &getphone
}
func (l *AccessTokenLogic) ReturnToken() (bool, string) {
	for i := 0; i < 3; i++ {
		if ok, token := l.GetToken(); ok {
			return true, token
		}
	}
	return false, ""
}

func (l *AccessTokenLogic) UpdateToken(Force bool) (bool, string) {
	if ok, token := l.GetToken(); Force && ok {
		return true, token
	}
	one, _ := l.svcCtx.AccessTokenModel.FindOne(l.ctx, svc.Tockenid)
	past := -one.Time.Sub(time.Now()).Seconds()
	if past > 600 {
		ok, token := l.GetToken()
		return ok, token
	}
	return false, ""
}
func (l *AccessTokenLogic) GetToken() (bool, string) {
	urlPath := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", l.svcCtx.Config.WxConf.AppId, l.svcCtx.Config.WxConf.Secret)
	TokenResp, err := httpc.Do(context.Background(), http.MethodGet, urlPath, nil)
	if err != nil {
		return false, ""
	}
	var ACyes Acyes
	body, _ := ioutil.ReadAll(TokenResp.Body)
	json.Unmarshal(body, &ACyes)
	TokenResp.Body.Close()
	if ACyes.AccessToken != "" && ACyes.Errcode == 0 {
		l.svcCtx.AccessTokenModel.Update(l.ctx, &cachemodel.AccessToken{Id: svc.Tockenid, Token: ACyes.AccessToken, Time: time.Now(), Overtime: int64(ACyes.ExpiresIn)})
		return true, ACyes.AccessToken
	}
	return false, ""
}
