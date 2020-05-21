package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type Any interface{}

// !! 重点学习 有助于理解编码的内涵
func UnicodeToUTF8(source string, bit int) string {
	res := []string{""}
	// 无论你是 \uXXXX 还是没有前缀的 XXXX 都可
	sUnicode := strings.Split(source, "\\u") // 切分出四位 unicode 的编码
	context := ""
	for _, v := range sUnicode {
		additional := ""
		if len(v) < 1 {
			continue
		}
		if len(v) > bit { // 长度大于 4
			rs := []rune(v)
			v = string(rs[:bit])          // 切分出前 4 位
			additional = string(rs[bit:]) // 后面的当作正常字符处理
		}
		temp, err := strconv.ParseInt(v, 16, bit*8) // <- 32bit 即允许四个字节整数
		if err != nil {
			context += v
		}
		context += fmt.Sprintf("%c", temp) // 使用 Go 原生支持的 UTF-8 的 %c 来输出该字符
		context += additional              // 添加上余下多出的一些正常字符
	}
	res = append(res, context)
	return strings.Join(res, "")
}
