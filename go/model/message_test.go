package model

import (
	"fmt"
	"testing"
)

func TestInsertMessage(t *testing.T) {
	InsertMessage(5, 1, "111")
}

func TestQueryMessageByUserId(t *testing.T) {
	QueryMessageByUserId(1)
}
func TestQert(t *testing.T) {
	_, i, _ := QueryNewestMessageByUserId(1)
	fmt.Println(i)
}
func TestQueryNewestMessageByUserIdAndToUserID(t *testing.T) {
	QueryNewestMessageByUserIdAndToUserID(1, 2)
}
func TestQueryMessageByUserIdAndToUserId(t *testing.T) {
	QueryMessageByUserIdAndToUserId(1, 2)
}
