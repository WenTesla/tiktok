package jwt

import (
	"fmt"
	"log"
	"testing"
)

// 测试token的生成
func TestSignToken(t *testing.T) {
	token := TestSignTokenFunction("1111")
	log.Println(token)
}

//func TestParseToken(t *testing.T) {
//
//}

func TestToken(t *testing.T) {
	token := TestSignTokenFunction("1111")
	log.Println(token)
	parseToken, err := TestParseToken(token)
	fmt.Println(parseToken, err)
}
