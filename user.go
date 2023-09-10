package main

import (
	"author/cachemodel"
	"author/internal/config"
	"author/internal/handler"
	"author/internal/svc"
	"author/internal/types"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io/ioutil"
	"net/http"
	"time"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	go refresscache()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
func refresscache() {
	for true {
		fmt.Println("开始刷新")
		time.Sleep(time.Second * 3)
		urlPath := "http://localhost:8888/refresh/refreshat"
		data := ForceRefresh{Force: false}
		resp, _ := httpc.Do(context.Background(), http.MethodPost, urlPath, data)
		if resp == nil {
			time.Sleep(time.Second * 50)
			continue
		} else {
			fmt.Println("结束刷新", resp)
			fmt.Println(resp.Body.Close())
			time.Sleep(time.Second * 50)
		}
	}
}

type ForceRefresh struct {
	Force bool `json:"force"`
}

func wxnmsl(svcCtx *svc.ServiceContext) {
	for true {
		time.Sleep(time.Second * 3)
		wxDeliveries, _ := svcCtx.WxDelivery.FindAll(context.Background())
		for _, wxDelivery := range wxDeliveries {
			giveMHTshit(svcCtx, wxDelivery)
			ConfirmMHTshit(svcCtx, wxDelivery)
		}
		time.Sleep(time.Second * 50)
	}
}
func ConfirmMHTshit(svcCtx *svc.ServiceContext, Payinfo *cachemodel.WxDelivery) {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	ctx := context.Background()
	accessToken, _ := svcCtx.AccessTokenModel.FindOne(ctx, 1)
	UrlPath := fmt.Sprintf("https://api.weixin.qq.com/wxa/sec/order/get_order?access_token=%s", accessToken.Token)
	resp, _ := httpc.Do(context.Background(), http.MethodPost, UrlPath, types.MsgDelivering{TransactionId: Payinfo.TransactionId})
	fmt.Println(resp.Body.Close())
	res := types.MsgReturn{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &res)
	if len(res.Order.Openid) > 1 && res.Order.OrderState == 2 {
		svcCtx.WxDelivery.UpdateDelivering(ctx, Payinfo.OutTradeNo)
	}
}
func giveMHTshit(svcCtx *svc.ServiceContext, Payinfo *cachemodel.WxDelivery) {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	ctx := context.Background()
	orders, _ := svcCtx.Order.FindAllByOutTradeNo(ctx, Payinfo.OutTradeNo)
	userinfos, _ := svcCtx.UserModel.FindOneByPhone(ctx, Payinfo.Phone)

	shippinginfo := make([]*types.ShippingList, 0)
	for _, order := range orders {
		shippinginfo = append(shippinginfo, &types.ShippingList{
			TrackingNo:     order.DeliverySn,
			ExpressCompany: "SF",
			ItemDesc:       order.ProductInfo,
			Contact: &types.Contact{
				ConsignorContact: "178****0845",
			},
		})
	}
	DataMsg := types.MsgData{
		OrderKey: &types.OrderKey{
			OrderNumberType: 2,
			TransactionId:   Payinfo.TransactionId,
			//Mchid:           "1652716843",
			//OutTradeNo:      "lt9i1DBXoZiefJD5pfPyS8zMV1P2i7GL",
		},
		LogisticsType:  1,
		DeliveryMode:   2,
		IsAllDelivered: true,
		ShippingList:   shippinginfo,
		UploadTime:     time.Now().Format("2006-01-02T15:04:05.000+08:00"),
		Payer: &types.Payer{
			Openid: userinfos.Openid,
		},
	}
	accessToken, _ := svcCtx.AccessTokenModel.FindOne(ctx, 1)
	UrlPath := fmt.Sprintf("https://api.weixin.qq.com/wxa/sec/order/upload_shipping_info?access_token=%s", accessToken.Token)
	resp, _ := httpc.Do(context.Background(), http.MethodPost, UrlPath, DataMsg)
	fmt.Println(resp.Body.Close())

}
