package cachemodel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AccessTokenModel = (*customAccessTokenModel)(nil)

type (
	// AccessTokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccessTokenModel.
	AccessTokenModel interface {
		accessTokenModel
	}

	customAccessTokenModel struct {
		*defaultAccessTokenModel
	}
)

// NewAccessTokenModel returns a model for the database table.
func NewAccessTokenModel(conn sqlx.SqlConn) AccessTokenModel {
	return &customAccessTokenModel{
		defaultAccessTokenModel: newAccessTokenModel(conn),
	}
}
