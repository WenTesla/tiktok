// Package storage 提供文件存储抽象，当前实现为本地文件系统存储（宿主卷挂载）。
// 视频/封面落盘后通过 gin 静态路由对外提供，PlayUrl/CoverUrl 存公网可达的本地 URL。
package storage

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Storage 存储抽象：Put 写入一个相对路径并返回可访问的公网 URL。
type Storage interface {
	// Put 将 reader 内容写入 relPath（相对存储根目录），返回公网可访问的 URL。
	Put(relPath string, r io.Reader) (publicURL string, err error)
	// LocalPath 返回 relPath 在容器/宿主上的绝对路径（供 ffmpeg 等本地进程直接读取）。
	LocalPath(relPath string) string
	// MkdirAll 创建运行所需的子目录（videos/covers/tmp）。
	MkdirAll() error
}

// LocalStorage 本地文件系统实现，BaseDir 为容器内切面（如 /data）。
type LocalStorage struct {
	BaseDir       string
	PublicBaseURL string
}

// NewLocal 构造本地存储。
func NewLocal(baseDir, publicBaseURL string) *LocalStorage {
	return &LocalStorage{BaseDir: baseDir, PublicBaseURL: publicBaseURL}
}

// LocalPath 返回 relPath 在本地磁盘上的绝对路径。
func (s *LocalStorage) LocalPath(relPath string) string {
	return filepath.Join(s.BaseDir, relPath)
}

// Put 将 reader 写入 BaseDir/relPath，并支持创建父目录；返回 PublicBaseURL/relPath。
func (s *LocalStorage) Put(relPath string, r io.Reader) (string, error) {
	full := s.LocalPath(relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", err
	}
	f, err := os.Create(full)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	return strings.TrimSuffix(s.PublicBaseURL, "/") + "/" + relPath, nil
}

// MkdirAll 创建运行所需的子目录：videos / covers / tmp。
func (s *LocalStorage) MkdirAll() error {
	for _, sub := range []string{"videos", "covers", "tmp"} {
		if err := os.MkdirAll(filepath.Join(s.BaseDir, sub), 0o755); err != nil {
			return err
		}
	}
	return nil
}
