package test

import (
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
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
	content := OpenSourceFile("samples/test.coral")
	testLexer := &Lexer{}
	testLexer.InitFromBytes(content)

	Convey("测试 PeekChar：", t, func() {
		firstUtf8Char := testLexer.PeekChar()
		So(firstUtf8Char.Rune, ShouldEqual, 'i')
	})
}
func TestPeekNextChar(t *testing.T) {
	content := OpenSourceFile("samples/test.coral")
	testLexer := &Lexer{}
	testLexer.InitFromBytes(content)

	Convey("测试 PeekNextChar：", t, func() {
		firstUtf8Char := testLexer.PeekChar()
		secondUtf8Char := testLexer.PeekNextChar(firstUtf8Char.ByteLength)
		So(secondUtf8Char.Rune, ShouldEqual, 'm')
	})
}
func TestPeekNextCharByStep(t *testing.T) {
	content := OpenSourceFile("samples/test.coral")
	testLexer := &Lexer{}
	testLexer.InitFromBytes(content)

	Convey("测试 PeekNextCharByStep：", t, func() {
		firstUtf8Char := testLexer.PeekChar()
		fourthUtf8Char := testLexer.PeekNextCharByStep(firstUtf8Char.ByteLength, 3)
		So(fourthUtf8Char.Rune, ShouldEqual, 'o')
	})
}
func TestGoNextChar(t *testing.T) {
	content := OpenSourceFile("samples/test.coral")
	testLexer := &Lexer{}
	testLexer.InitFromBytes(content)

	Convey("测试 GoNextChar", t, func() {
		testLexer.GoNextCharByStep(3)
		So(testLexer.PeekChar().Rune, ShouldEqual, 'o')
	})
}
func TestReadDecimal(t *testing.T) {
	testLexer1 := &Lexer{}
	testLexer1.InitFromString("386")

	Convey("测试读入十进制整数：normal", t, func() {
		gotToken, err := testLexer1.ReadDecimal(false)
		if err != nil {
			panic(err)
		}
		So(gotToken.Str, ShouldEqual, "386")
		So(gotToken.Kind, ShouldEqual, TokenTypeDecimalInteger)
	})

	testLexer2 := &Lexer{}
	testLexer2.InitFromString("000186")

	Convey("测试读入十进制整数：more zero", t, func() {
		gotToken, err := testLexer2.ReadDecimal(true)
		if err != nil {
			panic(err)
		}
		So(gotToken.Str, ShouldEqual, "186")
		So(gotToken.Kind, ShouldEqual, TokenTypeDecimalInteger)
	})
}
func TestReadFloat(t *testing.T) {
	successLexer := &Lexer{}
	successLexer.InitFromString("3.5681")

	Convey("读入小数: success", t, func() {
		gotToken, err := successLexer.ReadDecimal(false)
		if err != nil {
			CoralErrorCrashHandler(err)
		} else {
			So(gotToken.Str, ShouldEqual, "3.5681")
			So(gotToken.Kind, ShouldEqual, TokenTypeFloat)
		}
	})

	testLexer := &Lexer{}
	testLexer.InitFromString("000.186")

	Convey("测试读入零开头的小数：", t, func() {
		gotToken, err := testLexer.ReadDecimal(true)
		if err != nil {
			panic(err)
		}
		So(gotToken.Str, ShouldEqual, "0.186")
		So(gotToken.Kind, ShouldEqual, TokenTypeFloat)
	})

	failLexer := &Lexer{}
	failLexer.InitFromString("3.56.81")
	Convey("读入小数: error", t, func() {
		_, err := failLexer.ReadDecimal(true)
		So(err.ErrEnum, ShouldEqual, LexFloatFormatError)
	})
}
func TestReadExponent(t *testing.T) {
	successLexer := &Lexer{}
	successLexer.InitFromString("1.7e+2")

	Convey("读入科学记数法: success", t, func() {
		gotToken, err := successLexer.ReadDecimal(false)
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "1.7e+2")
		So(gotToken.Kind, ShouldEqual, TokenTypeExponent)
	})

	failLexer := &Lexer{}
	failLexer.InitFromString("5.4e-2e08")
	Convey("读入科学记数法: error", t, func() {
		_, err := failLexer.ReadDecimal(false)
		So(err.ErrEnum, ShouldEqual, LexExponentFormatError)
	})

	failLexer2 := &Lexer{}
	failLexer2.InitFromString("0e-5")
	Convey("无意义的 0e 开头：", t, func() {
		_, err := failLexer2.ReadDecimal(true)
		So(err.ErrEnum, ShouldEqual, LexExponentFormatError)
	})
}
func TestReadHexadecimal(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString("0xEf012a")

	Convey("测试读入十六进制整数", t, func() {
		gotToken, err := testLexer.ReadHexadecimal()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "0xEf012a")
		So(gotToken.Kind, ShouldEqual, TokenTypeHexadecimalInteger)
	})
}
func TestReadBinary(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString("0b101001")

	Convey("测试读入二进制整数", t, func() {
		gotToken, err := testLexer.ReadBinary()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "0b101001")
		So(gotToken.Kind, ShouldEqual, TokenTypeBinaryInteger)
	})
}
func TestReadOctal(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString("0o1073")

	Convey("测试读入八进制整数", t, func() {
		gotToken, err := testLexer.ReadOctal()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "0o1073")
		So(gotToken.Kind, ShouldEqual, TokenTypeOctalInteger)
	})
}
func TestReadString(t *testing.T) {
	testLexer1 := &Lexer{}
	testLexer1.InitFromString("\"我就是\\t想装个逼：\\u77e5道unicode是这样的\"")

	Convey("测试读入字符串（支持转义 \\u 字符）", t, func() {
		gotToken, err := testLexer1.ReadString()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "我就是\t想装个逼：知道unicode是这样的")
		So(gotToken.Kind, ShouldEqual, TokenTypeString)
	})

	testLexer2 := &Lexer{}
	testLexer2.InitFromString("\"来个单的：\\xD688\"")

	Convey("测试读入字符串（支持转义 \\x 字符）", t, func() {
		gotToken, err := testLexer2.ReadString()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "来个单的：Ö88")
		So(gotToken.Kind, ShouldEqual, TokenTypeString)
	})

	failLexer := &Lexer{}
	failLexer.InitFromString("\"\\u332\"")
	Convey("测试读入unicode编码但不足4位报错: error", t, func() {
		_, err := failLexer.ReadString()
		So(err.ErrEnum, ShouldEqual, LexUnicodeEscapeFormatError)
	})
}
func TestReadRune(t *testing.T) {
	Convey("测试读入字符 1", t, func() {
		testLexer := &Lexer{}
		testLexer.InitFromString("'Z'")
		gotToken, err := testLexer.ReadRune()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "Z")
		So(gotToken.Kind, ShouldEqual, TokenTypeRune)
	})

	Convey("测试读入字符 2：支持转义字符", t, func() {
		testLexer := &Lexer{}
		testLexer.InitFromString("'\\u94F8'")
		gotToken, err := testLexer.ReadRune()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "铸")
		So(gotToken.Kind, ShouldEqual, TokenTypeRune)
	})
}
func TestReadIdentifierAscii(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString("num_x")

	Convey("测试读入标识符: ascii", t, func() {
		gotToken, err := testLexer.ReadIdentifier()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "num_x")
		So(gotToken.Kind, ShouldEqual, TokenTypeIdentifier)
	})
}
func TestReadIdentifierUTF8(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString("大π∆变量1")

	Convey("测试读入标识符: UTF8", t, func() {
		gotToken, err := testLexer.ReadIdentifier()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "大π∆变量1")
		So(gotToken.Kind, ShouldEqual, TokenTypeIdentifier)
	})
}
func TestReadIdentifierButItIsKeyword(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString("var")

	Convey("测试读入关键字：", t, func() {
		gotToken, err := testLexer.ReadIdentifier()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		So(gotToken.Str, ShouldEqual, "var")
		So(gotToken.Kind, ShouldEqual, TokenTypeVar)
	})
}
func TestReadIdentifierThrowError(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString("2un√π")
	Convey("测试读入标识符: error", t, func() {
		_, err := testLexer.ReadIdentifier()
		So(err.ErrEnum, ShouldEqual, LexIdentifierFirstRuneCanNotBeDigit)
	})
}
func TestSkipLineComment(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString(`
	// this is a line comment
  val
	`)

	gotToken, err := testLexer.GetNextToken(false)
	if err != nil {
		CoralErrorCrashHandler(err)
	}
	Convey("测试 跳过行注释：", t, func() {
		So(gotToken.Str, ShouldEqual, "val")
		So(gotToken.Kind, ShouldEqual, TokenTypeVal)
	})
}
func TestSkipBlockComment(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString(`
	/* this is a block comment
     line two
     line three, oh my god it works !
  */
  import
	`)

	gotToken, err := testLexer.GetNextToken(false)
	if err != nil {
		CoralErrorCrashHandler(err)
	}
	Convey("测试 跳过块注释：", t, func() {
		So(gotToken.Str, ShouldEqual, "import")
		So(gotToken.Kind, ShouldEqual, TokenTypeImport)
	})
}
func TestUnclosedException(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString(`(([{])}`)

	Convey("测试括号未匹配报错：", t, func() {
		count := 0
		for _, err := testLexer.GetNextToken(false); count < 8; _, err = testLexer.GetNextToken(false) {
			if count == 7 {
				// 此时已经到达 BytePos 的末尾，但是圆括号仍未关闭完全
				So(err.ErrEnum, ShouldEqual, LexParenthesesUnclosed)
			}
			count++
		}
	})
}

