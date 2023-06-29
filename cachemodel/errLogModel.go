package cachemodel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ErrLogModel = (*customErrLogModel)(nil)

type (
	// ErrLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customErrLogModel.
	ErrLogModel interface {
		errLogModel
	}

	customErrLogModel struct {
		*defaultErrLogModel
	}
)

// NewErrLogModel returns a model for the database table.
func NewErrLogModel(conn sqlx.SqlConn) ErrLogModel {
	return &customErrLogModel{
		defaultErrLogModel: newErrLogModel(conn),
	}
}
