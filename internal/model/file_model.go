package model

type FileModel struct {
	Model
	Name         string `gorm:"column:name;type:varchar(64)"`          // 文件名称
	Path         string `gorm:"column:path;type:varchar(128)"`         // 文件存储路径
	Url          string `gorm:"column:url;type:varchar(256)"`          // 文件访问地址
	Size         int64  `gorm:"column:size;type:bigint"`               // 文件大小
	Hash         string `gorm:"column:hash;type:varchar(64)"`          // 文件哈希值
	Extension    string `gorm:"column:extension;type:varchar(16)"`     // 文件后缀
	BusinessType string `gorm:"column:business_type;type:varchar(32)"` // 业务类型
}

func (FileModel) TableName() string {
	return "file_info"
}
