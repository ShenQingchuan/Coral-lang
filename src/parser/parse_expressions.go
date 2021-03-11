package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
	"regexp"
)

func GetBinaryOperatorPriority(token *Token) int {
	if token == nil {
		return 99
	}

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

func (parser *Parser) ParseIdentifier(avoidAngleConfusingLater bool) *Identifier {
	if !parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		return nil
	}

	identifier := &Identifier{Token: parser.CurrentToken} // 以当前标识符为 operand
	if avoidAngleConfusingLater {
		parser.PeekNextTokenAvoidAngleConfusing()
	} else {
		parser.PeekNextToken()
	}
	return identifier
}

// IDENTIFIER ('.' IDENTIFIER)*
func (parser *Parser) ParseIdentifierList() []*Identifier {
	if !parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		return nil
	}

	// 先添加传入的 token，已确定其为 identifier
	var identifierList []*Identifier
	identifierList = append(identifierList, &Identifier{Token: parser.CurrentToken})

	for {
		parser.PeekNextToken()
		if !parser.MatchCurrentTokenType(TokenTypeDot) { // 需要一个 '.' 来分隔
			break
		}
		parser.PeekNextToken()                                  // 移过 '.'
		if !parser.MatchCurrentTokenType(TokenTypeIdentifier) { // '.' 后的 Identifier
			break
		}
		identifierList = append(identifierList, &Identifier{Token: parser.CurrentToken})
	}
	return identifierList
}

func (parser *Parser) ParseExpression() Expression {
	// 括号表达式优先级最高
	if parser.MatchCurrentTokenType(TokenTypeLeftParen) {
		currentLexerPos := parser.Lexer.BytePos
		tryLambdaLitExpression := parser.ParsePrimaryExpression() // 由于左圆括号的特殊性 先尝试解析 lambdaLit
		if tryLambdaLitExpression != nil {
			if _, isLambda := tryLambdaLitExpression.(*BasicPrimaryExpression).It.(*LambdaLit); isLambda {
				return parser.TryParseBinaryExpression(tryLambdaLitExpression)
			}
		} else {
			parser.Lexer.BytePos = currentLexerPos // 如果不是 lambda 恢复词法器位置
			parser.PeekNextToken()                 // 移过左括号
			inParenExpression := parser.ParseExpression()
			if !parser.AssertCurrentTokenIs(TokenTypeRightParen,
				"right parenthesis", "to close a parenthesis expression!") {
				return nil
			}
			return parser.TryParseBinaryExpression(inParenExpression)
		}
	}

	if unaryExpression := parser.ParseUnaryExpression(); unaryExpression != nil {
		return parser.TryParseBinaryExpression(unaryExpression)
	}
	if primaryExpr := parser.ParsePrimaryExpression(); primaryExpr != nil {
		return parser.TryParseBinaryExpression(primaryExpr)
	}
	if newInstanceExpression := parser.ParseNewInstanceExpression(); newInstanceExpression != nil {
		return parser.TryParseBinaryExpression(newInstanceExpression)
	}

	return nil
}

// 实质上是：解析基本表达式的 operand 部分
func (parser *Parser) ParsePrimaryExpression() PrimaryExpression {
	literal := parser.ParseLiteral()
	if literal != nil {
		return parser.TryEnhancePrimaryExpression(&BasicPrimaryExpression{It: literal})
	} // 如果 literal 为空则另一种情况

	operandName := parser.ParseOperandName()
	if operandName != nil {
		return parser.TryEnhancePrimaryExpression(&BasicPrimaryExpression{It: operandName})
	}

	return nil
}

