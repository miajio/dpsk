package md5_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/miajio/dpsk/pkg/md5"
)

func TestMD5(t *testing.T) {
	md5util := md5.New()

	// 字符串哈希
	str := "123456"
	hash := md5util.SumString(str)
	fmt.Printf("String MD5: %s\n", hash)

	// 验证字符串
	valid, err := md5util.VerifyString(str, hash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Verification: %v\n", valid)

	// 文件哈希
	fileHash, err := md5util.SumFile("C:\\Users\\Administrator\\Desktop\\20250422000054.csv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File MD5: %s\n", fileHash)

	// 大文件哈希 (自定义选项)
	largeFileHash, err := md5util.SumLargeFile("C:\\Users\\Administrator\\.lingma\\model\\0.1.0\\java.model",
		md5.WithChunkSize(2<<20), // 2MB chunks
		md5.WithWorkers(8),       // 8 workers
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Large File MD5: %s\n", largeFileHash)
}
