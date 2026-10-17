package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// 注意：GORM 生成的 SQL 会给列/表名加反引号，这里用表级宽松正则避免引号差异导致匹配失败。
func TestQueryLikeCountsByVideoIds(t *testing.T) {
	mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM `likes`").
		WithArgs(int64(3), int64(5), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"video_id", "cnt"}).
			AddRow(int64(3), int64(7)).
			AddRow(int64(5), int64(2)))
	got, err := QueryLikeCountsByVideoIds([]int64{3, 5})
	require.NoError(t, err)
	require.Equal(t, int64(7), got[3])
	require.Equal(t, int64(2), got[5])
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryIsLikeBatch(t *testing.T) {
	mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM `likes`").
		WithArgs(int64(9), int64(3), int64(5), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"video_id"}).
			AddRow(int64(3)).
			AddRow(int64(5)))
	got, err := QueryIsLikeBatch(9, []int64{3, 5})
	require.NoError(t, err)
	require.True(t, got[3])
	require.True(t, got[5])
	require.False(t, got[7])
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLikeCountsByUserIds(t *testing.T) {
	mock := newMockDB(t)
	// 点赞他人数
	mock.ExpectQuery("SELECT .* FROM `likes`").
		WithArgs(int64(1), int64(2), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"uid", "cnt"}).AddRow(int64(1), int64(2)).AddRow(int64(2), int64(3)))
	// 获赞总数（JOIN videos）
	mock.ExpectQuery("SELECT v.author_id as uid").
		WithArgs(int64(0), int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"uid", "cnt"}).AddRow(int64(1), int64(20)).AddRow(int64(2), int64(30)))
	fav, tot, err := GetLikeCountsByUserIds([]int64{1, 2})
	require.NoError(t, err)
	require.Equal(t, int64(2), fav[1])
	require.Equal(t, int64(3), fav[2])
	require.Equal(t, int64(20), tot[1])
	require.Equal(t, int64(30), tot[2])
	require.NoError(t, mock.ExpectationsWereMet())
}
