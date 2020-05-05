package test

import (
	"coral-lang/src/parser"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

/*
 !! [Important Notice] !!:
	Please run these unit tests separately,
	because some of them contain the opening of 'samples/test.coral' file.
	It's in order to avoid meaningless memory allocation.
*/

func TestPeekChar(t *testing.T) {
	content := parser.OpenSourceFile("samples/test.coral")
	testLexer := &parser.Lexer{}
	parser.InitLexer(testLexer, content)

	Convey("测试 PeekChar：", t, func() {
		firstUtf8Char := testLexer.PeekChar()
		So(firstUtf8Char.Rune, ShouldEqual, 'i')
	})
}
func TestPeekNextChar(t *testing.T) {
	content := parser.OpenSourceFile("samples/test.coral")
	testLexer := &parser.Lexer{}
	parser.InitLexer(testLexer, content)

	Convey("测试 PeekNextChar：", t, func() {
		firstUtf8Char := testLexer.PeekChar()
		secondUtf8Char := testLexer.PeekNextChar(firstUtf8Char.ByteLength)
		So(secondUtf8Char.Rune, ShouldEqual, 'm')
	})
}
func TestPeekNextCharByStep(t *testing.T) {
	content := parser.OpenSourceFile("samples/test.coral")
	testLexer := &parser.Lexer{}
	parser.InitLexer(testLexer, content)

	Convey("测试 PeekNextCharByStep：", t, func() {
		firstUtf8Char := testLexer.PeekChar()
		fourthUtf8Char := testLexer.PeekNextCharByStep(firstUtf8Char.ByteLength, 3)
		So(fourthUtf8Char.Rune, ShouldEqual, 'o')
	})
}
func TestGoNextChar(t *testing.T) {
	content := parser.OpenSourceFile("samples/test.coral")
	testLexer := &parser.Lexer{}
	parser.InitLexer(testLexer, content)

	Convey("测试 GoNextChar", t, func() {
		testLexer.GoNextCharByStep(3)
		So(testLexer.PeekChar().Rune, ShouldEqual, 'o')
	})
}
func TestReadDecimal(t *testing.T) {
	testLexer := &parser.Lexer{}
	parser.InitLexer(testLexer, []byte("386"))

	Convey("测试读入十进制整数", t, func() {
		gotToken := testLexer.ReadDecimal(false)
		So(gotToken.Str, ShouldEqual, "386")
		So(gotToken.Kind, ShouldEqual, parser.TokenTypeDecimalInteger)
		fmt.Println(gotToken)
	})
}
func TestReadHexadecimal(t *testing.T) {
	testLexer := &parser.Lexer{}
	parser.InitLexer(testLexer, []byte("0xEf012a"))

	Convey("测试读入十六进制整数", t, func() {
		gotToken := testLexer.ReadHexadecimal()
		So(gotToken.Str, ShouldEqual, "0xEf012a")
		So(gotToken.Kind, ShouldEqual, parser.TokenTypeHexadecimalInteger)
	})
}
