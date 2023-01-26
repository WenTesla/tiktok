package model

import (
	"log"
	"testing"
	"time"
)

func TestGetVideoByLastTime(t *testing.T) {
	now := time.Now()
	videos, err := GetVideoByLastTime(now)
	log.Printf("%v", videos)
	log.Printf("%v", err)
}
func TestGetVideoNextTime(t *testing.T) {
	now := time.Now()
	lastTime, err := GetVideoNextTime(now)
	log.Printf("%v", lastTime)
	log.Printf("%v", err)
}
