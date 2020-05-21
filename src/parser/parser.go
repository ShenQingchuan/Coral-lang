package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
	"fmt"
)

type Parser struct {
	Lexer   *Lexer
	Program *Program

	LastToken    *Token
	CurrentToken *Token
}

func InitParserFromBytes(parser *Parser, content []byte) {
	parser.Lexer = new(Lexer)
	InitLexerFromBytes(parser.Lexer, content)
	parser.PeekNextToken() // 统一获取到第一个 Token
}
func InitParserFromString(parser *Parser, content string) {
	parser.Lexer = new(Lexer)
	InitLexerFromString(parser.Lexer, content)
	parser.PeekNextToken() // 统一获取到第一个 Token
}
func (parser *Parser) PeekNextToken() {
	token, err := parser.Lexer.GetNextToken()
	if err != nil {
		CoralErrorCrashHandler(err)
	}

	parser.LastToken = parser.CurrentToken
	parser.CurrentToken = token
}
func (parser *Parser) GetCurrentTokenPos() string {
	return fmt.Sprintf("line %d:%d: ", parser.CurrentToken.Line, parser.CurrentToken.Col)
}
func GetBinaryOperatorPriority(token *Token) int {
	switch token.Kind {
	case TokenTypeDoubleStar:
		return 2
	case TokenTypeSlash, TokenTypeStar, TokenTypePercent:
		return 3
	case TokenTypePlus, TokenTypeMinus:
		return 4
	case TokenTypeDoubleLeftAngle, TokenTypeDoubleRightAngle:
		return 5
	case TokenTypeRightAngle, TokenTypeLeftAngle, TokenTypeRightAngleEqual, TokenTypeLeftAngleEqual:
		return 6
	case TokenTypeDoubleEqual, TokenTypeBangEqual:
		return 7
	case TokenTypeAmpersand:
		return 8
	case TokenTypeCaret:
		return 9
	case TokenTypeVertical:
		return 10
	case TokenTypeDoubleAmpersand:
		return 11
	case TokenTypeDoubleVertical:
		return 12
	case TokenTypeEqual, TokenTypeSlashEqual,
		TokenTypeStarEqual, TokenTypePercentEqual,
		TokenTypePlusEqual, TokenTypeMinusEqual,
		TokenTypeDoubleLeftAngleEqual, TokenTypeDoubleRightAngleEqual,
		TokenTypeAmpersandEqual, TokenTypeCaretEqual, TokenTypeVerticalEqual:
		return 14
	}

	return 99 // 返回一个极大值表示无优先级
}

// IDENTIFIER ('.' IDENTIFIER)*
func (parser *Parser) ParseIdentifierList() []*Identifier {
	if parser.CurrentToken.Kind != TokenTypeIdentifier {
		return nil
	}

	// 先添加传入的 token，已确定其为 identifier
	var identifierList []*Identifier
	identifierList = append(identifierList, &Identifier{Name: parser.CurrentToken})

	for {
		parser.PeekNextToken()
		if parser.CurrentToken.Kind != TokenTypeDot { // 需要一个 '.' 来分隔
			break
		}
		parser.PeekNextToken()
		if parser.CurrentToken.Kind != TokenTypeIdentifier { // '.' 后的 Identifier
			break
		}
		identifierList = append(identifierList, &Identifier{Name: parser.CurrentToken})
	}
	return identifierList
}

func (parser *Parser) ParseProgram() {
	for stmt := parser.ParseStatement(); stmt != nil; stmt = parser.ParseStatement() {
		// stmt 为 nil 的情况中其实早已被 CoralErrorCrashHandler 处理并退出了
		parser.Program.Root = append(parser.Program.Root, stmt)
	}
}

func (parser *Parser) ParseStatement() Statement {
	parser.PeekNextToken()

	if simpleStmt := parser.ParseSimpleStatement(); simpleStmt != nil {
		return simpleStmt
	}

	return nil
}

func (parser *Parser) ParseSimpleStatement() SimpleStatement {
	if expression := parser.ParseExpression(); expression != nil {
		parser.PeekNextToken()
		if parser.CurrentToken.Kind == TokenTypeSemi {
			return expression
		} else {
			CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
				"expected a semicolon to terminate this statement!", ParsingUnexpected))
		}
	}

	return nil
}

func (parser *Parser) ParseExpression() Expression {
	if primaryExpr := parser.ParsePrimaryExpression(); primaryExpr != nil {
		return primaryExpr
	}
	if newInstanceExpression := parser.ParseNewInstanceExpression(); newInstanceExpression != nil {
		return newInstanceExpression
	}
	if unaryExpression := parser.ParseUnaryExpression(); unaryExpression != nil {
		return unaryExpression
	}
	if binaryExpression := parser.ParseBinaryExpression(); binaryExpression != nil {
		return binaryExpression
	}

	return nil
}

