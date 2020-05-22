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
func (parser *Parser) AssertCurrentIsComma(situation string) {
	if parser.MatchCurrentTokenType(TokenTypeComma) {
		parser.PeekNextToken()
	} else {
		CoralErrorCrashHandler(NewCoralError("Compile Error",
			fmt.Sprintf("expected a comma for expression list %s!", situation), ParsingUnexpected))
	}
}

func GetBinaryOperatorPriority(token *Token) int {
	if token == nil {
		return 99
	}

	switch token.Kind {
	case TokenTypeLeftParen, TokenTypeLeftBracket:
		return 1
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
func (parser *Parser) MatchCurrentTokenType(tokenType TokenType) bool {
	if parser.CurrentToken != nil {
		return parser.CurrentToken.Kind == tokenType
	}
	return false
}

// IDENTIFIER ('.' IDENTIFIER)*
func (parser *Parser) ParseIdentifierList() []*Identifier {
	if !parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		return nil
	}

	// 先添加传入的 token，已确定其为 identifier
	var identifierList []*Identifier
	identifierList = append(identifierList, &Identifier{Name: parser.CurrentToken})

	for {
		parser.PeekNextToken()
		if !parser.MatchCurrentTokenType(TokenTypeDot) { // 需要一个 '.' 来分隔
			break
		}
		parser.PeekNextToken()
		if !parser.MatchCurrentTokenType(TokenTypeIdentifier) { // '.' 后的 Identifier
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
		if parser.MatchCurrentTokenType(TokenTypeSemi) {
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
		return parser.TryParseBinaryExpression(primaryExpr)
	}
	if newInstanceExpression := parser.ParseNewInstanceExpression(); newInstanceExpression != nil {
		return parser.TryParseBinaryExpression(newInstanceExpression)
	}
	if unaryExpression := parser.ParseUnaryExpression(); unaryExpression != nil {
		return parser.TryParseBinaryExpression(unaryExpression)
	}

	return nil
}

func (parser *Parser) ParseTypeDescription() TypeDescription {
	// 如果是左中括号
	if parser.MatchCurrentTokenType(TokenTypeLeftBracket) {
		parser.PeekNextToken()
		typeDescription := parser.ParseTypeDescription()
		if typeDescription != nil {
			arrayTypeLit := new(ArrayTypeLit)
			arrayTypeLit.ElementType = &typeDescription

			// 此时当前 token 应为 ']'
			if parser.MatchCurrentTokenType(TokenTypeRightBracket) {
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
	if parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		typeName := parser.ParseTypeName()
		if typeName != nil {
			// 结束 typeName 解析后如果是一个 '<' 则进入 Generics 解析
			if parser.MatchCurrentTokenType(TokenTypeLeftAngle) {
				genericsTypeLit := new(GenericsTypeLit)
				genericsTypeLit.BasicType = typeName

				parser.PeekNextToken() // 越过 '<'
				for {
					genericsArg := parser.ParseTypeName()
					if genericsArg != nil {
						genericsTypeLit.GenericsArgs = append(genericsTypeLit.GenericsArgs, genericsArg)
					}
					if parser.MatchCurrentTokenType(TokenTypeRightAngle) {
						parser.PeekNextToken() // 移过 '>'
						return genericsTypeLit // 结束泛型参数解析
					} else {
						parser.AssertCurrentIsComma(fmt.Sprintf(
							"for seperating generics arguments but got '%s'",
							parser.CurrentToken.Str))
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
	if !parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		return nil
	}

	// 先添加传入的 token，已确定其为 identifier
	typeName := new(TypeName)
	identifierList := parser.ParseIdentifierList()
	typeName.NameList = identifierList

	return typeName
}

// 实质上是：解析基本表达式的 operand 部分
func (parser *Parser) ParsePrimaryExpression() PrimaryExpression {
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
	case TokenTypeString:
		defer parser.PeekNextToken()
		return &StringLit{Value: parser.CurrentToken}
	case TokenTypeRune:
		defer parser.PeekNextToken()
		return &RuneLit{Value: parser.CurrentToken}
	case TokenTypeDecimalInteger:
		defer parser.PeekNextToken()
		return &DecimalLit{Value: parser.CurrentToken}
	case TokenTypeHexadecimalInteger:
		defer parser.PeekNextToken()
		return &HexadecimalLit{Value: parser.CurrentToken}
	case TokenTypeOctalInteger:
		defer parser.PeekNextToken()
		return &OctalLit{Value: parser.CurrentToken}
	case TokenTypeBinaryInteger:
		defer parser.PeekNextToken()
		return &BinaryLit{Value: parser.CurrentToken}
	case TokenTypeFloat:
		defer parser.PeekNextToken()
		return &FloatLit{Value: parser.CurrentToken}
	case TokenTypeExponent:
		defer parser.PeekNextToken()
		return &ExponentLit{Value: parser.CurrentToken}
	case TokenTypeNil:
		defer parser.PeekNextToken()
		return &NilLit{Value: parser.CurrentToken}
	case TokenTypeTrue:
		defer parser.PeekNextToken()
		return &TrueLit{Value: parser.CurrentToken}
	case TokenTypeFalse:
		defer parser.PeekNextToken()
		return &FalseLit{Value: parser.CurrentToken}
	}
}

// 解析 operand 的 operandName（变量名）情况
func (parser *Parser) ParseOperandName() *OperandName {
	if !parser.MatchCurrentTokenType(TokenTypeIdentifier) {
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
	if !parser.MatchCurrentTokenType(TokenTypeNew) {
		return nil
	}

	parser.PeekNextToken()
	typeDescription := parser.ParseTypeDescription()
	if typeDescription == nil {
		CoralErrorCrashHandler(NewCoralError("Compile Error",
			"expected a type description for object instance creating!", ParsingUnexpected))
	}
	if parser.MatchCurrentTokenType(TokenTypeLeftParen) {
		parser.PeekNextToken() // 移过当前的左括号，到下一个 token
		newInstanceExpression := new(NewInstanceExpression)
		newInstanceExpression.Class = typeDescription
		for expression := parser.ParseExpression(); expression != nil; expression = parser.ParseExpression() {
			newInstanceExpression.InitParams = append(newInstanceExpression.InitParams, expression)
			if parser.MatchCurrentTokenType(TokenTypeRightParen) {
				break
			} else {
				parser.AssertCurrentIsComma("in new object instance constructor")
			}
		}
		// 结束循环时，检测是否停留于 token ')'
		if parser.MatchCurrentTokenType(TokenTypeRightParen) {
			parser.PeekNextToken() // 移过 ')'
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

func (parser *Parser) TryParseBinaryExpression(left Expression) Expression {
	if priority := GetBinaryOperatorPriority(parser.CurrentToken); priority != 99 {
		if priority == 1 {
			_, isLeftSlice := left.(*SliceExpression)

			// try: slice/index
			if parser.MatchCurrentTokenType(TokenTypeLeftBracket) {
				parser.PeekNextToken()

				// 为什么先尝试 slice? : 因为可能存在 arr[:3]
				// 如果是中括号 '['
				// 如果确定是 '[:'
				if parser.MatchCurrentTokenType(TokenTypeColon) {
					sliceExpr := new(SliceExpression)
					sliceExpr.Operand = left

					// 解析切片终点表达式
					parser.PeekNextToken()
					end := parser.ParseExpression()
					if end == nil {
						CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
							"expected an expression to be end position for slice!", ParsingUnexpected))
					}
					sliceExpr.End = end

					if parser.MatchCurrentTokenType(TokenTypeRightBracket) {
						parser.PeekNextToken() // 移过 ']'
						return parser.TryParseBinaryExpression(sliceExpr)
					} else {
						CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
							"expected an close bracket for slice expression!", ParsingUnexpected))
					}
				}

				start := parser.ParseExpression()
				if start == nil {
					CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
						"expected an expression to be an index/key or a start position for slice!", ParsingUnexpected))
				}

				if parser.MatchCurrentTokenType(TokenTypeColon) {
					sliceExpr := new(SliceExpression)
					sliceExpr.Operand = left
					sliceExpr.Start = start

					// 解析切片终点表达式
					parser.PeekNextToken() // 移过冒号 ':'
					end := parser.ParseExpression()
					if end == nil {
						CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
							"expected an expression to be an end position for slice!", ParsingUnexpected))
					}
					sliceExpr.End = end

					if parser.MatchCurrentTokenType(TokenTypeRightBracket) {
						parser.PeekNextToken()
						return parser.TryParseBinaryExpression(sliceExpr)
					} else {
						CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
							"expected an close bracket for slice expression!", ParsingUnexpected))
					}
				} else if parser.MatchCurrentTokenType(TokenTypeRightBracket) {
					// 只有一个表达式就遇到了右括号
					indexExpr := new(IndexExpression)
					indexExpr.Operand = left
					indexExpr.Index = start

					parser.PeekNextToken()
					return parser.TryParseBinaryExpression(indexExpr)
				}
			} else if parser.MatchCurrentTokenType(TokenTypeLeftParen) {
				if isLeftSlice {
					CoralErrorCrashHandler(NewCoralError("Compile Error",
						"slice expression can't be called/executed as a function!", ParsingUnexpected))
				}
				parser.PeekNextToken() // 移过当前的左括号，到下一个 token
				callExpression := new(CallExpression)
				callExpression.Operand = left
				for expression := parser.ParseExpression(); expression != nil; expression = parser.ParseExpression() {
					callExpression.Params = append(callExpression.Params, expression)

					if parser.MatchCurrentTokenType(TokenTypeRightParen) {
						break // 虽然我这里定义了 break 的条件是需要当前 token 为右括号，但是可能语法有错误，例如根本没写右括号
					} else if parser.MatchCurrentTokenType(TokenTypeComma) {
						parser.PeekNextToken()
					}
				}
				// 结束循环时，检测是否停留于 token ')'
				if parser.MatchCurrentTokenType(TokenTypeRightParen) {
					parser.PeekNextToken()
					return parser.TryParseBinaryExpression(callExpression)
				} else {
					CoralErrorCrashHandler(NewCoralError("Compile Error",
						"expected a right parenthesis for function calling!", ParsingUnexpected))
				}
			}
		} else {
			// 证明是 二元运算符
			binaryExpression := new(BinaryExpression)
			binaryExpression.Operator = parser.CurrentToken
			binaryExpression.Left = left

			parser.PeekNextToken() // 移过当前操作符节点

			// 先暂时形成 左结构，然后解析右边节点（即体现：自左向右结合）
			if r := parser.ParseExpression(); r != nil {
				right, rightIsBinary := r.(*BinaryExpression)
				if rightIsBinary && GetBinaryOperatorPriority(right.Operator) > priority {
					// 右边节点也是二元表达式，需要进行优先级判断、旋转树
					// 即 右边节点的 priority 数值大，反而优先级别低，应作父节点
					binaryExpression.Right = right.Left // 补充原树的右节点
					right.Left = binaryExpression       // 而原树的成为左节点

					return right
				}
				// 否则就正常补充右节点
				binaryExpression.Right = r
				return binaryExpression
			} else {
				CoralErrorCrashHandler(NewCoralError("Compile Error",
					"expected a right node for binary expression!", ParsingUnexpected))
			}
		}
	}
	return left // 即一个基本的表达式而已
}