// 探寻基本表达式的其他可能性 index/slice/call/member
func (parser *Parser) TryEnhancePrimaryExpression(basic PrimaryExpression) PrimaryExpression {
	_, isSlice := basic.(*SliceExpression)

	// try: slice/index
	if parser.MatchCurrentTokenType(TokenTypeLeftBracket) {
		parser.PeekNextToken()

		// 为什么先尝试 slice? : 因为可能存在 arr[:3]
		// 如果是中括号 '['
		// 如果确定是 '[:'
		if parser.MatchCurrentTokenType(TokenTypeColon) {
			sliceExpr := new(SliceExpression)
			sliceExpr.Operand = basic

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
				return parser.TryEnhancePrimaryExpression(sliceExpr)
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
			sliceExpr.Operand = basic
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
				return parser.TryEnhancePrimaryExpression(sliceExpr)
			} else {
				CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
					"expected an close bracket for slice expression!", ParsingUnexpected))
			}
		} else if parser.MatchCurrentTokenType(TokenTypeRightBracket) {
			// 只有一个表达式就遇到了右括号
			indexExpr := new(IndexExpression)
			indexExpr.Operand = basic
			indexExpr.Index = start

			parser.PeekNextToken()
			return parser.TryEnhancePrimaryExpression(indexExpr)
		}
	} else if parser.MatchCurrentTokenType(TokenTypeLeftParen) {
		// try: call
		if isSlice {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"slice expression can't be called/executed as a function!", ParsingUnexpected))
			return nil
		}
		parser.PeekNextToken() // 移过当前的左括号，到下一个 token
		callExpression := new(CallExpression)
		callExpression.Operand = basic
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
			return parser.TryEnhancePrimaryExpression(callExpression)
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a right parenthesis for function calling!", ParsingUnexpected))
			return nil
		}
	} else if parser.MatchCurrentTokenType(TokenTypeDot) {
		parser.PeekNextToken() // 移过 '.' 点，到下一个 token
		memberExpression := new(MemberExpression)
		memberExpression.Operand = basic
		if idList := parser.ParseIdentifierList(); idList != nil {
			memberExpression.Member = new(MemberLinkNode)
			cursor := memberExpression.Member // 开始根据得到的 标识符列表构建成员链
			for i, id := range idList {
				cursor.It = id
				if i != len(idList)-1 {
					cursor.MemberNext = new(MemberLinkNode) // 结链
					cursor = cursor.MemberNext              // -> next
				}
			}

			return parser.TryEnhancePrimaryExpression(memberExpression)
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a member identifier list after dot!", ParsingUnexpected))
			return nil
		}
	}

	return basic
}

func (parser *Parser) ParseTableElement() *TableElement {
	if parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		tableElement := new(TableElement)
		tableElement.Key = &Identifier{Token: parser.CurrentToken}
		parser.PeekNextToken() // 移过标识符

		if !parser.AssertCurrentTokenIs(TokenTypeColon, "a colon",
			"in map literal element to separate key and value!") {
			return nil
		}
		if value := parser.ParseExpression(); value != nil {
			tableElement.Value = value
			return tableElement
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected an expression as value in map literal element!", ParsingUnexpected))
			return nil
		}
	}

	return nil
}

// 解析 operand 的 literal 情况
func (parser *Parser) ParseLiteral() Literal {
	if parser.CurrentToken != nil {
		switch parser.CurrentToken.Kind {
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
			valueToken := parser.CurrentToken
			floatLit := new(FloatLit)
			floatLit.Value = valueToken
			if reg := regexp.MustCompile(`\.(\d+)$`); len(reg.FindString(parser.CurrentToken.Str))-1 > 6 && len(reg.FindString(parser.CurrentToken.Str))-1 <= 15 {
				floatLit.Accuracy = 15
			} else {
				floatLit.Accuracy = 6
			}
			return floatLit
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
		case TokenTypeLeftBracket:
			parser.PeekNextToken() // 移过 '['
			var expressionList []Expression
			for expression := parser.ParseExpression(); expression != nil; expression = parser.ParseExpression() {
				expressionList = append(expressionList, expression)

				if parser.MatchCurrentTokenType(TokenTypeComma) {
					parser.PeekNextToken()
				} else {
					break
				}
			}
			if !parser.AssertCurrentTokenIs(TokenTypeRightBracket, "right bracket",
				"to close the array literal value!") {
				return nil
			}
			return &ArrayLit{ValueList: expressionList}
		case TokenTypeLeftBrace:
			parser.PeekNextToken() // 移过 '{'
			var elements []*TableElement
			for tableElement := parser.ParseTableElement(); tableElement != nil; tableElement = parser.ParseTableElement() {
				elements = append(elements, tableElement)

				if parser.MatchCurrentTokenType(TokenTypeComma) {
					parser.PeekNextToken() // 移过 ','
				} else {
					break
				}
			}
			if !parser.AssertCurrentTokenIs(TokenTypeRightBrace, "right brace",
				"in map literal definition!") {
				return nil
			}
			return &TableLit{KeyValueList: elements}
		case TokenTypeLeftParen:
			if signature := parser.ParseSignature(true, true); signature != nil {
				lambdaLit := new(LambdaLit)
				lambdaLit.Signature = signature
				if parser.MatchCurrentTokenType(TokenTypeRightArrow) {
					parser.PeekNextToken() // 移过尖头
					if lambdaBlock := parser.ParseBlockStatement(); lambdaBlock != nil {
						lambdaLit.Result = lambdaBlock
						return lambdaLit
					} else if lambdaExpr := parser.ParseExpression(); lambdaExpr != nil {
						lambdaLit.Result = lambdaExpr
						return lambdaLit
					} else {
						CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
							"expected a block for lambda lambda function!", ParsingUnexpected))
						return nil
					}
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected a right arrow '->' for lambda function!", ParsingUnexpected))
					return nil
				}
			}
			return nil
		case TokenTypeThis:
			defer parser.PeekNextToken()
			return &ThisLit{
				Token:     parser.CurrentToken,
				BelongsTo: nil,
			}
		case TokenTypeSuper:
			defer parser.PeekNextToken()
			return &SuperLit{
				Token:     parser.CurrentToken,
				BelongsTo: nil,
			}
		}
	}

	return nil
}

