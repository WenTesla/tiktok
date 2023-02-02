package model

import (
	"log"
	"testing"
)

func TestGetFollowingById(t *testing.T) {
	id, err := GetFollowingById(20053)
	log.Println(err)
	log.Println(id)
}

func TestGetFansById(t *testing.T) {
	id, err := GetFansById(20053)
	log.Println(id)
	log.Println(err)
}
func TestInsertFollow(t *testing.T) {
	InsertFollow(4, 2)
}
