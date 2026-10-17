package util

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// ffmpegBin 容器内 ffmpeg 静态二进制路径（FROM scratch 镜像通过 COPY 注入，无 shell，绝对路径调用）。
const ffmpegBin = "/src/ffmpeg"

// ExtractFirstFrame 从 inPath 视频抽取首帧，并缩放到 width 宽（保持比例），输出为 outPath（jpg）。
// ffmpeg 二进制缺失时返回明确错误（不 panic），调用方据此降级处理（封面留空，视频仍发布）。
func ExtractFirstFrame(inPath, outPath string, width int) error {
	// -tmpdir 指向输出文件所在目录（需可写）；单帧抽取几乎无临时文件，这里仅为稳健。
	tmpDir := filepath.Dir(outPath)
	cmd := exec.Command(ffmpegBin,
		"-y",
		"-i", inPath,
		"-vf", fmt.Sprintf("scale=%d:-1", width), // 缩放到 width 宽，高度按比例
		"-frames:v", "1",
		"-q:v", "2", // 封面质量
		outPath,
		"-tmpdir", tmpDir,
	)
	if err := cmd.Run(); err != nil {
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return fmt.Errorf("ffmpeg 未找到：需将静态二进制 ffmpeg 注入镜像 %s: %w", ffmpegBin, err)
		}
		return fmt.Errorf("ffmpeg 抽帧失败: %w", err)
	}
	return nil
}

// FfmpegAvailable 探测 ffmpeg 二进制是否存在（供测试/启动自检使用）。
func FfmpegAvailable() bool {
	_, err := exec.LookPath(ffmpegBin)
	if err == nil {
		return true
	}
	// LookPath 在部分 scratch 环境对绝对路径不可靠，再尝试直接 stat
	cmd := exec.Command(ffmpegBin, "-version")
	return cmd.Run() == nil
}
