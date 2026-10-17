package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalStoragePut(t *testing.T) {
	dir := t.TempDir()
	s := NewLocal(dir, "http://example/static")

	const rel = "videos/1_123.mp4"
	url, err := s.Put(rel, strings.NewReader("fake-video-bytes"))
	if err != nil {
		t.Fatalf("Put 失败: %v", err)
	}

	// 文件应落在 BaseDir/rel
	full := filepath.Join(dir, rel)
	if _, err := os.Stat(full); err != nil {
		t.Fatalf("文件未落盘: %v", err)
	}
	b, _ := os.ReadFile(full)
	if string(b) != "fake-video-bytes" {
		t.Fatalf("内容不符: %q", string(b))
	}

	// URL 应以 PublicBaseURL 开头并包含 rel
	want := "http://example/static/" + rel
	if url != want {
		t.Fatalf("URL 不符: got %q want %q", url, want)
	}
}

func TestLocalStorageMkdirAll(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	s := NewLocal(dir, "http://example/static")
	if err := s.MkdirAll(); err != nil {
		t.Fatalf("MkdirAll 失败: %v", err)
	}
	for _, sub := range []string{"videos", "covers", "tmp"} {
		if _, err := os.Stat(filepath.Join(dir, sub)); err != nil {
			t.Fatalf("子目录缺失 %s: %v", sub, err)
		}
	}
}

func TestLocalStorageLocalPath(t *testing.T) {
	s := NewLocal("/data", "http://h/static")
	want := filepath.Join("/data", "covers", "x.jpg")
	if got := s.LocalPath("covers/x.jpg"); got != want {
		t.Fatalf("LocalPath 不符: got %q want %q", got, want)
	}
}
