package myuser

import (
	"author/internal/logic/refresh"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io/ioutil"
	"net/http"

	"author/internal/svc"
	"author/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetphoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetphoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetphoneLogic {
	return &GetphoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetphoneLogic) Getphone(req *types.GetPhoneRes) (resp *types.GetPhoneResp, err error) {
	one, err := l.svcCtx.AccessTokenModel.FindOne(l.ctx, refresh.Tockenid)
	var getphone types.PhoneStruct
	urlPath := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", one.Token)
	PhoneResp, err := httpc.Do(context.Background(), http.MethodPost, urlPath, getphonestruct{Code: req.PhoneCode})
	body, _ := ioutil.ReadAll(PhoneResp.Body)
	json.Unmarshal(body, &getphone)
	PhoneResp.Body.Close()
	if getphone.Errcode != 0 {
		return &types.GetPhoneResp{Code: "4004", Msg: getphone.Errmsg}, nil
	}

	return &types.GetPhoneResp{Code: "10000", Msg: "success", Data: &types.GetPhoneRp{Phone: getphone.Phoneinfo.PurePhoneNumber, CountryCode: getphone.Phoneinfo.CountryCode}}, nil
}

type getphonestruct struct {
	Code string `json:"code"`
}
