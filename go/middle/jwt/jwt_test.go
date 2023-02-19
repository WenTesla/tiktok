package jwt

import (
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"testing"
	"time"
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

func TestParseTokenTime(t *testing.T) {
	time, err := ParseTokenTime("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiaXNzIjoidGlrdG9rIiwibmJmIjoxNjc1MTU5ODM4LCJzdWIiOiLnlKjmiLdJZCJ9.sPKBddA3foAeVVH19QFNoDNIckUIIipfZOrGvKLHxQk")
	if err != nil {
		t.Fail()
	}
	fmt.Println(time)
}

func TestGetDays(t *testing.T) {
	days := GetDays(1676775698, 1675998098)
	fmt.Println(days)
}

func TestLimit(t *testing.T) {
	r := rate.Every(1 * time.Millisecond)
	limit := rate.NewLimiter(r, 10)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if limit.Allow() {
			fmt.Printf("请求成功，当前时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("请求成功，但是被限流了。。。\n")
		}
	})

	_ = http.ListenAndServe(":8081", nil)
}

func GetApi() {
	api := "http://localhost:8081/"
	res, err := http.Get(api)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Printf("get api success\n")
	}
}

func Benchmark_Main(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetApi()
	}
}
