package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
	"fmt"
)

func (parser *Parser) ParseStatement() Statement {
	if simpleStmt := parser.ParseSimpleStatement(); simpleStmt != nil {
		return simpleStmt
	}
	if breakStatement := parser.ParseBreakStatement(); breakStatement != nil {
		return breakStatement
	}
	if continueStatement := parser.ParseBreakStatement(); continueStatement != nil {
		return continueStatement
	}
	if returnStatement := parser.ParseReturnStatement(); returnStatement != nil {
		return returnStatement
	}
	if importStatement := parser.ParseImportStatement(); importStatement != nil {
		return importStatement
	}
	if enumStatement := parser.ParseEnumStatement(); enumStatement != nil {
		return enumStatement
	}

	return nil
}

func (parser *Parser) ParseSimpleStatement() SimpleStatement {
	if assignListStatement := parser.ParseAssignListStatement(); assignListStatement != nil {
		return assignListStatement
	}
	if expression := parser.ParseExpression(); expression != nil {
		parser.PeekNextToken()
		if parser.MatchCurrentTokenType(TokenTypeSemi) {
			return expression
		} else if parser.MatchCurrentTokenType(TokenTypeDoublePlus) || parser.MatchCurrentTokenType(TokenTypeDoubleMinus) {
			incDecStatement := new(IncDecStatement)
			incDecStatement.Expression = expression
			incDecStatement.Operator = parser.CurrentToken

			parser.PeekNextToken() // 移过 '++'/'--'
			parser.AssertCurrentTokenIs(TokenTypeSemi, "semicolon",
				"to terminate increase/decrease statement")
			return incDecStatement
		} else {
			CoralErrorCrashHandler(NewCoralError(parser.GetCurrentTokenPos(),
				"expected a semicolon to terminate this statement!", ParsingUnexpected))
		}
	}
	if varDeclStatement := parser.ParseVarDeclStatement(); varDeclStatement != nil {
		return varDeclStatement
	}

	return nil
}

func (parser *Parser) ParseVarDeclElement() *VarDeclElement {
	if parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		varNameToken := parser.CurrentToken
		varDeclElement := new(VarDeclElement)
		varDeclElement.VarName = varNameToken
		parser.PeekNextToken() // 移过当前 identifier

		if typeDescription := parser.ParseTypeDescription(); typeDescription != nil {
			// 有类型标注 -> 那么允许无初始值
			// eg (1): var b rune = 'B'
			// eg (2): var a int, s Student
			varDeclElement.Type = typeDescription
			// 当前 token 必须：要么是 '=' 要么是 ','
			if parser.MatchCurrentTokenType(TokenTypeEqual) {
				parser.PeekNextToken() // 移过 '='
				if initValue := parser.ParseExpression(); initValue != nil {
					varDeclElement.InitValue = initValue
					// 一个变量定义元素完成，此时 token 应为 ',' 会在外部循环断言
					return varDeclElement
				} else {
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
						fmt.Sprintf("expected an expression as initial value for variable '%s'", varNameToken.Str),
						ParsingUnexpected))
				}

			} else if parser.MatchCurrentTokenType(TokenTypeComma) {
				// 那么一个变量定义元素可以结束了，不移过逗号 ',' 而等待外部断言
				return varDeclElement
			} else {
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
					fmt.Sprintf("unexpected token: '%s', incomplete variable declaration!", parser.CurrentToken.Str),
					ParsingUnexpected))
			}
		} else {
			// 无类型标注 -> 则当前 token 必须是等号、有初始值
			if parser.MatchCurrentTokenType(TokenTypeEqual) {
				parser.PeekNextToken() // 移过 '='
				if initValue := parser.ParseExpression(); initValue != nil {
					varDeclElement.InitValue = initValue
					return varDeclElement
				} else {
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
						fmt.Sprintf("expected an expression as initial value to type inferring"+
							"for the non-typed variable '%s'!", varNameToken.Str),
						ParsingUnexpected))
				}
			} else {
				// 当前是其他不正确的 token
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
					fmt.Sprintf("unexpected token '%s' for variabel declaration!", parser.CurrentToken.Str),
					ParsingUnexpected))
			}
		}
	}

	return nil
}

