package cachemodel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserBehaviourLogModel = (*customUserBehaviourLogModel)(nil)

type (
	// UserBehaviourLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserBehaviourLogModel.
	UserBehaviourLogModel interface {
		userBehaviourLogModel
	}

	customUserBehaviourLogModel struct {
		*defaultUserBehaviourLogModel
	}
)

// NewUserBehaviourLogModel returns a model for the database table.
func NewUserBehaviourLogModel(conn sqlx.SqlConn, c cache.CacheConf) UserBehaviourLogModel {
	return &customUserBehaviourLogModel{
		defaultUserBehaviourLogModel: newUserBehaviourLogModel(conn, c),
	}
}
