package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	Id           int64  `gorm:"column:id;type:bigint;primaryKey"`
	CreateUserId int64  `gorm:"column:createUserId;type:bigint;default:0"`                        // 创建人
	CreateTime   int64  `gorm:"column:createTime;type:bigint;autoCreateTime:nano"`                // 创建时间
	ModifyUserId int64  `gorm:"column:modifyUserId;type:bigint;default:0"`                        // 修改人
	ModifyTime   int64  `gorm:"column:modifyTime;type:bigint;default:0"`                          // 修改时间
	IsDeleted    string `gorm:"column:isDeleted;type:enum('NORMAL', 'DELETED') default 'NORMAL'"` // 删除状态 NORMAL 正常 DELETED 删除
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.IsDeleted == "" {
		m.IsDeleted = "NORMAL"
	}
	return
}

func (m *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifyTime = int64(time.Now().UnixNano())
	return
}

const (
	NORMAL  = "NORMAL"
	DELETED = "DELETED"
	YES     = "YES"
	NO      = "NO"
	ALL     = "ALL"
)
