package util

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFfmpegAvailable 仅记录 ffmpeg 是否可用，不强制（缺失时 publish 会降级处理）。
func TestFfmpegAvailable(t *testing.T) {
	t.Logf("ffmpeg available = %v", FfmpegAvailable())
}

// TestExtractFirstFrameInvalidInput 输入文件不存在时，函数必须返回错误而非静默成功
// （无论 ffmpeg 是否在镜像中），且不应产生输出文件。
func TestExtractFirstFrameInvalidInput(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.jpg")
	err := ExtractFirstFrame("/nonexistent/path/to/video.mp4", out, 480)
	if err == nil {
		t.Fatalf("期望错误，但 ExtractFirstFrame 返回 nil")
	}
	if _, statErr := os.Stat(out); statErr == nil {
		t.Fatalf("期望无输出文件，但 %s 已存在", out)
	}
}
