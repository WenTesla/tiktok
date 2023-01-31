package model

import (
	"fmt"
	"testing"
)

func TestQueryCommentByUserId(t *testing.T) {
	comments, err := QueryCommentByUserId(2)
	fmt.Printf("%v", comments)
	fmt.Printf("%v", err)
}
func TestQueryCommentByVideoId(t *testing.T) {
	comments, err := QueryCommentByVideoId(2)
	fmt.Printf("%v", comments)
	fmt.Printf("%v", err)
}
