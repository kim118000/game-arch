package model

import (
	"time"
)

type Base struct {
	ID         uint64    `json:"id,omitempty"            gorm:"type:BIGINT;column:id;primaryKey;autoIncrement"`
	RoleId     uint64    `json:"roleId,omitempty"        gorm:"type:BIGINT;column:role_id;not null;index:idx_role_id"`
	CreateTime time.Time `json:"createTime,omitempty"    gorm:"type:TIMESTAMP;column:create_time"`
	UpdateTime time.Time `json:"updateTime,omitempty"    gorm:"type:TIMESTAMP;column:update_time"`
}

func (base *Base) IsExist() bool {
	return base.ID > 0
}
