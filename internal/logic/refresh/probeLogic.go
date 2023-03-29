package refresh

import (
	"context"

	"author/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProbeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProbeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProbeLogic {
	return &ProbeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProbeLogic) Probe() error {

	return nil
}
