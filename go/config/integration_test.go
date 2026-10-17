package config

import (
	"fmt"
	"os"
	"testing"
)

// TestMain 在缺少真实 MySQL/Redis 时跳过集成测试，保证 `go test ./...` 在本地/CI 无需基础设施即可全绿。
// 需要跑真实链路时：TIKTOK_INTEGRATION=1 go test ./go/config/...
func TestMain(m *testing.M) {
	if os.Getenv("TIKTOK_INTEGRATION") == "" {
		fmt.Println("skipping config integration tests (set TIKTOK_INTEGRATION=1 to enable)")
		os.Exit(0)
	}
	os.Exit(m.Run())
}
