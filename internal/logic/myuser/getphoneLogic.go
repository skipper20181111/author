package myuser

import (
	"context"

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

func (l *GetphoneLogic) Getphone(req *types.LoginRes) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line

	return
}
