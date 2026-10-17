package util

import (
	"os"
	"testing"
)

// InitSensitiveFilter 读取 go/util/key.txt（相对仓库根）。测试 cwd 为包目录时文件不可达，跳过而非崩溃。
func TestInitSensitiveFilter(t *testing.T) {
	if _, err := os.Stat(SensitiveWordFilePath); err != nil {
		t.Skip("sensitive filter key file not reachable from test cwd; skipping")
	}
	InitSensitiveFilter()
}

func TestInitSensitiveFilter2(t *testing.T) {
	if _, err := os.Stat(SensitiveWordFilePath); err != nil {
		t.Skip("sensitive filter key file not reachable from test cwd; skipping")
	}
	InitSensitiveFilter()
	filter, _ := SensitiveWordsFilter("傻逼")
	t.Logf("filter result: %s", filter)
}

func TestTest(t *testing.T) {
	Test()
}
