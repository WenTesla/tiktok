package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetWorkCountsByUserIds(t *testing.T) {
	mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM `videos`").
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"author_id", "cnt"}).
			AddRow(int64(1), int64(5)).
			AddRow(int64(2), int64(8)))
	got, err := GetWorkCountsByUserIds([]int64{1, 2})
	require.NoError(t, err)
	require.Equal(t, int64(5), got[1])
	require.Equal(t, int64(8), got[2])
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIds(t *testing.T) {
	mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM `users`").
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"}).
			AddRow(int64(1), "alice", "hash1").
			AddRow(int64(2), "bob", "hash2"))
	users, err := GetUserByIds([]int64{1, 2})
	require.NoError(t, err)
	require.Len(t, users, 2)
	require.Equal(t, "alice", users[0].Name)
	require.Equal(t, "bob", users[1].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}
