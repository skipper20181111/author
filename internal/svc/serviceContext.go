package svc

import (
	"author/cachemodel"
	"author/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config             config.Config
	UserModel          cachemodel.UserinfosModel
	UserBehaviourModel cachemodel.UserBehaviourLogModel
	AccessTokenModel   cachemodel.AccessTokenModel
	UserPointsModel    cachemodel.UserPointsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		UserModel:          cachemodel.NewUserinfosModel(sqlx.NewMysql(c.DB.DataSource)),
		UserBehaviourModel: cachemodel.NewUserBehaviourLogModel(sqlx.NewMysql(c.DB.DataSource)),
		AccessTokenModel:   cachemodel.NewAccessTokenModel(sqlx.NewMysql(c.DB.DataSource)),
		UserPointsModel:    cachemodel.NewUserPointsModel(sqlx.NewMysql(c.DB.DataSource)),
	}
}