func (parser *Parser) ParseTypeDescription() TypeDescription {
	// 如果是左中括号
	if parser.CurrentToken.Kind == TokenTypeLeftBracket {
		parser.PeekNextToken()
		typeDescription := parser.ParseTypeDescription()
		if typeDescription != nil {
			arrayTypeLit := new(ArrayTypeLit)
			arrayTypeLit.ElementType = &typeDescription

			// 此时当前 token 应为 ']'
			if parser.CurrentToken.Kind == TokenTypeRightBracket {
				return arrayTypeLit
			} else {
				CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
					fmt.Sprintf("expected a right bracket for type literal but got '%s'", parser.CurrentToken.Str),
					ParsingUnexpected))
			}
		} else {
			CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
				fmt.Sprintf("expected a type description but got '%s'", parser.CurrentToken.Str),
				ParsingUnexpected))
		}
	}

	// 如果是 identifier 说明可能是 GenericsLit
	if parser.CurrentToken.Kind == TokenTypeIdentifier {
		typeName := parser.ParseTypeName()
		if typeName != nil {
			parser.PeekNextToken()
			// 结束 typeName 解析后如果是一个 '<' 则进入 Generics 解析
			if parser.CurrentToken.Kind == TokenTypeLeftAngle {
				genericsTypeLit := new(GenericsTypeLit)
				genericsTypeLit.BasicType = typeName

				parser.PeekNextToken() // 越过 '<'
				for {
					genericsArg := parser.ParseTypeName()
					if genericsArg != nil {
						genericsTypeLit.GenericsArgs = append(genericsTypeLit.GenericsArgs, genericsArg)
					}
					if parser.CurrentToken.Kind != TokenTypeComma {
						CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
							fmt.Sprintf("expected a comma to seperate generics arguments but got '%s'",
								parser.CurrentToken.Str), ParsingUnexpected))
					} else if parser.CurrentToken.Kind == TokenTypeRightAngle {
						return genericsTypeLit // 遇到 '>' 结束泛型参数解析
					}
				}
			} else {
				// 否则就将 typeName 返回作为该 typeDescription
				return typeName
			}
		}
	}

	return nil
}

func (parser *Parser) ParseTypeName() *TypeName {
	if parser.CurrentToken.Kind != TokenTypeIdentifier {
		return nil
	}

	// 先添加传入的 token，已确定其为 identifier
	typeName := new(TypeName)
	identifierList := parser.ParseIdentifierList()
	typeName.NameList = identifierList

	return typeName
}

func (parser *Parser) ParsePrimaryExpression() PrimaryExpression {
	operand := parser.ParseBasicPrimaryExpression()
	if operand == nil {
		return nil // 没有解析到 operand 说明不是
	}
	// 至此已经获取到了 operand，后续 三种情况如果解析中都不能匹配
	// 则最后返回 operand

	// try: slice/index
	// 为什么先尝试 slice? : 因为可能存在 arr[:3]
	parser.PeekNextToken()

	// 如果是中括号 '['
	if parser.CurrentToken.Kind == TokenTypeLeftBracket {
		parser.PeekNextToken()

		// 如果确定是 '[:'
		if parser.CurrentToken.Kind == TokenTypeColon {
			sliceExpr := new(SliceExpression)
			sliceExpr.Operand = operand

			// 解析切片终点表达式
			parser.PeekNextToken()
			end := parser.ParseExpression()
			if end == nil {
				CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
					"expected an expression to be end position for slice!", ParsingUnexpected))
			}
			sliceExpr.End = end
			return sliceExpr
		}

		start := parser.ParseExpression()
		if start == nil {
			CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
				"expected an expression to be an index/key or a start position for slice!", ParsingUnexpected))
		}

		parser.PeekNextToken()
		if parser.CurrentToken.Kind == TokenTypeColon {
			sliceExpr := new(SliceExpression)
			sliceExpr.Operand = operand
			sliceExpr.Start = start

			// 解析切片终点表达式
			parser.PeekNextToken()
			end := parser.ParseExpression()
			if end == nil {
				CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
					"expected an expression to be an end position for slice!", ParsingUnexpected))
			}
			sliceExpr.End = end
			return sliceExpr
		} else if parser.CurrentToken.Kind == TokenTypeRightBracket {
			// 只有一个表达式就遇到了右括号
			indexExpr := new(IndexExpression)
			indexExpr.Operand = operand
			indexExpr.Index = start
			return indexExpr
		}
	} else if parser.CurrentToken.Kind == TokenTypeLeftParen {
		parser.PeekNextToken() // 移过当前的左括号，到下一个 token
		callExpression := new(CallExpression)
		callExpression.Operand = operand
		for expression := parser.ParseExpression(); expression != nil; expression = parser.ParseExpression() {
			callExpression.Params = append(callExpression.Params, &expression)
		}
		// 结束循环时，检测是否停留于 token ')'
		if parser.CurrentToken.Kind == TokenTypeRightParen {
			return callExpression
		} else {
			CoralErrorCrashHandler(NewCoralError("Compile Error",
				"expected a right parenthesis for function calling!", ParsingUnexpected))
		}
	}

	return operand
}

// 解析基本表达式的 operand 部分
func (parser *Parser) ParseBasicPrimaryExpression() *BasicPrimaryExpression {
	literal := parser.ParseLiteral()
	if literal != nil {
		return &BasicPrimaryExpression{Operand: literal}
	} // 如果 literal 为空则另一种情况

	operandName := parser.ParseOperandName()
	if operandName != nil {
		return &BasicPrimaryExpression{Operand: operandName}
	}

	return nil
}

