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
func TestBinaryExpression(t *testing.T) {
	Convey("测试二元表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "num* 3+ 0x3f21")
		So(parser.CurrentToken.Str, ShouldEqual, "num")

		binaryExpression, isBinary := parser.ParseExpression().(*BinaryExpression)
		So(isBinary, ShouldEqual, true)

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
func TestIndexSliceCallMemberExpression(t *testing.T) {
	Convey("测试索引表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "arr[i]")
		So(parser.CurrentToken.Str, ShouldEqual, "arr")

		indexExpression, isIndex := parser.ParseExpression().(*IndexExpression)
		So(isIndex, ShouldEqual, true)

		So(indexExpression.Operand.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "arr")
		So(indexExpression.Index.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "i")
	})

	Convey("测试切片表达式 1：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "arr[:4]")
		So(parser.CurrentToken.Str, ShouldEqual, "arr")

		sliceExpression, isSlice := parser.ParseExpression().(*SliceExpression)
		So(isSlice, ShouldEqual, true)

		So(sliceExpression.Operand.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "arr")
		So(sliceExpression.End.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str,
			ShouldEqual, "4")
	})

	Convey("测试切片表达式 2：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "arr[kk:5]")
		So(parser.CurrentToken.Str, ShouldEqual, "arr")

		sliceExpression, isSlice := parser.ParseExpression().(*SliceExpression)
		So(isSlice, ShouldEqual, true)

		So(sliceExpression.Operand.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "arr")
		So(sliceExpression.Start.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "kk")
		So(sliceExpression.End.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str,
			ShouldEqual, "5")
	})

	Convey("测试函数调用表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "funcA(i, 3.11)")
		So(parser.CurrentToken.Str, ShouldEqual, "funcA")

		callExpression, isCall := parser.ParseExpression().(*CallExpression)
		So(isCall, ShouldEqual, true)

		So(callExpression.Operand.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "funcA")
		So(callExpression.Params[0].(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "i")
		So(callExpression.Params[1].(*BasicPrimaryExpression).Operand.(*FloatLit).Value.Str,
			ShouldEqual, "3.11")
	})

	Convey("测试成员访问表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "request.query.page")
		So(parser.CurrentToken.Str, ShouldEqual, "request")

		memberExpression, isMemberExpr := parser.ParseExpression().(*MemberExpression)
		So(isMemberExpr, ShouldEqual, true)

		So(memberExpression.Operand.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "request")
		So(memberExpression.Member.Operand.Token.Str, ShouldEqual, "query")
		So(memberExpression.Member.MemberNext.Operand.Token.Str, ShouldEqual, "page")
	})

	Convey("混合测试", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "makeArr(dd +7, 3.11e2)[mm:6]")
		So(parser.CurrentToken.Str, ShouldEqual, "makeArr")

		sliceExpression, isSlice := parser.ParseExpression().(*SliceExpression)
		So(isSlice, ShouldEqual, true)

		callExpression, isIndexOperandCall := sliceExpression.Operand.(*CallExpression)
		So(isIndexOperandCall, ShouldEqual, true)

		So(callExpression.Operand.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "makeArr")
		So(callExpression.Params[0].(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypePlus)
		So(callExpression.Params[0].(*BinaryExpression).Left.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "dd")
		So(callExpression.Params[0].(*BinaryExpression).Right.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str,
			ShouldEqual, "7")

		So(sliceExpression.Start.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "mm")
		So(sliceExpression.End.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str,
			ShouldEqual, "6")
	})
}
func TestNewInstanceExpression(t *testing.T) {
	Convey("测试对象实例新建表达式：无泛型参数", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `new Student("Peter", 18)`)
		So(parser.CurrentToken.Str, ShouldEqual, "new")

		newInstanceExpression, isNewInstance := parser.ParseExpression().(*NewInstanceExpression)
		So(isNewInstance, ShouldEqual, true)
		So(newInstanceExpression, ShouldNotEqual, nil)
		So(newInstanceExpression.Class.(*TypeName).GetFullName(),
			ShouldEqual, "Student")
		So(newInstanceExpression.InitParams[0].(*BasicPrimaryExpression).Operand.(*StringLit).Value.Str,
			ShouldEqual, "Peter")
		So(newInstanceExpression.InitParams[1].(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str,
			ShouldEqual, "18")
	})

	Convey("测试对象实例新建表达式：含泛型参数", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "new Array<string>(3)")
		So(parser.CurrentToken.Str, ShouldEqual, "new")

		newInstanceExpression, isNewInstance := parser.ParseExpression().(*NewInstanceExpression)
		So(isNewInstance, ShouldEqual, true)
		So(newInstanceExpression, ShouldNotEqual, nil)
		So(newInstanceExpression.Class.(*GenericsTypeLit).BasicType.GetFullName(),
			ShouldEqual, "Array")
		So(newInstanceExpression.Class.(*GenericsTypeLit).GenericsArgs[0].GetFullName(),
			ShouldEqual, "string")
		So(newInstanceExpression.InitParams[0].(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str,
			ShouldEqual, "3")
	})
}
func TestAssignListStatement(t *testing.T) {
	Convey("测试赋值列表：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `num_a, x[1]= 3, "hello";`)
		So(parser.CurrentToken.Str, ShouldEqual, "num_a")

		assignListStmt, isAssignList := parser.ParseStatement().(*AssignListStatement)
		So(isAssignList, ShouldEqual, true)

		So(assignListStmt.Targets[0].(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "num_a")
		So(assignListStmt.Targets[1].(*IndexExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "x")
		So(assignListStmt.Targets[1].(*IndexExpression).Index.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str,
			ShouldEqual, "1")
		So(assignListStmt.Values[0].(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str,
			ShouldEqual, "3")
		So(assignListStmt.Values[1].(*BasicPrimaryExpression).Operand.(*StringLit).Value.Str,
			ShouldEqual, "hello")
	})
}
