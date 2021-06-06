package test

import (
	. "coral-lang/src/utils"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUnicodeToUTF8(t *testing.T) {
	Convey("Unicode 转 UTF8", t, func() {
		So(UnicodeToUTF8("\\u77e5", 4), ShouldEqual, "知")
		So(UnicodeToUTF8("\\u94F8", 4), ShouldEqual, "铸")
		So(UnicodeToUTF8("72fc", 4), ShouldEqual, "狼")
		So(UnicodeToUTF8("21", 2), ShouldEqual, "!")
	})
}

func TestColoredPrint(t *testing.T) {
	fmt.Println(Red("红色字符串"))
	fmt.Println(Yellow("黄色字符串"))
	fmt.Println(Blue("蓝色字符串"))
	fmt.Println(Green("绿色字符串"))
	fmt.Println(Purple("紫色字符串"))
	fmt.Println(Gray("灰色字符串"))
}
