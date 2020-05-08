package test

import (
	. "coral-lang/src/utils"
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
