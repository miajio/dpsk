package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	Id           uint64 `gorm:"primaryKey"`
	CreateUserId uint64 `gorm:"default:0"`                                       // 创建人
	CreateTime   uint64 `gorm:"autoCreateTime:nano"`                             // 创建时间
	ModifyUserId uint64 `gorm:"default:0"`                                       // 修改人
	ModifyTime   uint64 `gorm:"default:0"`                                       // 修改时间
	IsDeleted    string `gorm:"type:enum('NORMAL', 'DELETED') default 'NORMAL'"` // 删除状态 NORMAL 正常 DELETED 删除
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.IsDeleted == "" {
		m.IsDeleted = "NORMAL"
	}
	return
}

func (m *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifyTime = uint64(time.Now().UnixNano())
	return
}

const (
	NORMAL  = "NORMAL"
	DELETED = "DELETED"
	YES     = "YES"
	NO      = "NO"
	ALL     = "ALL"
)
