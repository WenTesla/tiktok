package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestQueryCommentCountsByVideoIds(t *testing.T) {
	mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM `comments`").
		WithArgs(int64(3), int64(5), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"video_id", "cnt"}).
			AddRow(int64(3), int64(1)).
			AddRow(int64(5), int64(9)))
	got, err := QueryCommentCountsByVideoIds([]int64{3, 5})
	require.NoError(t, err)
	require.Equal(t, int64(1), got[3])
	require.Equal(t, int64(9), got[5])
	require.NoError(t, mock.ExpectationsWereMet())
}