// 解析 operand 的 operandName（变量名）情况
func (parser *Parser) ParseOperandName() *OperandName {
	if !parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		return nil
	}

	// 先添加传入的 token，已确定其为 identifier
	operandName := new(OperandName)
	identifier := parser.ParseIdentifier(false)
	operandName.Name = identifier

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
		CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
			"expected a type description for object instance creating!", ParsingUnexpected))
		return nil
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
				if !parser.AssertCurrentTokenIs(TokenTypeComma,
					"a comma", "in new object instance constructor") {
					return nil
				}
			}
		}
		// 结束循环时，检测是否停留于 token ')'
		if parser.MatchCurrentTokenType(TokenTypeRightParen) {
			parser.PeekNextToken() // 移过 ')'
			return newInstanceExpression
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a right parenthesis for constructor method's ending!", ParsingUnexpected))
			return nil
		}
	} else {
		CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
			"expected a left parenthesis for constructor method's calling!", ParsingUnexpected))
		return nil
	}
}

// 解析 单目表达式
func (parser *Parser) ParseUnaryExpression() *UnaryExpression {
	if parser.CurrentToken != nil {
		switch parser.CurrentToken.Kind {
		case TokenTypeMinus, TokenTypeBang, TokenTypeWavy:
			unaryExpression := new(UnaryExpression)
			unaryExpression.Operator = parser.CurrentToken
			parser.PeekNextToken() // 移过该单目运算符

			if operand := parser.ParsePrimaryExpression(); operand != nil {
				unaryExpression.Operand = operand
				return unaryExpression
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"missing operand for unary expression!", ParsingUnexpected))
				return nil
			}
		}
	}

	return nil
}

// 递归尝试解析 二元表达式
func (parser *Parser) TryParseBinaryExpression(left Expression) Expression {
	if parser.MatchCurrentTokenType(TokenTypeAs) {
		parser.PeekNextToken() // 移过 'as'
		castExpression := new(CastExpression)
		castExpression.Source = left

		if typeDescription := parser.ParseTypeDescription(); typeDescription != nil {
			castExpression.Type = typeDescription
			return castExpression
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a type name for typing cast!", ParsingUnexpected))
			return nil
		}
	} else if leftPrimary, isPrimary := left.(PrimaryExpression); // 针对 区间表达式 的特判
	isPrimary && (parser.MatchCurrentTokenType(TokenTypeEllipsis) || parser.MatchCurrentTokenType(TokenTypeDoubleDot)) {
		rangeExpression := new(RangeExpression)
		rangeExpression.Start = leftPrimary
		rangeExpression.IncludeEnd = parser.CurrentToken.Kind == TokenTypeEllipsis // 三点表示闭区间，包括终点
		parser.PeekNextToken()                                                     // 移动过 三点或两点 符号
		if right := parser.ParsePrimaryExpression(); right != nil {
			rangeExpression.End = right
			return rangeExpression // 区间表达式比较独立，不需要再额外操作
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a primary expression as range endpoint!", ParsingUnexpected))
			return nil
		}
	} else if priority := GetBinaryOperatorPriority(parser.CurrentToken); priority != 99 {
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
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a right node for binary expression!", ParsingUnexpected))
			return nil
		}
	}
	return left // 即一个基本的表达式而已
}

// 解析一个以逗号分隔的 表达式 列表
func (parser *Parser) ParseExpressionList() []Expression {
	var exprList []Expression
	for expression := parser.ParseExpression(); expression != nil; expression = parser.ParseExpression() {
		exprList = append(exprList, expression)
		if parser.MatchCurrentTokenType(TokenTypeComma) {
			parser.PeekNextToken() // 移过 ',' 逗号, continue
		} else {
			return exprList
		}
	}

	return nil
}

// 解析一个以逗号分隔 基本表达式 列表
func (parser *Parser) ParsePrimaryExpressionList() []PrimaryExpression {
	var primaryExprList []PrimaryExpression
	for primaryExpr := parser.ParsePrimaryExpression(); primaryExpr != nil; primaryExpr = parser.ParsePrimaryExpression() {
		primaryExprList = append(primaryExprList, primaryExpr)
		if parser.MatchCurrentTokenType(TokenTypeComma) {
			parser.PeekNextToken() // 移过 ',' 逗号, continue
		} else {
			return primaryExprList
		}
	}

	return nil
}
