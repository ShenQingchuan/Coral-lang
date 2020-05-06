package test

import (
	. "coral-lang/src/exception"
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
	parser.InitLexerFromBytes(testLexer, content)

	Convey("测试 PeekChar：", t, func() {
		firstUtf8Char := testLexer.PeekChar()
		So(firstUtf8Char.Rune, ShouldEqual, 'i')
	})
}
func TestPeekNextChar(t *testing.T) {
	content := parser.OpenSourceFile("samples/test.coral")
	testLexer := &parser.Lexer{}
	parser.InitLexerFromBytes(testLexer, content)

	Convey("测试 PeekNextChar：", t, func() {
		firstUtf8Char := testLexer.PeekChar()
		secondUtf8Char := testLexer.PeekNextChar(firstUtf8Char.ByteLength)
		So(secondUtf8Char.Rune, ShouldEqual, 'm')
	})
}
func TestPeekNextCharByStep(t *testing.T) {
	content := parser.OpenSourceFile("samples/test.coral")
	testLexer := &parser.Lexer{}
	parser.InitLexerFromBytes(testLexer, content)

	Convey("测试 PeekNextCharByStep：", t, func() {
		firstUtf8Char := testLexer.PeekChar()
		fourthUtf8Char := testLexer.PeekNextCharByStep(firstUtf8Char.ByteLength, 3)
		So(fourthUtf8Char.Rune, ShouldEqual, 'o')
	})
}
func TestGoNextChar(t *testing.T) {
	content := parser.OpenSourceFile("samples/test.coral")
	testLexer := &parser.Lexer{}
	parser.InitLexerFromBytes(testLexer, content)

	Convey("测试 GoNextChar", t, func() {
		testLexer.GoNextCharByStep(3)
		So(testLexer.PeekChar().Rune, ShouldEqual, 'o')
	})
}
func TestReadDecimal(t *testing.T) {
	testLexer := &parser.Lexer{}
	parser.InitLexerFromString(testLexer, "386")

	Convey("测试读入十进制整数", t, func() {
		gotToken, err := testLexer.ReadDecimal(false)
		if err != nil {
			panic(err)
		}
		So(gotToken.Str, ShouldEqual, "386")
		So(gotToken.Kind, ShouldEqual, parser.TokenTypeDecimalInteger)
		fmt.Println(gotToken)
	})
}
func TestReadFloat(t *testing.T) {
	successLexer := &parser.Lexer{}
	parser.InitLexerFromString(successLexer, "3.5681")

	Convey("读入小数: success", t, func() {
		gotToken, err := successLexer.ReadDecimal(false)
		if err != nil {
			CoralErrorHandler(err)
		} else {
			So(gotToken.Str, ShouldEqual, "3.5681")
			So(gotToken.Kind, ShouldEqual, parser.TokenTypeFloat)
		}
	})

	failLexer := &parser.Lexer{}
	parser.InitLexerFromString(failLexer, "3.56.81")
	Convey("读入小数: error", t, func() {
		_, err := failLexer.ReadDecimal(false)
		So(err.ErrEnum, ShouldEqual, LexFloatFormatError)
	})
}
func TestReadExponent(t *testing.T) {
	successLexer := &parser.Lexer{}
	parser.InitLexerFromString(successLexer, "1.7e+2")

	Convey("读入科学记数法: success", t, func() {
		gotToken, err := successLexer.ReadDecimal(false)
		if err != nil {
			CoralErrorHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "1.7e+2")
		So(gotToken.Kind, ShouldEqual, parser.TokenTypeExponent)
	})

	failLexer := &parser.Lexer{}
	parser.InitLexerFromString(failLexer, "5.4e-2e08")
	Convey("读入科学记数法: error", t, func() {
		_, err := failLexer.ReadDecimal(false)
		So(err.ErrEnum, ShouldEqual, LexExponentFormatError)
	})
}
func TestReadHexadecimal(t *testing.T) {
	testLexer := &parser.Lexer{}
	parser.InitLexerFromString(testLexer, "0xEf012a")

	Convey("测试读入十六进制整数", t, func() {
		gotToken, err := testLexer.ReadHexadecimal()
		if err != nil {
			CoralErrorHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "0xEf012a")
		So(gotToken.Kind, ShouldEqual, parser.TokenTypeHexadecimalInteger)
	})
}
func TestReadBinary(t *testing.T) {
	testLexer := &parser.Lexer{}
	parser.InitLexerFromString(testLexer, "0b101001")

	Convey("测试读入二进制整数", t, func() {
		gotToken, err := testLexer.ReadBinary()
		if err != nil {
			CoralErrorHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "0b101001")
		So(gotToken.Kind, ShouldEqual, parser.TokenTypeBinaryInteger)
	})
}
func TestReadOctal(t *testing.T) {
	testLexer := &parser.Lexer{}
	parser.InitLexerFromString(testLexer, "0o1073")

	Convey("测试读入八进制整数", t, func() {
		gotToken, err := testLexer.ReadOctal()
		if err != nil {
			CoralErrorHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "0o1073")
		So(gotToken.Kind, ShouldEqual, parser.TokenTypeOctalInteger)
	})
}
