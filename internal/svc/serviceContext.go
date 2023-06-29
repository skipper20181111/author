package svc

import (
	"author/cachemodel"
	"author/internal/config"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

const localCacheExpire = time.Duration(time.Second * 100)
const Tockenid = 1

type ServiceContext struct {
	Config             config.Config
	UserModel          cachemodel.UserinfosModel
	UserBehaviourModel cachemodel.UserBehaviourLogModel
	AccessTokenModel   cachemodel.AccessTokenModel
	UserPointsModel    cachemodel.UserPointsModel
	ErrLog             cachemodel.ErrLogModel
	LocalCache         *collection.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	localCache, err := collection.NewCache(localCacheExpire)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:             c,
		UserModel:          cachemodel.NewUserinfosModel(sqlx.NewMysql(c.DB.DataSource)),
		UserBehaviourModel: cachemodel.NewUserBehaviourLogModel(sqlx.NewMysql(c.DB.DataSource)),
		AccessTokenModel:   cachemodel.NewAccessTokenModel(sqlx.NewMysql(c.DB.DataSource)),
		UserPointsModel:    cachemodel.NewUserPointsModel(sqlx.NewMysql(c.DB.DataSource)),
		ErrLog:             cachemodel.NewErrLogModel(sqlx.NewMysql(c.DB.DataSource)),
		LocalCache:         localCache,
	}
}
