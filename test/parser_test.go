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

	Convey("测试解析字面量值：数组字面量", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "[1,44, 9, 65]")
		So(parser.CurrentToken.Str, ShouldEqual, "[")

		a := parser.ParseLiteral()
		_, isArrayLit := a.(*ArrayLit)
		So(isArrayLit, ShouldEqual, true)
	})

	Convey("测试解析字面量值：映射表字面量", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `{
			key1: 1+33,
			mama: Color.Dark,
		}`)
		So(parser.CurrentToken.Str, ShouldEqual, "{")

		a := parser.ParseLiteral()
		tableLit, isTableLit := a.(*TableLit)
		So(isTableLit, ShouldEqual, true)

		So(tableLit.KeyValueList[0].Key.Token.Str, ShouldEqual, "key1")
		So(tableLit.KeyValueList[0].Value.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypePlus)

		So(tableLit.KeyValueList[1].Key.Token.Str, ShouldEqual, "mama")
		_, isMemberExpr := tableLit.KeyValueList[1].Value.(*MemberExpression)
		So(isMemberExpr, ShouldEqual, true)
	})

	Convey("测试解析字面量值：lambda", t, func() {
		parser1 := new(Parser)
		InitParserFromString(parser1, `
		var a = (m ,n int) float -> {
			println((m+n) * 2);
		};`)
		So(parser1.CurrentToken.Str, ShouldEqual, "var")

		lambdaVar, isVarDecl := parser1.ParseStatement().(*VarDeclStatement)
		So(isVarDecl, ShouldEqual, true)
		So(lambdaVar.Mutable, ShouldEqual, true)
		lambdaLit := lambdaVar.Declarations[0].InitValue.(*BasicPrimaryExpression).It.(*LambdaLit)
		So(lambdaLit.Signature.Arguments[0].Type.(*TypeName).Identifier.Token.Str, ShouldEqual, "int")
		So(lambdaLit.Signature.Arguments[0].Name.Token.Str, ShouldEqual, "m")
		So(lambdaVar.Declarations[0].VarName.Str, ShouldEqual, "a")
		So(lambdaLit.Signature.Returns[0].(*TypeName).Identifier.Token.Str, ShouldEqual,
			"float")
		So(lambdaLit.Result.(*BlockStatement).Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.GetName(),
			ShouldEqual, "println")

		parser2 := new(Parser)
		InitParserFromString(parser2, `
		friends.forEach((f) -> {
			f.greet();
		});`)
		So(parser2.CurrentToken.Str, ShouldEqual, "friends")

		call, isCall := parser2.ParseStatement().(*ExpressionStatement)
		So(isCall, ShouldEqual, true)
		So(call.Expression.(*CallExpression).Operand.(*MemberExpression).Member.Operand.Token.Str,
			ShouldEqual, "forEach")
		So((call.Expression.(*CallExpression)).Params[0].(*BasicPrimaryExpression).It.(*LambdaLit).Signature.Arguments[0].Name.Token.Str,
			ShouldEqual, "f")
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

		num := binaryExpression.Left.(*BinaryExpression).Left.(*BasicPrimaryExpression).It.(*OperandName)
		So(num.GetFullName(), ShouldEqual, "num")

		three := binaryExpression.Left.(*BinaryExpression).Right.(*BasicPrimaryExpression).It.(*DecimalLit)
		So(three.Value.Kind, ShouldEqual, TokenTypeDecimalInteger)
		So(three.Value.Str, ShouldEqual, "3")

		hex := binaryExpression.Right.(*BasicPrimaryExpression).It.(*HexadecimalLit)
		So(hex.Value.Kind, ShouldEqual, TokenTypeHexadecimalInteger)
		So(hex.Value.Str, ShouldEqual, "0x3f21")
	})
}
func TestCastExpression(t *testing.T) {
	Convey("测试类型强转表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "(3.1415 as float64)")
		So(parser.CurrentToken.Str, ShouldEqual, "(")

		castExpression, isCast := parser.ParseExpression().(*CastExpression)
		So(isCast, ShouldEqual, true)

		So(castExpression.Type.(*TypeName).Identifier.Token.Str, ShouldEqual, "float64")
		So(castExpression.Source.(*BasicPrimaryExpression).It.(*FloatLit).Value.Str, ShouldEqual,
			"3.1415")
		So(castExpression.Source.(*BasicPrimaryExpression).It.(*FloatLit).Accuracy, ShouldEqual,
			6)
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

		So(rangeExpression.Start.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual, "0")
		So(rangeExpression.End.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"arr")
		So(rangeExpression.End.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual, "length")
		So(rangeExpression.End.(*MemberExpression).Member.MemberNext, ShouldEqual, nil)
	})
}
func TestIndexSliceCallMemberExpression(t *testing.T) {
	Convey("测试索引表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "arr[i]")
		So(parser.CurrentToken.Str, ShouldEqual, "arr")

		indexExpression, isIndex := parser.ParseExpression().(*IndexExpression)
		So(isIndex, ShouldEqual, true)

		So(indexExpression.Operand.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "arr")
		So(indexExpression.Index.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "i")
	})

	Convey("测试切片表达式 1：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "arr[:4]")
		So(parser.CurrentToken.Str, ShouldEqual, "arr")

		sliceExpression, isSlice := parser.ParseExpression().(*SliceExpression)
		So(isSlice, ShouldEqual, true)

		So(sliceExpression.Operand.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "arr")
		So(sliceExpression.End.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str,
			ShouldEqual, "4")
	})

	Convey("测试切片表达式 2：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "arr[kk:5]")
		So(parser.CurrentToken.Str, ShouldEqual, "arr")

		sliceExpression, isSlice := parser.ParseExpression().(*SliceExpression)
		So(isSlice, ShouldEqual, true)

		So(sliceExpression.Operand.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "arr")
		So(sliceExpression.Start.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "kk")
		So(sliceExpression.End.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str,
			ShouldEqual, "5")
	})

	Convey("测试函数调用表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "funcA(i, 3.11)")
		So(parser.CurrentToken.Str, ShouldEqual, "funcA")

		callExpression, isCall := parser.ParseExpression().(*CallExpression)
		So(isCall, ShouldEqual, true)

		So(callExpression.Operand.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "funcA")
		So(callExpression.Params[0].(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "i")
		So(callExpression.Params[1].(*BasicPrimaryExpression).It.(*FloatLit).Value.Str,
			ShouldEqual, "3.11")
	})

	Convey("测试函数调用时、参数为 lambda 的类型自动推导 (Parser 部分体现为允许 无类型标注)：",
		t, func() {
			parser := new(Parser)
			InitParserFromString(parser, `friends.forEach((f) -> {
			f.greet();
		})`)
			So(parser.CurrentToken.Str, ShouldEqual, "friends")

			callExpression, isCall := parser.ParseExpression().(*CallExpression)
			So(isCall, ShouldEqual, true)
			So(callExpression.Operand.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str,
				ShouldEqual, "friends")
			So(callExpression.Operand.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual, "forEach")
			So(callExpression.Params[0].(*BasicPrimaryExpression).It.(*LambdaLit).Signature.Arguments[0].Name.Token.Str,
				ShouldEqual, "f")
			So(callExpression.Params[0].(*BasicPrimaryExpression).It.(*LambdaLit).Signature.Arguments[0].Type,
				ShouldEqual, nil)
		})

	Convey("测试成员访问表达式：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "request.query.page")
		So(parser.CurrentToken.Str, ShouldEqual, "request")

		memberExpression, isMemberExpr := parser.ParseExpression().(*MemberExpression)
		So(isMemberExpr, ShouldEqual, true)

		So(memberExpression.Operand.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
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

		So(callExpression.Operand.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "makeArr")
		So(callExpression.Params[0].(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypePlus)
		So(callExpression.Params[0].(*BinaryExpression).Left.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "dd")
		So(callExpression.Params[0].(*BinaryExpression).Right.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str,
			ShouldEqual, "7")

		So(sliceExpression.Start.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "mm")
		So(sliceExpression.End.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str,
			ShouldEqual, "6")
	})
}
func TestVarValDeclarationStatement(t *testing.T) {
	Convey("测试变量定义：1", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "var a int[] = [6, 7, 11];")
		So(parser.CurrentToken.Str, ShouldEqual, "var")

		varDeclStatement, isVarDecl := parser.ParseStatement().(*VarDeclStatement)
		So(isVarDecl, ShouldEqual, true)

		So(varDeclStatement.Mutable, ShouldEqual, true)
		So(varDeclStatement.Declarations[0].VarName.Str, ShouldEqual,
			"a")
		So(varDeclStatement.Declarations[0].Type.(*ArrayTypeLit).ElementType.(*TypeName).Identifier.Token.Str, ShouldEqual,
			"int")
		So(varDeclStatement.Declarations[0].InitValue.(*BasicPrimaryExpression).It.(*ArrayLit).ValueList[0].(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual,
			"6")
	})

	Convey("测试变量定义：1", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "var a rune[3] = ['c', 'd', '我'];")
		So(parser.CurrentToken.Str, ShouldEqual, "var")

		varDeclStatement, isVarDecl := parser.ParseStatement().(*VarDeclStatement)
		So(isVarDecl, ShouldEqual, true)

		So(varDeclStatement.Mutable, ShouldEqual, true)
		So(varDeclStatement.Declarations[0].VarName.Str, ShouldEqual,
			"a")
		So(varDeclStatement.Declarations[0].Type.(*ArrayTypeLit).ElementType.(*TypeName).Identifier.Token.Str, ShouldEqual,
			"rune")
		So(varDeclStatement.Declarations[0].Type.(*ArrayTypeLit).ArrayLength, ShouldEqual,
			3)
		So(varDeclStatement.Declarations[0].InitValue.(*BasicPrimaryExpression).It.(*ArrayLit).ValueList[2].(*BasicPrimaryExpression).It.(*RuneLit).Value.Str, ShouldEqual,
			"我")
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
		So(varDeclStatement.Declarations[0].InitValue.(*BasicPrimaryExpression).It.(*FloatLit).Value.Str, ShouldEqual,
			"3.14")

		So(varDeclStatement.Declarations[1].VarName.Str, ShouldEqual,
			"光速")
		So(varDeclStatement.Declarations[1].Type, ShouldEqual, nil)
		So(varDeclStatement.Declarations[1].InitValue.(*BasicPrimaryExpression).It.(*ExponentLit).Value.Str, ShouldEqual,
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
		So(leftUnary.Operand.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "m")
		So(leftUnary.Operand.(*MemberExpression).Member.Operand.Token.Str,
			ShouldEqual, "tt")

		rightUnary, isRightUnary := binaryExpression.Right.(*UnaryExpression)
		So(isRightUnary, ShouldEqual, true)
		So(rightUnary.Operator.Kind, ShouldEqual, TokenTypeWavy)
		So(rightUnary.Operand.(*IndexExpression).Operand.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "ss")
		So(rightUnary.Operand.(*IndexExpression).Operand.(*MemberExpression).Member.Operand.Token.Str,
			ShouldEqual, "k")
		So(rightUnary.Operand.(*IndexExpression).Index.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str,
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
		So(newInstanceExpression.Class.(*TypeName).Identifier.Token.Str,
			ShouldEqual, "Student")
		So(newInstanceExpression.InitParams[0].(*BasicPrimaryExpression).It.(*StringLit).Value.Str,
			ShouldEqual, "Peter")
		So(newInstanceExpression.InitParams[1].(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str,
			ShouldEqual, "18")
	})

	Convey("测试对象实例新建表达式：含泛型参数", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, "new Array<string>(3)")
		So(parser.CurrentToken.Str, ShouldEqual, "new")

		newInstanceExpression, isNewInstance := parser.ParseExpression().(*NewInstanceExpression)
		So(isNewInstance, ShouldEqual, true)
		So(newInstanceExpression, ShouldNotEqual, nil)
		So(newInstanceExpression.Class.(*GenericsTypeLit).BasicType.Identifier.Token.Str,
			ShouldEqual, "Array")
		So(newInstanceExpression.Class.(*GenericsTypeLit).GenericsArgs[0].(*TypeName).Identifier.Token.Str,
			ShouldEqual, "string")
		So(newInstanceExpression.InitParams[0].(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str,
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

		So(assignListStmt.Targets[0].(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "num_a")
		So(assignListStmt.Targets[1].(*IndexExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).GetFullName(),
			ShouldEqual, "x")
		So(assignListStmt.Targets[1].(*IndexExpression).Index.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str,
			ShouldEqual, "1")
		So(assignListStmt.Values[0].(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str,
			ShouldEqual, "3")
		So(assignListStmt.Values[1].(*BasicPrimaryExpression).It.(*StringLit).Value.Str,
			ShouldEqual, "hello")
	})
}
func TestImportStatement(t *testing.T) {
	Convey("测试导入模块语句：1", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `import stdlib as std;`)
		So(parser.CurrentToken.Str, ShouldEqual, "import")

		stmt := parser.ParseStatement()
		singleGlobalImportStatement, isSingleGlobal := stmt.(*SingleGlobalImportStatement)
		So(isSingleGlobal, ShouldEqual, true)
		So(singleGlobalImportStatement.Element.ModuleName.GetName(), ShouldEqual, "stdlib")
		So(singleGlobalImportStatement.Element.As.Token.Str, ShouldEqual, "std")
	})

	Convey("测试导入模块语句：2", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `from httplib import Request as Req;`)
		So(parser.CurrentToken.Str, ShouldEqual, "from")

		stmt := parser.ParseStatement()
		singleImportStatement, isSingleImport := stmt.(*SingleFromImportStatement)
		So(isSingleImport, ShouldEqual, true)
		So(singleImportStatement.From.GetName(), ShouldEqual, "httplib")
		So(singleImportStatement.Element.ModuleName.GetName(), ShouldEqual, "Request")
		So(singleImportStatement.Element.As.Token.Str, ShouldEqual, "Req")
	})

	Convey("测试导入模块语句：3", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `from httplib import {
      Request as Req,
			Response as Resp
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "from")

		listImportStatement, isListImport := parser.ParseStatement().(*ListImportStatement)
		So(isListImport, ShouldEqual, true)
		So(listImportStatement.From.GetName(), ShouldEqual, "httplib")
		So(listImportStatement.Elements[0].ModuleName.GetName(), ShouldEqual, "Request")
		So(listImportStatement.Elements[0].As.Token.Str, ShouldEqual, "Req")
		So(listImportStatement.Elements[1].ModuleName.GetName(), ShouldEqual, "Response")
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
		So(ifStatement.If.Condition.(*UnaryExpression).Operand.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"screen")
		So(ifStatement.If.Condition.(*UnaryExpression).Operand.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"closed")
		So(ifStatement.If.Condition.(*UnaryExpression).Operand.(*MemberExpression).Member.MemberNext, ShouldEqual,
			nil)
		So(ifStatement.If.Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(ifStatement.If.Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).It.(*StringLit).Value.Str, ShouldEqual,
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

		So(ifStatement.If.Condition.(*BinaryExpression).Left.(*CallExpression).Operand.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"t")
		So(ifStatement.If.Condition.(*BinaryExpression).Left.(*CallExpression).Operand.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"getYear")
		So(len(ifStatement.If.Condition.(*BinaryExpression).Left.(*CallExpression).Params), ShouldEqual,
			0)
		So(ifStatement.If.Condition.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeDoubleEqual)
		So(ifStatement.If.Block.Statements[0].(*IncDecStatement).Operator.Kind, ShouldEqual, TokenTypeDoublePlus)
		So(ifStatement.If.Block.Statements[0].(*IncDecStatement).Expression.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"wawa")
		So(ifStatement.Elif[0].Condition.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeRightAngle)
		So(ifStatement.Elif[0].Condition.(*BinaryExpression).Left.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"m")
		So(ifStatement.Elif[0].Condition.(*BinaryExpression).Left.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"what")
		So(ifStatement.Elif[0].Condition.(*BinaryExpression).Right.(*BasicPrimaryExpression).It.(*FloatLit).Value.Str, ShouldEqual,
			"3.55")
		So(ifStatement.Elif[0].Block.Statements[0].(*AssignListStatement).Targets[0].(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"x")
		So(ifStatement.Elif[0].Block.Statements[0].(*AssignListStatement).Targets[1].(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"y")
		So(ifStatement.Elif[0].Block.Statements[0].(*AssignListStatement).Values[0].(*BasicPrimaryExpression).It.(*ExponentLit).Value.Str, ShouldEqual,
			"1.3e6")
		So(ifStatement.Elif[0].Block.Statements[0].(*AssignListStatement).Values[1].(*BasicPrimaryExpression).It.(*RuneLit).Value.Str, ShouldEqual,
			"Z")
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeEqual)
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Left.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"bb")
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Right.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypePlus)
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Right.(*BinaryExpression).Right.(*BinaryExpression).Operator.Kind, ShouldEqual,
			TokenTypeStar)
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Right.(*BinaryExpression).Right.(*BinaryExpression).Left.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual,
			"3")
		So(ifStatement.Else.Statements[0].(*ExpressionStatement).Expression.(*BinaryExpression).Right.(*BinaryExpression).Right.(*BinaryExpression).Right.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
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
		So(switchStatement.Entry.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"tom")
		So(switchStatement.Entry.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"grade")
		So(switchStatement.Default.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(switchStatement.Default.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).It.(*StringLit).Value.Str, ShouldEqual,
			"incorrect grade number!")

		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Range.Start.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual,
			"0")
		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Range.IncludeEnd, ShouldEqual,
			true)
		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Range.End.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual,
			"59")
		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(switchStatement.Cases[0].(*SwitchStatementRangeCase).Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).It.(*StringLit).Value.Str, ShouldEqual,
			"Failed.")

		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Range.Start.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual,
			"60")
		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Range.IncludeEnd, ShouldEqual,
			true)
		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Range.End.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual,
			"100")
		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(switchStatement.Cases[1].(*SwitchStatementRangeCase).Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).It.(*StringLit).Value.Str, ShouldEqual,
			"Lucky pass!")
	})
}
func TestWhileStatement(t *testing.T) {
	Convey("测试 while 循环解析：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `while num < 10 {
			println(num);
			num++;
			if num > 3 {
				break;
			} else {
        continue;
			}
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "while")

		whileStatement, isWhile := parser.ParseStatement().(*WhileStatement)
		So(isWhile, ShouldEqual, true)

		So(whileStatement.Condition.(*BinaryExpression).Left.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"num")
		So(whileStatement.Condition.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeLeftAngle)
		So(whileStatement.Condition.(*BinaryExpression).Right.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual,
			"10")

		So(whileStatement.Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"println")
		So(whileStatement.Block.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"num")

		So(whileStatement.Block.Statements[1].(*IncDecStatement).Operator.Kind, ShouldEqual, TokenTypeDoublePlus)
		So(whileStatement.Block.Statements[1].(*IncDecStatement).Expression.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"num")

		So(whileStatement.Block.Statements[2].(*IfStatement).If.Condition.(*BinaryExpression).Left.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"num")
		So(whileStatement.Block.Statements[2].(*IfStatement).If.Condition.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeRightAngle)
		So(whileStatement.Block.Statements[2].(*IfStatement).If.Condition.(*BinaryExpression).Right.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual,
			"3")
		So(whileStatement.Block.Statements[2].(*IfStatement).If.Block.Statements[0].(*BreakStatement).Token.Str, ShouldEqual,
			"break")
		So(whileStatement.Block.Statements[2].(*IfStatement).Else.Statements[0].(*ContinueStatement).Token.Str, ShouldEqual,
			"continue")
	})
}
func TestForStatement(t *testing.T) {
	Convey("测试 for 循环语句：1", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `for var i = 0, j = arr.length; i <= j; i++, j-- {
      println(i + j);
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "for")

		forStatement, isFor := parser.ParseStatement().(*ForStatement)
		So(isFor, ShouldEqual, true)

		So(forStatement.Initial.(*VarDeclStatement).Mutable, ShouldEqual, true)
		So(forStatement.Initial.(*VarDeclStatement).Declarations[0].VarName.Str, ShouldEqual,
			"i")
		So(forStatement.Initial.(*VarDeclStatement).Declarations[0].InitValue.(*BasicPrimaryExpression).It.(*DecimalLit).Value.Str, ShouldEqual,
			"0")
		So(forStatement.Initial.(*VarDeclStatement).Declarations[1].VarName.Str, ShouldEqual,
			"j")
		So(forStatement.Initial.(*VarDeclStatement).Declarations[1].InitValue.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"arr")
		So(forStatement.Initial.(*VarDeclStatement).Declarations[1].InitValue.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"length")

		So(forStatement.Condition.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeLeftAngleEqual)
		So(forStatement.Condition.(*BinaryExpression).Left.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"i")
		So(forStatement.Condition.(*BinaryExpression).Right.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"j")

		So(forStatement.Appendix[0].(*IncDecStatement).Operator.Kind, ShouldEqual, TokenTypeDoublePlus)
		So(forStatement.Appendix[0].(*IncDecStatement).Expression.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"i")
		So(forStatement.Appendix[1].(*IncDecStatement).Operator.Kind, ShouldEqual, TokenTypeDoubleMinus)
		So(forStatement.Appendix[1].(*IncDecStatement).Expression.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"j")
	})

	Convey("测试 for 循环语句：无初始化、无尾缀操作", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `for ;i < arr.length; {
      println(arr[i]);
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "for")

		forStatement, isFor := parser.ParseStatement().(*ForStatement)
		So(isFor, ShouldEqual, true)

		So(forStatement.Initial, ShouldEqual, nil)

		So(forStatement.Condition.(*BinaryExpression).Operator.Kind, ShouldEqual, TokenTypeLeftAngle)
		So(forStatement.Condition.(*BinaryExpression).Left.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"i")
		So(forStatement.Condition.(*BinaryExpression).Right.(*MemberExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual,
			"arr")
		So(forStatement.Condition.(*BinaryExpression).Right.(*MemberExpression).Member.Operand.Token.Str, ShouldEqual,
			"length")
		So(forStatement.Condition.(*BinaryExpression).Right.(*MemberExpression).Member.MemberNext, ShouldEqual,
			nil)

		So(len(forStatement.Appendix), ShouldEqual, 0)
	})
}
func TestEachStatement(t *testing.T) {
	Convey("测试 each 循环语句：1 无 key", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `each num in [1,4,5,99] {
			println(num);
		}`)
		So(parser.CurrentToken.Str, ShouldEqual, "each")

		eachStatement, isEach := parser.ParseStatement().(*EachStatement)
		So(isEach, ShouldEqual, true)

		So(eachStatement.Element.Token.Str, ShouldEqual, "num")
		So(eachStatement.Key, ShouldEqual, nil)
		So(len(eachStatement.Target.(*BasicPrimaryExpression).It.(*ArrayLit).ValueList), ShouldEqual, 4)
	})

	Convey("测试 each 循环语句：2 有 key", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `each num, i in [1,4,5,99] {
			println("No." + i + num);
		}`)
		So(parser.CurrentToken.Str, ShouldEqual, "each")

		eachStatement, isEach := parser.ParseStatement().(*EachStatement)
		So(isEach, ShouldEqual, true)

		So(eachStatement.Element.Token.Str, ShouldEqual, "num")
		So(eachStatement.Key.Token.Str, ShouldEqual, "i")
		So(len(eachStatement.Target.(*BasicPrimaryExpression).It.(*ArrayLit).ValueList), ShouldEqual, 4)
	})
}
func TestFnStatement(t *testing.T) {
	Convey("测试函数定义语句：1", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `fn fibonacci(n int) int {
        var a = n % 2, b = 1;
        for var i = 0; i < n/2; i++ {
            a += b;
            b += a;
        }

        return a;
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "fn")

		fnStatement, isFn := parser.ParseStatement().(*FunctionDeclarationStatement)
		So(isFn, ShouldEqual, true)

		So(fnStatement.Name.Token.Str, ShouldEqual, "fibonacci")
		So(fnStatement.Signature.Arguments[0].Name.Token.Str, ShouldEqual, "n")
		So(fnStatement.Signature.Arguments[0].Type.(*TypeName).Identifier.Token.Str, ShouldEqual, "int")
		So(fnStatement.Signature.Returns[0].(*TypeName).Identifier.Token.Str, ShouldEqual, "int")
	})

	Convey("测试函数定义语句：2", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `fn initMapWithAPair<T, K>(
			n1 T, 
			n2 K
		) Map<T, K>, bool throws NullPointerException {
			return new Map(n1, n2), true;
		}`)
		So(parser.CurrentToken.Str, ShouldEqual, "fn")

		fnStatement, isFn := parser.ParseStatement().(*FunctionDeclarationStatement)
		So(isFn, ShouldEqual, true)

		So(fnStatement.Name.Token.Str, ShouldEqual, "initMapWithAPair")
		So(fnStatement.Generics.Args[0].ArgName.Token.Str, ShouldEqual, "T")
		So(fnStatement.Generics.Args[1].ArgName.Token.Str, ShouldEqual, "K")
		So(fnStatement.Signature.Arguments[0].Name.Token.Str, ShouldEqual, "n1")
		So(fnStatement.Signature.Arguments[0].Type.(*TypeName).Identifier.Token.Str, ShouldEqual, "T")
		So(fnStatement.Signature.Arguments[1].Name.Token.Str, ShouldEqual, "n2")
		So(fnStatement.Signature.Arguments[1].Type.(*TypeName).Identifier.Token.Str, ShouldEqual, "K")
		So(fnStatement.Signature.Returns[0].(*GenericsTypeLit).BasicType.Identifier.Token.Str, ShouldEqual,
			"Map")
		So(fnStatement.Signature.Returns[0].(*GenericsTypeLit).GenericsArgs[0].(*TypeName).Identifier.Token.Str, ShouldEqual,
			"T")
		So(fnStatement.Signature.Returns[0].(*GenericsTypeLit).GenericsArgs[1].(*TypeName).Identifier.Token.Str, ShouldEqual,
			"K")
		So(fnStatement.Signature.Returns[1].(*TypeName).Identifier.Token.Str, ShouldEqual, "bool")
		So(fnStatement.Signature.Throws[0].(*TypeName).Identifier.Token.Str, ShouldEqual, "NullPointerException")
	})
}
func TestClassStatement(t *testing.T) {
	Convey("测试类定义语句：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `class VideoDisk<T, K> : Disk<K> <- Playable<T> {
        var time Date,
						movieName string,
        		movieDirector string,
        		productionCompany string,
						renter Customer;

        var hasBeenRented bool = false;

				fn VideoDisk(name string, director string, company string) {
					this.movieName = name;
					this.movieDirector = director;
					this.productionCompany = company;
				}

        public fn rent(c Customer) {
            this.renter = Customer;
            this.hasBeenRented = true;
        }
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "class")

		classStatement, isClass := parser.ParseStatement().(*ClassDeclarationStatement)
		So(isClass, ShouldEqual, true)

		So(classStatement.Definition.Name.Token.Str, ShouldEqual, "VideoDisk")
		So(classStatement.Definition.Generics.Args[0].ArgName.Token.Str, ShouldEqual, "T")
		So(classStatement.Definition.Generics.Args[1].ArgName.Token.Str, ShouldEqual, "K")

		So(classStatement.Extends.Name.Token.Str, ShouldEqual, "Disk")
		So(classStatement.Extends.Generics.Args[0].ArgName.Token.Str, ShouldEqual, "K")

		So(classStatement.Implements[0].Name.Token.Str, ShouldEqual, "Playable")
		So(classStatement.Implements[0].Generics.Args[0].ArgName.Token.Str, ShouldEqual, "T")

		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[0].VarName.Str, ShouldEqual,
			"time")
		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[0].Type.(*TypeName).Identifier.Token.Str, ShouldEqual,
			"Date")

		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[1].VarName.Str, ShouldEqual,
			"movieName")
		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[1].Type.(*TypeName).Identifier.Token.Str, ShouldEqual,
			"string")

		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[2].VarName.Str, ShouldEqual,
			"movieDirector")
		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[2].Type.(*TypeName).Identifier.Token.Str, ShouldEqual,
			"string")

		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[3].VarName.Str, ShouldEqual,
			"productionCompany")
		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[3].Type.(*TypeName).Identifier.Token.Str, ShouldEqual,
			"string")

		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[4].VarName.Str, ShouldEqual,
			"renter")
		So(classStatement.Members[0].(*ClassMemberVar).VarDecl.Declarations[4].Type.(*TypeName).Identifier.Token.Str, ShouldEqual,
			"Customer")

		So(classStatement.Members[3].(*ClassMemberMethod).MethodDecl.Name.Token.Str, ShouldEqual,
			"rent")
		So(classStatement.Members[3].(*ClassMemberMethod).Scope, ShouldEqual, ClassMemberScopePublic)
	})
}
func TestInterfaceStatement(t *testing.T) {
	Convey("测试接口定义语句：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `interface A<T> : B<T> {
        public  fn cc<T>() string throws MMException;
        private fn dd() int;
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "interface")

		interfaceStatement, isInterface := parser.ParseStatement().(*InterfaceDeclarationStatement)
		So(isInterface, ShouldEqual, true)

		So(interfaceStatement.Definition.Name.Token.Str, ShouldEqual, "A")
		So(interfaceStatement.Definition.Generics.Args[0].ArgName.Token.Str, ShouldEqual, "T")

		So(interfaceStatement.Extends.Name.Token.Str, ShouldEqual, "B")
		So(interfaceStatement.Extends.Generics.Args[0].ArgName.Token.Str, ShouldEqual, "T")

		So(interfaceStatement.Methods[0].Name.Token.Str, ShouldEqual,
			"cc")
		So(interfaceStatement.Methods[0].Generics.Args[0].ArgName.Token.Str, ShouldEqual, "T")
		So(interfaceStatement.Methods[0].Signature.Returns[0].(*TypeName).Identifier.Token.Str, ShouldEqual,
			"string")
		So(interfaceStatement.Methods[0].Signature.Throws[0].(*TypeName).Identifier.Token.Str, ShouldEqual,
			"MMException")
		So(interfaceStatement.Methods[0].Scope, ShouldEqual, ClassMemberScopePublic)
	})
}
func TestTryCatchStatement(t *testing.T) {
	Convey("测试接口定义语句：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `try {
			val n = 3 / 0;
    } catch e MathException {
			println(e.message());
    } catch e CCException {
			println("Just test.");
		} finally {
      println("hahaha, it's ok");
    }`)
		So(parser.CurrentToken.Str, ShouldEqual, "try")

		tryCatchStmt, isTryCatch := parser.ParseStatement().(*TryCatchStatement)
		So(isTryCatch, ShouldEqual, true)

		So(tryCatchStmt.TryBlock.Statements[0].(*VarDeclStatement).Mutable, ShouldEqual, false)
		So(tryCatchStmt.TryBlock.Statements[0].(*VarDeclStatement).Declarations[0].VarName.Str, ShouldEqual, "n")
		So(tryCatchStmt.TryBlock.Statements[0].(*VarDeclStatement).Declarations[0].InitValue.(*BinaryExpression).Operator.Kind,
			ShouldEqual, TokenTypeSlash)

		So(tryCatchStmt.Handlers[0].Name.Token.Str, ShouldEqual, "e")
		So(tryCatchStmt.Handlers[0].ErrorType.(*TypeName).Identifier.Token.Str, ShouldEqual, "MathException")

		So(tryCatchStmt.Handlers[1].Name.Token.Str, ShouldEqual, "e")
		So(tryCatchStmt.Handlers[1].ErrorType.(*TypeName).Identifier.Token.Str, ShouldEqual, "CCException")

		So(tryCatchStmt.Finally.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Operand.(*BasicPrimaryExpression).It.(*OperandName).Name.Token.Str, ShouldEqual, "println")
		So(tryCatchStmt.Finally.Statements[0].(*ExpressionStatement).Expression.(*CallExpression).Params[0].(*BasicPrimaryExpression).It.(*StringLit).Value.Str, ShouldEqual, "hahaha, it's ok")
	})
}
func TestPackageStatement(t *testing.T) {
	Convey("测试包名定义语句：", t, func() {
		parser := new(Parser)
		InitParserFromString(parser, `package test;`)
		So(parser.CurrentToken.Str, ShouldEqual, "package")

		pkgStmt, isPkg := parser.ParseStatement().(*PackageStatement)
		So(isPkg, ShouldEqual, true)
		So(pkgStmt.Name.GetName(), ShouldEqual, "test")
	})
}