func (parser *Parser) ParseVarDeclStatement() *VarDeclStatement {
	if parser.MatchCurrentTokenType(TokenTypeVar) || parser.MatchCurrentTokenType(TokenTypeVal) {
		varDeclStatement := new(VarDeclStatement)
		varDeclStatement.Mutable = parser.CurrentToken.Kind == TokenTypeVar
		parser.PeekNextToken() // 移过 'var'/'val'

		// 开始循环遍历读取 varDeclElement
		for varDeclElement := parser.ParseVarDeclElement(); varDeclElement != nil; varDeclElement = parser.ParseVarDeclElement() {
			varDeclStatement.Declarations = append(varDeclStatement.Declarations, varDeclElement)
			if parser.MatchCurrentTokenType(TokenTypeSemi) {
				// 分号即应该结束此段定义语句
				parser.PeekNextToken() // 移过 ';'
				return varDeclStatement
			} else {
				parser.AssertCurrentTokenIs(TokenTypeComma, "comma",
					"to separate multiple variable declarations!")
			}
		}
	}

	return nil
}

func (parser *Parser) ParseBreakStatement() *BreakStatement {
	if parser.MatchCurrentTokenType(TokenTypeBreak) {
		breakToken := parser.CurrentToken
		parser.PeekNextToken() // 移过 'break'

		parser.AssertCurrentTokenIs(TokenTypeSemi, "semicolon",
			"to terminate a break statement!")
		return &BreakStatement{Token: breakToken}
	}

	return nil
}

func (parser *Parser) ParseContinueStatement() *ContinueStatement {
	if parser.MatchCurrentTokenType(TokenTypeContinue) {
		continueToken := parser.CurrentToken
		parser.PeekNextToken() // 移过 'continue'

		parser.AssertCurrentTokenIs(TokenTypeSemi, "semicolon",
			"to terminate a continue statement!")
		return &ContinueStatement{Token: continueToken}
	}

	return nil
}

func (parser *Parser) ParseReturnStatement() *ReturnStatement {
	if parser.MatchCurrentTokenType(TokenTypeReturn) {
		returnToken := parser.CurrentToken
		parser.PeekNextToken() // 移过 'return'

		if expressionList := parser.ParseExpressionList(); expressionList != nil {
			return &ReturnStatement{
				Token:      returnToken,
				Expression: expressionList,
			}
		} else {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
				"expected an expression for return statement!", ParsingUnexpected))
		}
	}

	return nil
}

func (parser *Parser) ParseAssignListStatement() *AssignListStatement {
	if primaryExprList := parser.ParsePrimaryExpressionList(); primaryExprList != nil {
		assignListStatement := new(AssignListStatement)
		assignListStatement.Targets = primaryExprList

		if !parser.MatchCurrentTokenType(TokenTypeEqual) {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
				"expected a equal mark for assignment list!", ParsingUnexpected))
		}
		assignListStatement.Token = parser.CurrentToken
		parser.PeekNextToken() // 移过 '='

		if valueList := parser.ParseExpressionList(); valueList != nil {
			assignListStatement.Values = valueList
			parser.AssertCurrentTokenIs(TokenTypeSemi, "semicolon",
				"to terminate a assignment list!")
			return assignListStatement
		} else {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
				"expected a list of expression as values for assignment list!", ParsingUnexpected))
		}
	}

	return nil
}

func (parser *Parser) ParseImportElement() *ImportElement {
	if nameUnits := parser.ParseIdentifierList(); nameUnits != nil {
		moduleName := &ModuleName{NameUnits: nameUnits}
		importElement := &ImportElement{
			ModuleName: moduleName,
			As:         nil,
		}
		if parser.MatchCurrentTokenType(TokenTypeAs) {
			parser.PeekNextToken() // 移过 'as'
			if asName := parser.ParseIdentifier(); asName != nil {
				importElement.As = asName
				return importElement
			} else {
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
					"expected an identifier as target's another name for import statement!", ParsingUnexpected))
			}
		}

		return importElement
	}

	return nil
}

