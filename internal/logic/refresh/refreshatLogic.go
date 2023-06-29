package refresh

import (
	"author/cachemodel"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"author/internal/svc"
	"author/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const Tockenid = 1

type RefreshatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshatLogic {
	return &RefreshatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshatLogic) Refreshat(req *types.RefreshRes) (resp *types.RefreshResp, err error) {

	if req.Force == true {
		refresh, err := l.refresh()
		if err != nil {
			return &types.RefreshResp{Code: "4004", Msg: refresh.Errmsg}, nil
		}
		l.svcCtx.AccessTokenModel.Update(l.ctx, &cachemodel.AccessToken{Id: Tockenid, Token: refresh.AccessToken, Time: time.Now()})
		return &types.RefreshResp{Code: "10000", Msg: "刷新成功"}, nil
	}
	one, err := l.svcCtx.AccessTokenModel.FindOne(l.ctx, Tockenid)
	past := -one.Time.Sub(time.Now()).Seconds()
	if (past*past)/25000000.1 > rand.Float64() {
		refresh, err := l.refresh()
		if err != nil {
			return &types.RefreshResp{Code: "4004", Msg: refresh.Errmsg}, nil
		}
		l.svcCtx.AccessTokenModel.Update(l.ctx, &cachemodel.AccessToken{Id: Tockenid, Token: refresh.AccessToken, Time: time.Now()})
		return &types.RefreshResp{Code: "10000", Msg: "刷新成功"}, nil
	}
	return &types.RefreshResp{Code: "10000", Msg: "下次再说"}, nil

}
func (l *RefreshatLogic) refresh() (*Acyes, error) {
	urlPath := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", l.svcCtx.Config.WxConf.AppId, l.svcCtx.Config.WxConf.Secret)
	TokenResp, err := httpc.Do(context.Background(), http.MethodGet, urlPath, nil)
	var ACyes Acyes
	body, _ := ioutil.ReadAll(TokenResp.Body)
	json.Unmarshal(body, &ACyes)

	TokenResp.Body.Close()
	return &ACyes, err

}

type Acyes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}
