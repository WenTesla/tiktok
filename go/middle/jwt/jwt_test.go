package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tiktok/go/config"
	"tiktok/go/model"
)

// 确保签名有效期足够长，避免测试时 token 立即过期（config 在测试中可能为默认值 0）。
func init() {
	config.AppConfig.JWT.Expire = 24
}

// 签名 -> 解析 往返测试
func TestSignTokenAndParse(t *testing.T) {
	token := SignToken(model.User{Id: 42})
	require.NotEmpty(t, token)

	id, err := ParseToken(token)
	require.NoError(t, err)
	// JSON 数字解码为 float64
	assert.Equal(t, float64(42), id)
}

// 过期 token 应被自动拒绝（验证 exp 声明生效）
func TestParseTokenExpired(t *testing.T) {
	claims := jwt.MapClaims{
		"id": 7,
		"exp": time.Now().Add(-time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(jwtSecret)
	require.NoError(t, err)

	_, err = ParseToken(signed)
	assert.Error(t, err)
}

// ParseTokenTime 应能从合法 token 中取出 nbf
func TestParseTokenTime(t *testing.T) {
	token := SignToken(model.User{Id: 1})
	nbf, err := ParseTokenTime(token)
	require.NoError(t, err)
	assert.NotNil(t, nbf)
}

// GetDays 天数差计算
func TestGetDays(t *testing.T) {
	// 相差约 9 天
	d := GetDays(1675998098, 1676775698)
	assert.Equal(t, 9, d)
}
