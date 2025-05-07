package models

// 文件状态
type FileStatus string

// 删除状态
type DeletedStatus string

const (
	FileStatusNormal  FileStatus = "NORMAL"  // 文件状态: 正常
	FileStatusDisable FileStatus = "DISABLE" // 文件状态: 禁用

	DeleteNormal   DeletedStatus = "NORMAL"  // 删除状态: 正常
	DeletedDeleted DeletedStatus = "DELETED" // 删除状态: 删除
)

// IsDeleted 是删除状态
func (d DeletedStatus) IsDeleted() bool {
	return d == DeletedDeleted
}

// IsNormal 是正常状态
func (f DeletedStatus) IsNormal() bool {
	return f == DeletedDeleted
}
