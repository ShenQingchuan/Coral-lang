package test

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/lexer"
	. "coral-lang/src/parser"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseLiteral(t *testing.T) {
	Convey("测试解析字面量值：十进制", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "8997")
		So(parser.CurrentToken.Str, ShouldEqual, "8997")

		a := parser.ParseLiteral()
		_, isDecimalLit := a.(*DecimalLit)
		So(isDecimalLit, ShouldEqual, true)
	})

	Convey("测试解析字面量值：浮点数", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "3.11")
		So(parser.CurrentToken.Str, ShouldEqual, "3.11")

		a := parser.ParseLiteral()
		_, isFloatLit := a.(*FloatLit)
		So(isFloatLit, ShouldEqual, true)
	})

	Convey("测试解析字面量值：含e科学记数法的指数", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "3.63e-2")
		So(parser.CurrentToken.Str, ShouldEqual, "3.63e-2")

		a := parser.ParseLiteral()
		_, isExponentLit := a.(*ExponentLit)
		So(isExponentLit, ShouldEqual, true)
	})
}
func TestParseOperandName(t *testing.T) {
	Convey("测试解析操作数名称：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "httplib.HttpRequest")
		So(parser.CurrentToken.Str, ShouldEqual, "httplib")

		operandName := parser.ParseOperandName()
		So(operandName.GetFullName(), ShouldEqual, "httplib.HttpRequest")
	})
}
func TestBinaryExpression(t *testing.T) {
	Convey("测试二元表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "num* 3+ 0x3f21")
		So(parser.CurrentToken.Str, ShouldEqual, "num")

		binaryExpression, isBinary := parser.ParseExpression().(*BinaryExpression)
		So(isBinary, ShouldEqual, true) // 证明右边也是一颗二元表达式节点

		So(binaryExpression.Operator.Kind, ShouldEqual, TokenTypePlus)

		num := binaryExpression.Left.(*BinaryExpression).Left.(*BasicPrimaryExpression).Operand.(*OperandName)
		So(num.GetFullName(), ShouldEqual, "num")

		three := binaryExpression.Left.(*BinaryExpression).Right.(*BasicPrimaryExpression).Operand.(*DecimalLit)
		So(three.Value.Kind, ShouldEqual, TokenTypeDecimalInteger)
		So(three.Value.Str, ShouldEqual, "3")

		hex := binaryExpression.Right.(*BasicPrimaryExpression).Operand.(*HexadecimalLit)
		So(hex.Value.Kind, ShouldEqual, TokenTypeHexadecimalInteger)
		So(hex.Value.Str, ShouldEqual, "0x3f21")
	})
}
