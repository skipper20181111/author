package myuser

import (
	"author/internal/svc"
	"author/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetphoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	atl    *AccessTokenLogic
}

func NewGetphoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetphoneLogic {
	return &GetphoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		atl:    NewAccessTokenLogic(ctx, svcCtx),
	}
}

func (l *GetphoneLogic) Getphone(req *types.GetPhoneRes) (resp *types.GetPhoneResp, err error) {
	_, getphone := l.atl.GetPhone(req.PhoneCode)
	if getphone.Errcode != 0 {
		return &types.GetPhoneResp{Code: "10000", Msg: getphone.Errmsg}, nil
	}
	return &types.GetPhoneResp{Code: "10000", Msg: "success", Data: &types.GetPhoneRp{Phone: getphone.Phoneinfo.PurePhoneNumber, CountryCode: getphone.Phoneinfo.CountryCode}}, nil
}