// 大测试：给一段源代码来生成 token 流
func TestGetNextToken(t *testing.T) {
	testLexer := &Lexer{}
	testLexer.InitFromString(`
	from httplib import {
	 	HttpRequest as req
  }

  var i int = 1 *3 <<= 5, e = esm(6 |= 0.2);
  val b double = functionA() ** 6.3 -1e-5;
	`)

	expectedTokens := []struct {
		expectedKind TokenType
		expectedStr  string
	}{
		{TokenTypeFrom, "from"},
		{TokenTypeIdentifier, "httplib"},
		{TokenTypeImport, "import"},
		{TokenTypeLeftBrace, "{"},
		{TokenTypeIdentifier, "HttpRequest"},
		{TokenTypeAs, "as"},
		{TokenTypeIdentifier, "req"},
		{TokenTypeRightBrace, "}"},
		{TokenTypeVar, "var"},
		{TokenTypeIdentifier, "i"},
		{TokenTypeIdentifier, "int"},
		{TokenTypeEqual, "="},
		{TokenTypeDecimalInteger, "1"},
		{TokenTypeStar, "*"},
		{TokenTypeDecimalInteger, "3"},
		{TokenTypeDoubleLeftAngleEqual, "<<="},
		{TokenTypeDecimalInteger, "5"},
		{TokenTypeComma, ","},
		{TokenTypeIdentifier, "e"},
		{TokenTypeEqual, "="},
		{TokenTypeIdentifier, "esm"},
		{TokenTypeLeftParen, "("},
		{TokenTypeDecimalInteger, "6"},
		{TokenTypeVerticalEqual, "|="},
		{TokenTypeFloat, "0.2"},
		{TokenTypeRightParen, ")"},
		{TokenTypeSemi, ";"},
		{TokenTypeVal, "val"},
		{TokenTypeIdentifier, "b"},
		{TokenTypeIdentifier, "double"},
		{TokenTypeEqual, "="},
		{TokenTypeIdentifier, "functionA"},
		{TokenTypeLeftParen, "("},
		{TokenTypeRightParen, ")"},
		{TokenTypeDoubleStar, "**"},
		{TokenTypeFloat, "6.3"},
		{TokenTypeMinus, "-"},
		{TokenTypeExponent, "1e-5"},
		{TokenTypeSemi, ";"},
	}

	Convey("测试读入 GetToken 流：", t, func() {
		for _, t := range expectedTokens {
			gotToken, err := testLexer.GetNextToken(false)
			if err != nil {
				CoralErrorCrashHandler(err)
			}
			So(gotToken.Str, ShouldEqual, t.expectedStr)
			So(gotToken.Kind, ShouldEqual, t.expectedKind)
		}
	})
}
