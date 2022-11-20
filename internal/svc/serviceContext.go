package svc

import (
	"author/cachemodel"
	"author/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel cachemodel.UserinfosModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: cachemodel.NewUserinfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
