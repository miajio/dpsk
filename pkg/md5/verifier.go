package md5

// Verifier 验证接口
type Verifier interface {
	Verify(data []byte, hash string) (bool, error)                                  // 验证数据与哈希是否匹配
	VerifyString(s string, hash string) (bool, error)                               // 验证字符串与哈希是否匹配
	VerifyFile(filePath string, hash string) (bool, error)                          // 验证文件内容与哈希是否匹配
	VerifyLargeFile(filePath string, hash string, opts ...FileOption) (bool, error) // 验证大文件内容与哈希是否匹配
}
