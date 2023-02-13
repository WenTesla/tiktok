package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const SensitiveWordFilePath = "go/util/key.txt"

var sensitiveWords []string

//用于过滤

var wordReg *regexp.Regexp

func InitSensitiveFilter() {
	//sensitiveWords := []string{
	//	"傻逼",
	//	"傻叉",
	//	"垃圾",
	//	"妈的",
	//	"sb",
	//}
	//text := "什么垃圾打野，傻逼一样，叫你来开龙不来，sb"
	//
	//// 构造正则匹配字符
	//regStr := strings.Join(sensitiveWords, "|")
	//println("regStr -> ", regStr)
	//傻逼|傻叉|垃圾|妈的|sb
	//wordReg := regexp.MustCompile(regStr)
	//text = wordReg.ReplaceAllString(text, "*")
	//
	//println("text -> ", text)
	//func ReadFile(name string) ([]byte, error) {}
	content, err := os.ReadFile(SensitiveWordFilePath)
	if err != nil {
		LogError(err.Error())
		panic(err)
	}
	s := string(content)
	wordReg = regexp.MustCompile(s)
}

func SensitiveWordsFilter(context string) (string, error) {
	test := wordReg.ReplaceAllString(context, "*")
	fmt.Printf("%s", test)
	return test, nil
}

func Test() {
	sensitiveWords := []string{
		"傻逼",
		"傻叉",
		"垃圾",
		"妈的",
		"sb",
	}
	text := "什么垃圾打野，傻逼一样，叫你来开龙不来，sb"
	// 构造正则匹配字符
	regStr := strings.Join(sensitiveWords, "|")
	println("regStr -> ", regStr)
	wordReg := regexp.MustCompile(regStr)
	text = wordReg.ReplaceAllString(text, "*")

	println("text -> ", text)

}
