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
func TestRangeExpression(t *testing.T) {
	Convey("测试区间表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "0..arr.length")
		So(parser.CurrentToken.Str, ShouldEqual, "0")

		rangeExpression, isRange := parser.ParseExpression().(*RangeExpression)
		So(isRange, ShouldEqual, true)
		So(rangeExpression.IncludeEnd, ShouldEqual, false)

		So(rangeExpression.Start.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str, ShouldEqual, "0")
		So(rangeExpression.End.(*MemberExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"arr")
		So(rangeExpression.End.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual, "length")
		So(rangeExpression.End.(*MemberExpression).Member.MemberNext.Operand, ShouldEqual, nil)
		So(rangeExpression.End.(*MemberExpression).Member.MemberNext.MemberNext, ShouldEqual, nil)
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
func TestVarValDeclarationStatement(t *testing.T) {
	Convey("测试变量定义：1", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "var a int = 6;")
		So(parser.CurrentToken.Str, ShouldEqual, "var")

		varDeclStatement, isVarDecl := parser.ParseStatement().(*VarDeclStatement)
		So(isVarDecl, ShouldEqual, true)

		So(varDeclStatement.Mutable, ShouldEqual, true)
		So(varDeclStatement.Declarations[0].VarName.Str, ShouldEqual,
			"a")
		So(varDeclStatement.Declarations[0].Type.(*TypeName).GetFullName(), ShouldEqual,
			"int")
		So(varDeclStatement.Declarations[0].InitValue.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str, ShouldEqual,
			"6")
	})

	Convey("测试变量定义：2", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "val 圆周率 = 3.14, 光速 = 3e8;")
		So(parser.CurrentToken.Str, ShouldEqual, "val")

		varDeclStatement, isVarDecl := parser.ParseStatement().(*VarDeclStatement)
		So(isVarDecl, ShouldEqual, true)

		So(varDeclStatement.Mutable, ShouldEqual, false)
		So(varDeclStatement.Declarations[0].VarName.Str, ShouldEqual,
			"圆周率")
		So(varDeclStatement.Declarations[0].Type, ShouldEqual, nil)
		So(varDeclStatement.Declarations[0].InitValue.(*BasicPrimaryExpression).Operand.(*FloatLit).Value.Str, ShouldEqual,
			"3.14")

		So(varDeclStatement.Declarations[1].VarName.Str, ShouldEqual,
			"光速")
		So(varDeclStatement.Declarations[1].Type, ShouldEqual, nil)
		So(varDeclStatement.Declarations[1].InitValue.(*BasicPrimaryExpression).Operand.(*ExponentLit).Value.Str, ShouldEqual,
			"3e8")
	})
}
func TestUnaryExpression(t *testing.T) {
	Convey("测试单目运算符解析：1", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "!m.tt&&~ss.k[1]")
		So(parser.CurrentToken.Str, ShouldEqual, "!")

		binaryExpression, isBinary := parser.ParseExpression().(*BinaryExpression)
		So(isBinary, ShouldEqual, true)
		So(binaryExpression.Operator.Kind, ShouldEqual, TokenTypeDoubleAmpersand)

		leftUnary, isLeftUnary := binaryExpression.Left.(*UnaryExpression)
		So(isLeftUnary, ShouldEqual, true)
		So(leftUnary.Operator.Kind, ShouldEqual, TokenTypeBang)
		So(leftUnary.Operand.(*MemberExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "m")
		So(leftUnary.Operand.(*MemberExpression).Member.Operand.Token.Str,
			ShouldEqual, "tt")

		rightUnary, isRightUnary := binaryExpression.Right.(*UnaryExpression)
		So(isRightUnary, ShouldEqual, true)
		So(rightUnary.Operator.Kind, ShouldEqual, TokenTypeWavy)
		So(rightUnary.Operand.(*IndexExpression).Operand.(*MemberExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).GetFullName(),
			ShouldEqual, "ss")
		So(rightUnary.Operand.(*IndexExpression).Operand.(*MemberExpression).Member.Operand.Token.Str,
			ShouldEqual, "k")
		So(rightUnary.Operand.(*IndexExpression).Index.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str,
			ShouldEqual, "1")
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
func TestImportStatement(t *testing.T) {
	Convey("测试导入模块语句：1", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `from SeaCoral import Request as Req;`)
		So(parser.CurrentToken.Str, ShouldEqual, "from")

		stmt := parser.ParseStatement()
		singleImportStatement, isSingleImport := stmt.(*SingleImportStatement)
		So(isSingleImport, ShouldEqual, true)
		So(singleImportStatement.From.GetFullModuleName(), ShouldEqual, "SeaCoral")
		So(singleImportStatement.Element.ModuleName.GetFullModuleName(), ShouldEqual, "Request")
		So(singleImportStatement.Element.As.Token.Str, ShouldEqual, "Req")
	})

	Convey("测试导入模块语句：2", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `from SeaCoral import {
      Request as Req,
			Response as Resp
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "from")

		listImportStatement, isListImport := parser.ParseStatement().(*ListImportStatement)
		So(isListImport, ShouldEqual, true)
		So(listImportStatement.From.GetFullModuleName(), ShouldEqual, "SeaCoral")
		So(listImportStatement.Elements[0].ModuleName.GetFullModuleName(), ShouldEqual, "Request")
		So(listImportStatement.Elements[0].As.Token.Str, ShouldEqual, "Req")
		So(listImportStatement.Elements[1].ModuleName.GetFullModuleName(), ShouldEqual, "Response")
		So(listImportStatement.Elements[1].As.Token.Str, ShouldEqual, "Resp")
	})
}
func TestEnumStatement(t *testing.T) {
	Convey("测试枚举定义语句解析：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `enum Sex {
        FEMALE = 0,
        MALE,
        SECRET
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "enum")

		enumStatement, isEnum := parser.ParseStatement().(*EnumStatement)
		So(isEnum, ShouldEqual, true)
		So(enumStatement.Name.Token.Str, ShouldEqual, "Sex")
		So(enumStatement.Elements[0].Name.Token.Str, ShouldEqual, "FEMALE")
		So(enumStatement.Elements[0].Value.Value.Str, ShouldEqual, "0")
		So(enumStatement.Elements[1].Name.Token.Str, ShouldEqual, "MALE")
		So(enumStatement.Elements[2].Name.Token.Str, ShouldEqual, "SECRET")
	})
}
func TestIfStatement(t *testing.T) {
	Convey("测试条件语句解析：1", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `if !screen.closed {
			println("屏幕还没关！");
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "if")

		ifStatement, isIfStmt := parser.ParseStatement().(*IfStatement)
		So(isIfStmt, ShouldEqual, true)

		So(ifStatement.If.Condition.(*UnaryExpression).Operator.Kind, ShouldEqual, TokenTypeBang)
		So(ifStatement.If.Condition.(*UnaryExpression).Operand.(*MemberExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"screen")
		So(ifStatement.If.Condition.(*UnaryExpression).Operand.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"closed")
		So(ifStatement.If.Condition.(*UnaryExpression).Operand.(*MemberExpression).Member.MemberNext.Operand, ShouldEqual,
			nil)
		So(ifStatement.If.Condition.(*UnaryExpression).Operand.(*MemberExpression).Member.MemberNext.MemberNext, ShouldEqual,
			nil)
		So(ifStatement.If.Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(ifStatement.If.Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).Operand.(*StringLit).Value.Str, ShouldEqual,
			"屏幕还没关！")
	})

	Convey("测试条件语句解析：2", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `if t.getYear() == 2020 {
			wawa++;
    } elif m.what > 3.55 {
			x, y = 1.3e6, 'Z';
    } else {
			bb = 6 + 3 * dd;
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "if")

		ifStatement, isIfStmt := parser.ParseStatement().(*IfStatement)
		So(isIfStmt, ShouldEqual, true)

		So(ifStatement.If.Condition.(*BinaryExpression).Left.(*CallExpression).Operand.(*MemberExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"t")
		So(ifStatement.If.Condition.(*BinaryExpression).Left.(*CallExpression).Operand.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"getYear")
		So(len(ifStatement.If.Condition.(*BinaryExpression).Left.(*CallExpression).Params), ShouldEqual,
			0)
		So(ifStatement.If.Condition.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeDoubleEqual)
		So(ifStatement.If.Block.Statements[0].(*IncDecStatement).Operator.Kind, ShouldEqual, TokenTypeDoublePlus)
		So(ifStatement.If.Block.Statements[0].(*IncDecStatement).Expression.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"wawa")
		So(ifStatement.Elif[0].Condition.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeRightAngle)
		So(ifStatement.Elif[0].Condition.(*BinaryExpression).Left.(*MemberExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"m")
		So(ifStatement.Elif[0].Condition.(*BinaryExpression).Left.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"what")
		So(ifStatement.Elif[0].Condition.(*BinaryExpression).Right.(*BasicPrimaryExpression).Operand.(*FloatLit).Value.Str, ShouldEqual,
			"3.55")
		So(ifStatement.Elif[0].Block.Statements[0].(*AssignListStatement).Targets[0].(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"x")
		So(ifStatement.Elif[0].Block.Statements[0].(*AssignListStatement).Targets[1].(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"y")
		So(ifStatement.Elif[0].Block.Statements[0].(*AssignListStatement).Values[0].(*BasicPrimaryExpression).Operand.(*ExponentLit).Value.Str, ShouldEqual,
			"1.3e6")
		So(ifStatement.Elif[0].Block.Statements[0].(*AssignListStatement).Values[1].(*BasicPrimaryExpression).Operand.(*RuneLit).Value.Str, ShouldEqual,
			"Z")
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeEqual)
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Left.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"bb")
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Right.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypePlus)
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Right.(*BinaryExpression).Right.(*BinaryExpression).Operator.Kind, ShouldEqual,
			TokenTypeStar)
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Right.(*BinaryExpression).Right.(*BinaryExpression).Left.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str, ShouldEqual,
			"3")
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Right.(*BinaryExpression).Right.(*BinaryExpression).Right.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"dd")
	})
}
func TestSwitchStatement(t *testing.T) {
	Convey("测试分支语句: ", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `switch tom.grade {
        default {
            println("incorrect grade number!");
        }
        case 0...59 {
            println("Failed.");
        }
        case 60...100 {
            println("Lucky pass!");
        }
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "switch")

		switchStatement, isSwitch := parser.ParseStatement().(*SwitchStatement)
		So(isSwitch, ShouldEqual, true)
		So(switchStatement.Entry.(*MemberExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"tom")
		So(switchStatement.Entry.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"grade")
		So(switchStatement.Default.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(switchStatement.Default.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).Operand.(*StringLit).Value.Str, ShouldEqual,
			"incorrect grade number!")

		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Range.Start.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str, ShouldEqual,
			"0")
		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Range.IncludeEnd, ShouldEqual,
			true)
		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Range.End.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str, ShouldEqual,
			"59")
		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).Operand.(*StringLit).Value.Str, ShouldEqual,
			"Failed.")

		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Range.Start.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str, ShouldEqual,
			"60")
		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Range.IncludeEnd, ShouldEqual,
			true)
		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Range.End.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str, ShouldEqual,
			"100")
		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).Operand.(*StringLit).Value.Str, ShouldEqual,
			"Lucky pass!")
	})
}
func TestWhileStatement(t *testing.T) {
	Convey("测试 while 循环解析：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `while num < 10 {
			println(num);
			num++;
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "while")

		whileStatement, isWhile := parser.ParseStatement().(*WhileStatement)
		So(isWhile, ShouldEqual, true)

		So(whileStatement.Condition.(*BinaryExpression).Left.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"num")
		So(whileStatement.Condition.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeLeftAngle)
		So(whileStatement.Condition.(*BinaryExpression).Right.(*BasicPrimaryExpression).Operand.(*DecimalLit).Value.Str, ShouldEqual,
			"10")

		So(whileStatement.Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(whileStatement.Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"num")

		So(whileStatement.Block.Statements[1].(*IncDecStatement).Operator.Kind, ShouldEqual, TokenTypeDoublePlus)
		So(whileStatement.Block.Statements[1].(*IncDecStatement).Expression.(*BasicPrimaryExpression).Operand.(*OperandName).Name.Token.Str, ShouldEqual,
			"num")
	})
}
