package models

import (
	"gorm.io/gorm"
)

// SysFile 系统文件表
type SysFile struct {
	gorm.Model
	ID            int64         `gorm:"primaryKey;autoIncrement"`                       // 文件ID
	Name          string        `gorm:"size:256;not null"`                              // 文件名称
	Type          string        `gorm:"size:128;not null"`                              // 文件类型
	Extension     string        `gorm:"size:64;not null"`                               // 文件扩展名
	RelativePath  string        `gorm:"size:256;not null"`                              // 文件相对路径
	AbsoluteURL   string        `gorm:"size:256;not null"`                              // 文件绝对路径
	Hash          string        `gorm:"size:32;not null;index"`                         // 文件哈希值
	Size          int64         `gorm:"not null"`                                       // 文件大小
	Status        FileStatus    `gorm:"type:enum('NORMAL','DISABLE');default:'NORMAL'"` // 文件状态
	IsDeleted     DeletedStatus `gorm:"type:enum('NORMAL','DELETED');default:'NORMAL'"` // 文件删除状态
	DisableReason string        `gorm:"size:256"`                                       // 文件禁用原因
}

func (f *SysFile) TableName() string {
	return "sys_file"
}

// 在model/file.go中添加
type SysChunkFile struct {
	gorm.Model
	UploadID    string `gorm:"size:36;not null;index"`
	ChunkNumber int    `gorm:"not null"`
	TotalChunks int    `gorm:"not null"`
	FilePath    string `gorm:"size:500;not null"`
	FileType    string `gorm:"size:100;not null"`
	FileName    string `gorm:"size:255;not null"`
	FileSize    int64  `gorm:"not null"`
	ChunkSize   int64  `gorm:"not null"`
	Hash        string `gorm:"size:32;not null"`
	IsComplete  bool   `gorm:"default:false"`
}

func (f *SysChunkFile) TableName() string {
	return "sys_chunk_file"
}
