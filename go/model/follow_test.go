package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetFollowCountsByUserIds(t *testing.T) {
	mock := newMockDB(t)
	// 关注数
	mock.ExpectQuery("SELECT .* FROM `follows`").
		WithArgs(int64(1), int64(2), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"uid", "cnt"}).AddRow(int64(1), int64(3)).AddRow(int64(2), int64(4)))
	// 粉丝数
	mock.ExpectQuery("SELECT .* FROM `follows`").
		WithArgs(int64(1), int64(2), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"uid", "cnt"}).AddRow(int64(1), int64(10)).AddRow(int64(2), int64(20)))
	following, fans, err := GetFollowCountsByUserIds([]int64{1, 2})
	require.NoError(t, err)
	require.Equal(t, int64(3), following[1])
	require.Equal(t, int64(4), following[2])
	require.Equal(t, int64(10), fans[1])
	require.Equal(t, int64(20), fans[2])
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryIsFollowBatch(t *testing.T) {
	mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM `follows`").
		WithArgs(int64(9), int64(1), int64(2), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"follower_id"}).AddRow(int64(1)))
	got, err := QueryIsFollowBatch(9, []int64{1, 2})
	require.NoError(t, err)
	require.True(t, got[1])
	require.False(t, got[2])
	require.NoError(t, mock.ExpectationsWereMet())
}