// 解析 operand 的 literal 情况
func (parser *Parser) ParseLiteral() Literal {
	switch parser.CurrentToken.Kind {
	default:
		return nil
	case TokenTypeDecimalInteger:
		return &DecimalLit{Value: parser.CurrentToken}
	case TokenTypeHexadecimalInteger:
		return &HexadecimalLit{Value: parser.CurrentToken}
	case TokenTypeOctalInteger:
		return &OctalLit{Value: parser.CurrentToken}
	case TokenTypeBinaryInteger:
		return &BinaryLit{Value: parser.CurrentToken}
	case TokenTypeFloat:
		return &FloatLit{Value: parser.CurrentToken}
	case TokenTypeExponent:
		return &ExponentLit{Value: parser.CurrentToken}
	case TokenTypeNil:
		return &NilLit{Value: parser.CurrentToken}
	case TokenTypeTrue:
		return &TrueLit{Value: parser.CurrentToken}
	case TokenTypeFalse:
		return &FalseLit{Value: parser.CurrentToken}
	}
}

// 解析 operand 的 operandName（变量名）情况
func (parser *Parser) ParseOperandName() *OperandName {
	if parser.CurrentToken.Kind != TokenTypeIdentifier {
		return nil
	}

	// 先添加传入的 token，已确定其为 identifier
	operandName := new(OperandName)
	identifierList := parser.ParseIdentifierList()
	operandName.NameList = identifierList

	return operandName
}

// 解析 新建对象实例 表达式
func (parser *Parser) ParseNewInstanceExpression() *NewInstanceExpression {
	if parser.CurrentToken.Kind != TokenTypeNew {
		return nil
	}

	parser.PeekNextToken()
	typeDescription := parser.ParseTypeDescription()
	if typeDescription == nil {
		CoralErrorCrashHandler(NewCoralError("Compile Error",
			"expected a type description for object instance creating!", ParsingUnexpected))
	}
	parser.PeekNextToken()
	if parser.CurrentToken.Kind == TokenTypeLeftParen {
		parser.PeekNextToken() // 移过当前的左括号，到下一个 token
		newInstanceExpression := new(NewInstanceExpression)
		newInstanceExpression.Class = &typeDescription
		for expression := parser.ParseExpression(); expression != nil; expression = parser.ParseExpression() {
			newInstanceExpression.InitParams = append(newInstanceExpression.InitParams, &expression)
		}
		// 结束循环时，检测是否停留于 token ')'
		if parser.CurrentToken.Kind == TokenTypeRightParen {
			return newInstanceExpression
		} else {
			CoralErrorCrashHandler(NewCoralError("Compile Error",
				"expected a right parenthesis for constructor method's ending!", ParsingUnexpected))
		}
	} else {
		CoralErrorCrashHandler(NewCoralError("Compile Error",
			"expected a left parenthesis for constructor method's calling!", ParsingUnexpected))
	}

	return nil
}

// 解析 单目表达式
func (parser *Parser) ParseUnaryExpression() *UnaryExpression {
	switch parser.CurrentToken.Kind {
	case TokenTypeMinus, TokenTypeBang, TokenTypeWavy:
		unaryExpression := new(UnaryExpression)
		unaryExpression.Operator = parser.CurrentToken
		if operand := parser.ParseExpression(); operand != nil {
			unaryExpression.Operand = operand
			return unaryExpression
		} else {
			CoralErrorCrashHandler(NewCoralError("Compile Error",
				"missing operand for unary expression!", ParsingUnexpected))
		}
	}

	return nil
}

func (parser *Parser) ParseBinaryExpression() Expression {
	if left := parser.ParseExpression(); left != nil {
		if priority := GetBinaryOperatorPriority(parser.CurrentToken); priority != 99 {
			// 证明是 二元运算符
			binaryExpression := new(BinaryExpression)
			binaryExpression.Operator = parser.CurrentToken
			binaryExpression.Left = left

			parser.PeekNextToken() // 移过当前操作符节点

			// 先暂时形成 左结构，然后解析右边节点（即体现：自左向右结合）
			if r := parser.ParseExpression(); r != nil {
				right, rightIsBinary := r.(*BinaryExpression)
				if rightIsBinary { // 右边节点也是二元表达式，需要进行优先级判断、旋转树
					if GetBinaryOperatorPriority(right.Operator) > priority {
						// 即 右边节点的 priority 数值大，反而优先级别低，应作父节点
						binaryExpression.Right = right.Left // 补充原树的右节点
						right.Left = binaryExpression       // 而原树的成为左节点

						return right
					}
				}
				// 否则就正常补充右节点
				binaryExpression.Right = right
				return binaryExpression
			} else {
				CoralErrorCrashHandler(NewCoralError("Compile Error",
					"expected a right node for binary expression!", ParsingUnexpected))
			}
		} else {
			return left // 即一个基本的表达式而已
		}
	}

	return nil
}
