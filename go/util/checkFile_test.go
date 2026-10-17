package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileType(t *testing.T) {
	// mp4 文件头魔数：00 00 00 20 66 74 79 70 6d 70（"ftypmp"）
	mp4 := []byte{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, 0x6d, 0x70}
	assert.Equal(t, "mp4", GetFileType(mp4))

	// PDF 文件头魔数：%PDF-1.5（与 fileTypeMap 中存储的魔数一致）
	pdf := []byte("%PDF-1.5")
	assert.Equal(t, "pdf", GetFileType(pdf))

	// 普通文本不应被判为 mp4（GetFileType 是启发式，可能命中其它类型，这里只验证不是 mp4）
	assert.NotEqual(t, "mp4", GetFileType([]byte("hello world")))

	// 空输入返回空
	assert.Equal(t, "", GetFileType(nil))
	assert.Equal(t, "", GetFileType([]byte{}))
}
