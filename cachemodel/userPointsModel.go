package cachemodel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserPointsModel = (*customUserPointsModel)(nil)

type (
	// UserPointsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserPointsModel.
	UserPointsModel interface {
		userPointsModel
	}

	customUserPointsModel struct {
		*defaultUserPointsModel
	}
)

// NewUserPointsModel returns a model for the database table.
func NewUserPointsModel(conn sqlx.SqlConn) UserPointsModel {
	return &customUserPointsModel{
		defaultUserPointsModel: newUserPointsModel(conn),
	}
}
