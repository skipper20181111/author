package refresh

import (
	"author/internal/logic/myuser"
	"author/internal/svc"
	"author/internal/types"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	atl    *myuser.AccessTokenLogic
}

func NewRefreshatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshatLogic {
	return &RefreshatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		atl:    myuser.NewAccessTokenLogic(ctx, svcCtx),
	}
}
func (l *RefreshatLogic) Refreshat(req *types.RefreshRes) (resp *types.RefreshResp, err error) {
	ok, _ := l.atl.UpdateToken(req.Force)
	if ok {
		return &types.RefreshResp{Code: "10000", Msg: "刷新成功"}, nil
	} else {
		return &types.RefreshResp{Code: "4004", Msg: "失败"}, nil
	}
}
