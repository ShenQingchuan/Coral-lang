package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
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

func (parser *Parser) ParseIdentifier() *Identifier {
	if !parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		return nil
	}

	identifier := new(Identifier)
	identifier.Token = parser.CurrentToken // 以当前标识符为 operand
	parser.PeekNextToken()
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
		inParenExpression := parser.ParseExpression()
		parser.AssertCurrentTokenIs(TokenTypeRightParen,
			"right parenthesis", "to close a parenthesis expression!")
		return inParenExpression
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
		return parser.TryEnhancePrimaryExpression(&BasicPrimaryExpression{Operand: literal})
	} // 如果 literal 为空则另一种情况

	operandName := parser.ParseOperandName()
	if operandName != nil {
		return parser.TryEnhancePrimaryExpression(&BasicPrimaryExpression{Operand: operandName})
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
			CoralErrorCrashHandler(NewCoralError("Compile Error",
				"slice expression can't be called/executed as a function!", ParsingUnexpected))
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
			CoralErrorCrashHandler(NewCoralError("Compile Error",
				"expected a right parenthesis for function calling!", ParsingUnexpected))
		}
	} else if parser.MatchCurrentTokenType(TokenTypeDot) {
		parser.PeekNextToken() // 移过 '.' 点，到下一个 token
		memberExpression := new(MemberExpression)
		memberExpression.Operand = basic
		if idList := parser.ParseIdentifierList(); idList != nil {
			memberExpression.Member = new(MemberLinkNode)
			cursor := memberExpression.Member // 开始根据得到的 标识符列表构建成员链
			for _, id := range idList {
				cursor.Operand = id
				cursor.MemberNext = new(MemberLinkNode) // 结链
				cursor = cursor.MemberNext              // -> next
			}

			return parser.TryEnhancePrimaryExpression(memberExpression)
		} else {
			CoralErrorCrashHandler(NewCoralError("Compile Error",
				"expected a member identifier list after dot!", ParsingUnexpected))
		}
	}

	return basic
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
	identifier := parser.ParseIdentifier()
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
				parser.AssertCurrentTokenIs(TokenTypeComma, "comma", "in new object instance constructor")
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
		parser.PeekNextToken() // 移过该单目运算符

		if operand := parser.ParsePrimaryExpression(); operand != nil {
			unaryExpression.Operand = operand
			return unaryExpression
		} else {
			CoralErrorCrashHandler(NewCoralError("Compile Error",
				"missing operand for unary expression!", ParsingUnexpected))
		}
	}

	return nil
}

// 递归尝试解析 二元表达式
func (parser *Parser) TryParseBinaryExpression(left Expression) Expression {
	if priority := GetBinaryOperatorPriority(parser.CurrentToken); priority != 99 {
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
