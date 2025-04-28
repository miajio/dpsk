package md5

// FileHasher 文件哈希值计算器
type FileHasher interface {
	SumFile(filepath string) (string, error)                          // 计算文件的哈希值
	SumLargeFile(filePath string, opts ...FileOption) (string, error) // 计算大文件的哈希值
}

// fileOptions 文件处理选项
type fileOptions struct {
	chunkSize int64
	workers   int
}

// FileOption 文件处理选项函数类型
type FileOption func(*fileOptions)

// WithChunkSize 设置分块大小
func WithChunkSize(size int64) FileOption {
	return func(o *fileOptions) {
		o.chunkSize = size
	}
}

// WithWorkers 设置工作协程数
func WithWorkers(n int) FileOption {
	return func(o *fileOptions) {
		o.workers = n
	}
}

// defaultFileOptions 默认文件处理选项
func defaultFileOptions() *fileOptions {
	return &fileOptions{
		chunkSize: 1 << 20, // 1MB
		workers:   4,
	}
}
