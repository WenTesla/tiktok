package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 纯函数测试：bcrypt 哈希与校验
func TestHashPasswordAndCheck(t *testing.T) {
	hash, err := HashPassword("secret123")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, CheckPassword(hash, "secret123"))
	assert.False(t, CheckPassword(hash, "wrong-password"))
}

// 纯函数测试：邮箱格式校验
func TestVerifyEmailFormat(t *testing.T) {
	assert.True(t, VerifyEmailFormat("alice@example.com"))
	assert.False(t, VerifyEmailFormat("not-an-email"))
}

// 纯函数测试：手机号格式校验
func TestVerifyMobileFormat(t *testing.T) {
	assert.True(t, VerifyMobileFormat("13800138000"))
	assert.False(t, VerifyMobileFormat("12345"))
}

// service 层测试（go-sqlmock）：UserServiceBatch 应一次性批量取数，消除 N+1
func TestUserServiceBatchLoggedOut(t *testing.T) {
	mock := newMockDB(t)
	// 1) 批量取作者基础信息
	mock.ExpectQuery("SELECT .* FROM `users`").WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"}).
			AddRow(int64(1), "alice", "h1").AddRow(int64(2), "bob", "h2"))
	// 2) 关注数
	mock.ExpectQuery("SELECT .* FROM `follows`").WithArgs(int64(1), int64(2), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"uid", "cnt"}).AddRow(int64(1), int64(3)).AddRow(int64(2), int64(4)))
	// 3) 粉丝数
	mock.ExpectQuery("SELECT .* FROM `follows`").WithArgs(int64(1), int64(2), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"uid", "cnt"}).AddRow(int64(1), int64(10)).AddRow(int64(2), int64(20)))
	// 4) 点赞他人数
	mock.ExpectQuery("SELECT .* FROM `likes`").WithArgs(int64(1), int64(2), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"uid", "cnt"}).AddRow(int64(1), int64(2)).AddRow(int64(2), int64(3)))
	// 5) 获赞总数
	mock.ExpectQuery("SELECT v.author_id as uid").WithArgs(int64(0), int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"uid", "cnt"}).AddRow(int64(1), int64(20)).AddRow(int64(2), int64(30)))
	// 6) 作品数
	mock.ExpectQuery("SELECT .* FROM `videos`").WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"author_id", "cnt"}).AddRow(int64(1), int64(5)).AddRow(int64(2), int64(8)))

	infos, err := UserServiceBatch([]int64{1, 2}, -1)
	require.NoError(t, err)
	require.Len(t, infos, 2)

	alice := infos[1]
	assert.Equal(t, "alice", alice.Name)
	assert.Equal(t, int64(3), alice.FollowCount)
	assert.Equal(t, int64(10), alice.FollowerCount)
	assert.Equal(t, int64(2), alice.FavoriteCount)
	assert.Equal(t, int64(20), alice.TotalFavorited)
	assert.Equal(t, int64(5), alice.WorkCount)
	assert.False(t, alice.IsFollow) // 未登录态

	require.NoError(t, mock.ExpectationsWereMet())
}
