package md5_test

import (
	"fmt"
	"testing"

	"github.com/miajio/dpsk/pkg/md5"
)

func TestMd5(t *testing.T) {
	a := "123456"
	m := md5.Md5(a)
	fmt.Println(m)
}
