package cachemodel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserinfosModel = (*customUserinfosModel)(nil)

type (
	// UserinfosModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserinfosModel.
	UserinfosModel interface {
		userinfosModel
	}

	customUserinfosModel struct {
		*defaultUserinfosModel
	}
)

// NewUserinfosModel returns a model for the database table.
func NewUserinfosModel(conn sqlx.SqlConn) UserinfosModel {
	return &customUserinfosModel{
		defaultUserinfosModel: newUserinfosModel(conn),
	}
}
