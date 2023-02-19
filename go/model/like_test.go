package model

import (
	"fmt"
	"testing"
)

func TestLikeVideoByUserId(t *testing.T) {

}

func TestUnLikeVideoByUserId(t *testing.T) {
	err := UpdateLikeVideoByUserId(2, 1, 1)
	fmt.Printf("%v", err)
}
func TestQueryDuplicateLikeData(t *testing.T) {
	data, err := QueryDuplicateLikeData(1, 10)
	fmt.Printf("%v %v", data, err)
}
func TestQueryVideoByUserId(t *testing.T) {
	QueryVideoByUserId(2)
}
func TestQueryLikeByVideoId(t *testing.T) {
	count, _ := QueryLikeByVideoId(3)
	fmt.Print(count)
}
func TestQueryTotalFavorited(t *testing.T) {
	QueryTotalFavorited(1)
}
