package cachemodel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ WxDeliveryModel = (*customWxDeliveryModel)(nil)

type (
	// WxDeliveryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWxDeliveryModel.
	WxDeliveryModel interface {
		wxDeliveryModel
	}

	customWxDeliveryModel struct {
		*defaultWxDeliveryModel
	}
)

// NewWxDeliveryModel returns a model for the database table.
func NewWxDeliveryModel(conn sqlx.SqlConn) WxDeliveryModel {
	return &customWxDeliveryModel{
		defaultWxDeliveryModel: newWxDeliveryModel(conn),
	}
}
