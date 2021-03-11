package test

import (
	. "coral-lang/src/parser"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMultipleErrorsAndWarnings(t *testing.T) {
	Convey("测试允许编译集合多个错误/警告：", t, func() {
		parser := new(Parser)
		parser.InitFromString(`fn add(x, y) {
  println("sum: ", x+y)
}`)
		parser.ParseProgram()
	})
}