func (parser *Parser) ParseImportElementList() []*ImportElement {
	var elementList []*ImportElement
	for element := parser.ParseImportElement(); element != nil; element = parser.ParseImportElement() {
		elementList = append(elementList, element)
		if parser.MatchCurrentTokenType(TokenTypeComma) {
			parser.PeekNextToken() // 移过 ','
		} else {
			break
		}
	}

	if len(elementList) < 2 {
		return nil
	}
	return elementList
}

func (parser *Parser) ParseImportStatement() ImportStatement {
	if parser.MatchCurrentTokenType(TokenTypeFrom) {
		parser.PeekNextToken() // 移过 'from'
		if nameUnits := parser.ParseIdentifierList(); nameUnits != nil && len(nameUnits) > 0 {
			from := &ModuleName{NameUnits: nameUnits}

			if parser.MatchCurrentTokenType(TokenTypeImport) {
				parser.PeekNextToken() // 移过 'import'
				// 进入分支判断：
				if parser.MatchCurrentTokenType(TokenTypeLeftBrace) {
					// 当前的如果是大括号说明是 群组引入
					parser.PeekNextToken() // 移过 '{'
					if elementList := parser.ParseImportElementList(); elementList != nil {
						listImportStatement := &ListImportStatement{
							From:     from,
							Elements: elementList,
						}

						if parser.MatchCurrentTokenType(TokenTypeRightBrace) {
							parser.PeekNextToken() // 移过 '}'
							return listImportStatement
						} else {
							CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
								"expected a right brace as ending for a block import statement!", ParsingUnexpected))
						}
					} else {
						CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
							"expected at least two import element for a block import statement!", ParsingUnexpected))
					}
				} else if importElement := parser.ParseImportElement(); importElement != nil {
					singleImportStatement := new(SingleImportStatement)
					singleImportStatement.From = from
					singleImportStatement.Element = importElement
					if parser.MatchCurrentTokenType(TokenTypeSemi) {
						parser.PeekNextToken() // 移过 ';'
						return singleImportStatement
					} else {
						CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
							"expected a semicolon as ending for import statement!", ParsingUnexpected))
					}
				} else {
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
						"expected a module name as target for import statement!", ParsingUnexpected))
				}
			}
		} else {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
				"expected a module name as source for import statement!", ParsingUnexpected))
		}
	}

	return nil
}

func (parser *Parser) ParseEnumElement() *EnumElement {
	if parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		enumElement := new(EnumElement)
		enumElement.Name = &Identifier{Token: parser.CurrentToken}
		parser.PeekNextToken() // 移过当前这个名称标识符
		// 尝试解析等于号，看是否有赋值
		if parser.MatchCurrentTokenType(TokenTypeEqual) {
			parser.PeekNextToken() // 移过 '='
			if decimalLit, isDecimal := parser.ParseLiteral().(*DecimalLit); isDecimal {
				enumElement.Value = decimalLit
				return enumElement
			} else {
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
					"expected a decimal literal as the enum element's value!", ParsingUnexpected))
			}
		}

		if !parser.MatchCurrentTokenType(TokenTypeComma) && !parser.MatchCurrentTokenType(TokenTypeRightBrace) {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
				"expected a comma to separate multiple enum elements!", ParsingUnexpected))
		}
		return enumElement
	}

	return nil
}

func (parser *Parser) ParseEnumStatement() *EnumStatement {
	if parser.MatchCurrentTokenType(TokenTypeEnum) {
		parser.PeekNextToken() // 移过 'enum'
		if enumName := parser.ParseIdentifier(); enumName != nil {
			enumStatement := new(EnumStatement)
			enumStatement.Name = enumName

			if parser.MatchCurrentTokenType(TokenTypeLeftBrace) {
				parser.PeekNextToken() // 移过 '{'

				for enumElement := parser.ParseEnumElement(); enumElement != nil; enumElement = parser.ParseEnumElement() {
					enumStatement.Elements = append(enumStatement.Elements, enumElement)
					if parser.MatchCurrentTokenType(TokenTypeComma) {
						parser.PeekNextToken() // 移过 ','
					} else {
						break
					}
				}

				if parser.MatchCurrentTokenType(TokenTypeRightBrace) {
					parser.PeekNextToken() // 移过 '}'
					return enumStatement
				} else {
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile Error",
						"expected a right brace as ending for enum definition!", ParsingUnexpected))
				}
			}

		}
	}

	return nil
}
