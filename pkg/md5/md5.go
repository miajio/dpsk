package md5

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"os"
	"sync"
)

// MD5 MD5处理器
type MD5 struct {
}

// New MD5实例
func New() *MD5 {
	return &MD5{}
}

// Sum 计算数据的 MD5 哈希值
func (m *MD5) Sum(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

// SumString 计算字符串的 MD5 哈希值
func (m *MD5) SumString(s string) string {
	return m.Sum([]byte(s))
}

// SumFile 计算整个文件的 MD5 哈希值
func (m *MD5) SumFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SumLargeFile 并发计算大文件的 MD5 哈希值
func (m *MD5) SumLargeFile(filePath string, opts ...FileOption) (string, error) {
	options := defaultFileOptions()
	for _, opt := range opts {
		opt(options)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	fileSize := fileInfo.Size()
	if options.chunkSize <= 0 {
		options.chunkSize = fileSize / int64(options.workers)
		if options.chunkSize == 0 {
			options.chunkSize = fileSize
		}
	}

	type chunk struct {
		index int
		data  []byte
	}

	chunks := make(chan chunk, options.workers)
	results := make(chan string, options.workers)
	errChan := make(chan error, 1)

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < options.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for c := range chunks {
				hash := md5.Sum(c.data)
				results <- hex.EncodeToString(hash[:])
			}
		}()
	}

	// 读取文件分块
	go func() {
		defer close(chunks)
		position := int64(0)
		index := 0
		for position < fileSize {
			size := options.chunkSize
			if position+size > fileSize {
				size = fileSize - position
			}

			buf := make([]byte, size)
			_, err := file.ReadAt(buf, position)
			if err != nil && err != io.EOF {
				errChan <- err
				return
			}

			chunks <- chunk{index: index, data: buf}
			position += size
			index++
		}
	}()

	// 收集结果
	go func() {
		wg.Wait()
		close(results)
	}()

	// 合并所有分块的哈希
	finalHash := md5.New()
	for result := range results {
		partHash, err := hex.DecodeString(result)
		if err != nil {
			return "", err
		}
		finalHash.Write(partHash)
	}

	select {
	case err := <-errChan:
		return "", err
	default:
		return hex.EncodeToString(finalHash.Sum(nil)), nil
	}
}

// Verify 验证数据与哈希是否匹配
func (m *MD5) Verify(data []byte, hash string) (bool, error) {
	decoded, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	dataHash := md5.Sum(data)
	return subtle.ConstantTimeCompare(dataHash[:], decoded) == 1, nil
}

// VerifyString 验证字符串与哈希是否匹配
func (m *MD5) VerifyString(s string, hash string) (bool, error) {
	return m.Verify([]byte(s), hash)
}

// VerifyFile 验证文件内容与哈希是否匹配
func (m *MD5) VerifyFile(filePath string, hash string) (bool, error) {
	fileHash, err := m.SumFile(filePath)
	if err != nil {
		return false, err
	}
	return fileHash == hash, nil
}

// VerifyLargeFile 验证大文件内容与哈希是否匹配
func (m *MD5) VerifyLargeFile(filePath string, hash string, opts ...FileOption) (bool, error) {
	fileHash, err := m.SumLargeFile(filePath, opts...)
	if err != nil {
		return false, err
	}
	return fileHash == hash, nil
}
