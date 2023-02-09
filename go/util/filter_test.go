package util

import (
	"fmt"
	"testing"
)

func TestInitSensitiveFilter(t *testing.T) {
	InitSensitiveFilter()
}

func TestInitSensitiveFilter2(t *testing.T) {
	InitSensitiveFilter()
	filter, _ := SensitiveWordsFilter("傻逼")
	fmt.Printf("%s", filter)
}

func TestTest(t *testing.T) {
	Test()
}
