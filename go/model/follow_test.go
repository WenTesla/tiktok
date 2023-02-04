package model

import (
	"log"
	"testing"
)

func TestGetFollowingById(t *testing.T) {
	id, err := GetFollowingById(2)
	log.Println(err)
	log.Println(id)
}

func TestGetFansById(t *testing.T) {
	id, err := GetFansById(2)
	log.Println(id)
	log.Println(err)
}
func TestInsertFollow(t *testing.T) {
	InsertFollow(4, 2)
}
func TestCancelFollow(t *testing.T) {
	CancelFollow(1, 2)
}
func TestQueryFollowUsersByUserId(t *testing.T) {
	users, err := QueryFollowUsersByUserId(1)
	if err != nil {
		log.Println(err)
	}
	log.Println(users)
}
