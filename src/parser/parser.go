package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
)

type Parser struct {
	Lexer   *Lexer
	Program *Program
}

func (parser *Parser) InitParserFromBytes(content []byte) {
	InitLexerFromBytes(parser.Lexer, content)
}
func (parser *Parser) InitParserFromString(content string) {
	InitLexerFromString(parser.Lexer, content)
}

func (parser *Parser) ParseProgram() {
	for stmt := parser.ParseStatement(); stmt != nil; stmt = parser.ParseStatement() {
		// stmt 为 nil 的情况中其实早已被 CoralErrorCrashHandler 处理并退出了
		parser.Program.Root = append(parser.Program.Root, stmt)
	}
}

func (parser *Parser) ParseStatement() Statement {
	currentToken, err := parser.Lexer.GetNextToken()
	if err != nil {
		CoralErrorCrashHandler(err)
	}

	if simpleStmt := parser.ParseSimpleStatement(currentToken); simpleStmt != nil {
		return simpleStmt
	}

	return nil
}

func (parser *Parser) ParseSimpleStatement(currentToken *Token) SimpleStatement {
	expression := parser.ParseExpression(currentToken)
	if expression != nil {
		return expression
	}

	// TODO: simpleStatement 的其他情况

	return nil
}

func (parser *Parser) ParseExpression(currentToken *Token) Expression {
	if primaryExpr := parser.ParsePrimaryExpression(currentToken); primaryExpr != nil {
		return primaryExpr
	}

	return nil
}

func (parser *Parser) ParseTypeDescription(currentToken *Token) TypeDescription {
	// TODO：解析 类型标注

	return nil
}

func (parser *Parser) ParsePrimaryExpression(currentToken *Token) PrimaryExpression {
	operand := parser.ParseBasicPrimaryExpression(currentToken)
	if operand == nil {
		return nil // 没有解析到 operand 说明不是
	}
	// 至此已经获取到了 operand，后续 三种情况如果解析中都不能匹配
	// 则最后返回 operand

	// try: slice/index
	// 为什么先尝试 slice? : 因为可能存在 arr[:3]
	currentToken, err := parser.Lexer.GetNextToken()
	if err != nil {
		CoralErrorCrashHandler(err)
	}

	// 如果是中括号 '['
	if currentToken.Kind == TokenTypeLeftBracket {
		currentToken, err = parser.Lexer.GetNextToken()
		if err != nil {
			CoralErrorCrashHandler(err)
		}

		// 如果确定是 '[:'
		if currentToken.Kind == TokenTypeColon {
			sliceExpr := new(SliceExpression)
			sliceExpr.Operand = operand

			// 解析切片终点表达式
			currentToken, err = parser.Lexer.GetNextToken()
			if err != nil {
				CoralErrorCrashHandler(err)
			}
			end := parser.ParseExpression(currentToken)
			if end == nil {
				CoralErrorCrashHandler(NewCoralError("Parsing",
					"expected an expression to be end position for slice!", ParsingUnexpected))
			}
			sliceExpr.End = &end
			return sliceExpr
		}

		start := parser.ParseExpression(currentToken)
		if start == nil {
			CoralErrorCrashHandler(NewCoralError("Parsing",
				"expected an expression to be an index/key or a start position for slice!", ParsingUnexpected))
		}

		currentToken, err = parser.Lexer.GetNextToken()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		if currentToken.Kind == TokenTypeColon {
			sliceExpr := new(SliceExpression)
			sliceExpr.Operand = operand
			sliceExpr.Start = &start

			// 解析切片终点表达式
			currentToken, err = parser.Lexer.GetNextToken()
			if err != nil {
				CoralErrorCrashHandler(err)
			}
			end := parser.ParseExpression(currentToken)
			if end == nil {
				CoralErrorCrashHandler(NewCoralError("Parsing",
					"expected an expression to be an end position for slice!", ParsingUnexpected))
			}
			sliceExpr.End = &end
			return sliceExpr
		} else if currentToken.Kind == TokenTypeRightBracket {
			// 只有一个表达式就遇到了右括号
			indexExpr := new(IndexExpression)
			indexExpr.Operand = operand
			indexExpr.Index = &start
			return indexExpr
		}
	}

	return operand
}

// 解析基本表达式的 operand 部分
func (parser *Parser) ParseBasicPrimaryExpression(currentToken *Token) *BasicPrimaryExpression {
	literal := parser.ParseLiteral(currentToken)
	if literal != nil {
		return &BasicPrimaryExpression{Operand: literal}
	} // 如果 literal 为空则另一种情况

	operandName := parser.ParseOperandName(currentToken)
	if operandName != nil {
		return &BasicPrimaryExpression{Operand: operandName}
	}

	return nil
}

// 解析 operand 的 literal 情况
func (parser *Parser) ParseLiteral(currentToken *Token) Literal {
	switch currentToken.Kind {
	default:
		return nil
	case TokenTypeDecimalInteger:
		return &DecimalLit{Value: currentToken}
	case TokenTypeHexadecimalInteger:
		return &HexadecimalLit{Value: currentToken}
	case TokenTypeOctalInteger:
		return &OctalLit{Value: currentToken}
	case TokenTypeBinaryInteger:
		return &BinaryLit{Value: currentToken}
	case TokenTypeNil:
		return &NilLit{Value: currentToken}
	case TokenTypeTrue:
		return &TrueLit{Value: currentToken}
	case TokenTypeFalse:
		return &FalseLit{Value: currentToken}
	}
}

// 解析 operand 的 operandName（变量名）情况
func (parser *Parser) ParseOperandName(currentToken *Token) *OperandName {
	if currentToken.Kind != TokenTypeIdentifier {
		return nil
	}

	// 先添加传入的 token，已确定其为 identifier
	operandName := new(OperandName)
	operandName.NameList = append(operandName.NameList, &Identifier{Name: currentToken})

	for {
		currentToken, err := parser.Lexer.GetNextToken()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		if currentToken.Kind != TokenTypeDot { // 需要一个 '.' 来分隔
			break
		}
		currentToken, err = parser.Lexer.GetNextToken()
		if err != nil {
			CoralErrorCrashHandler(err)
		}
		if currentToken.Kind != TokenTypeIdentifier { // '.' 后的 Identifier
			break
		}
		operandName.NameList = append(operandName.NameList, &Identifier{Name: currentToken})
	}
	return operandName
}
