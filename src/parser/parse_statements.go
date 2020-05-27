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
	if blockStatement := parser.ParseBlockStatement(); blockStatement != nil {
		return blockStatement
	}
	if ifStatement := parser.ParseIfStatement(); ifStatement != nil {
		return ifStatement
	}
	if switchStatement := parser.ParseSwitchStatement(); switchStatement != nil {
		return switchStatement
	}

	return nil
}

func (parser *Parser) ParseSimpleStatement() SimpleStatement {
	if expression := parser.ParseExpression(); expression != nil {
		if parser.MatchCurrentTokenType(TokenTypeSemi) {
			parser.PeekNextToken() // 移过分号 ';'
			return &ExpressionStatement{Expression: expression}
		} else if expression, isPrimary := expression.(PrimaryExpression); isPrimary && parser.MatchCurrentTokenType(TokenTypeComma) {
			parser.PeekNextToken() // 移过 ','
			primaryExprList := []PrimaryExpression{expression}
			for primaryExpr := parser.ParsePrimaryExpression(); primaryExpr != nil; primaryExpr = parser.ParsePrimaryExpression() {
				primaryExprList = append(primaryExprList, primaryExpr)
				if parser.MatchCurrentTokenType(TokenTypeComma) {
					parser.PeekNextToken() // 移过 ',' 逗号, continue
				} else {
					break // primaryExpressionList 收集完毕
				}
			}
			assignListStatement := new(AssignListStatement)
			assignListStatement.Targets = primaryExprList

			if !parser.MatchCurrentTokenType(TokenTypeEqual) {
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
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
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
					"expected a list of expression as values for assignment list!", ParsingUnexpected))
			}

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
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
						fmt.Sprintf("expected an expression as initial value for variable '%s'", varNameToken.Str),
						ParsingUnexpected))
				}

			} else if parser.MatchCurrentTokenType(TokenTypeComma) {
				// 那么一个变量定义元素可以结束了，不移过逗号 ',' 而等待外部断言
				return varDeclElement
			} else {
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
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
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
						fmt.Sprintf("expected an expression as initial value to type inferring"+
							"for the non-typed variable '%s'!", varNameToken.Str),
						ParsingUnexpected))
				}
			} else {
				// 当前是其他不正确的 token
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
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
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
				"expected an expression for return statement!", ParsingUnexpected))
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
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
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
							CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
								"expected a right brace as ending for a block import statement!", ParsingUnexpected))
						}
					} else {
						CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
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
						CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
							"expected a semicolon as ending for import statement!", ParsingUnexpected))
					}
				} else {
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
						"expected a module name as target for import statement!", ParsingUnexpected))
				}
			}
		} else {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
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
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
					"expected a decimal literal as the enum element's value!", ParsingUnexpected))
			}
		}

		if !parser.MatchCurrentTokenType(TokenTypeComma) && !parser.MatchCurrentTokenType(TokenTypeRightBrace) {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
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
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
						"expected a right brace as ending for enum definition!", ParsingUnexpected))
				}
			}

		}
	}

	return nil
}

func (parser *Parser) ParseBlockStatement() *BlockStatement {
	if parser.MatchCurrentTokenType(TokenTypeLeftBrace) {
		parser.PeekNextToken() // 移过 '{'

		blockStatement := new(BlockStatement)
		for stmt := parser.ParseStatement(); stmt != nil; stmt = parser.ParseStatement() {
			blockStatement.Statements = append(blockStatement.Statements, stmt)
			if parser.MatchCurrentTokenType(TokenTypeRightBrace) {
				parser.PeekNextToken()
				return blockStatement
			}
		}

		// 能结束循环到此处说明有问题、没有正常解析到右括号
		CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
			"expected a right brace as ending for block statement!", ParsingUnexpected))
	}

	return nil
}

func (parser *Parser) ParseIfElement() *IfElement {
	if condition := parser.ParseExpression(); condition != nil {
		ifElement := new(IfElement)
		ifElement.Condition = condition

		if block := parser.ParseBlockStatement(); block != nil {
			ifElement.Block = block
			return ifElement
		} else {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
				`expected a block as a "if" block for "if" statement!`, ParsingUnexpected))
		}
	} else {
		CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
			`expected an expression as a condition for "if" statement!`, ParsingUnexpected))
	}

	return nil
}

func (parser *Parser) ParseIfStatement() *IfStatement {
	if parser.MatchCurrentTokenType(TokenTypeIf) {
		parser.PeekNextToken() // 移过 'if'
		ifStatement := new(IfStatement)
		if ifElement := parser.ParseIfElement(); ifElement != nil {
			ifStatement.If = ifElement

			// 解析可能存在的 一些 elif
			if elifElements := parser.ParseElifStatements(); elifElements != nil {
				ifStatement.Elif = elifElements
			}

			// 解析可能存在的 else
			if parser.MatchCurrentTokenType(TokenTypeElse) {
				parser.PeekNextToken() // 移过 'else'
				if elseBlock := parser.ParseBlockStatement(); elseBlock != nil {
					ifStatement.Else = elseBlock
				} else {
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
						`expected a block statement for "else" statement!`, ParsingUnexpected))
				}
			}

			return ifStatement
		}
	}

	return nil
}

func (parser *Parser) ParseElifStatements() []*IfElement {
	var elifElements []*IfElement
	for {
		if parser.MatchCurrentTokenType(TokenTypeElif) {
			parser.PeekNextToken() // 移过 'elif'
			if elifElement := parser.ParseIfElement(); elifElement != nil {
				elifElements = append(elifElements, elifElement)
			} else {
				CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
					`expected an condition and block statement for "elif" statement!`, ParsingUnexpected))
			}
		} else {
			break
		}
	}

	if len(elifElements) == 0 {
		return nil
	}
	return elifElements
}

func (parser *Parser) ParseSwitchCase() (SwitchStatementCase, bool) {
	if parser.MatchCurrentTokenType(TokenTypeCase) {
		parser.PeekNextToken() // 移过 'case'
		if caseExpr := parser.ParseExpression(); caseExpr != nil {
			rangeExpr, isRange := caseExpr.(*RangeExpression)
			if isRange {
				rangeCase := new(SwitchStatementRangeCase)
				rangeCase.Range = rangeExpr
				if block := parser.ParseBlockStatement(); block != nil {
					rangeCase.Block = block
					return rangeCase, false
				} else {
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
						"expected a block statement as case handler!", ParsingUnexpected))
				}
			} else {
				normalCase := new(SwitchStatementNormalCase)
				normalCase.Conditions = append(normalCase.Conditions, caseExpr)
				// 要根据逗号情况：
				if parser.MatchCurrentTokenType(TokenTypeComma) {
					parser.PeekNextToken() // 移过 ','
					for condition := parser.ParseExpression(); condition != nil; condition = parser.ParseExpression() {
						normalCase.Conditions = append(normalCase.Conditions, condition)
						if parser.MatchCurrentTokenType(TokenTypeComma) {
							parser.PeekNextToken() // 移过 ','
						} else {
							break
						}
					}
					if normalBlock := parser.ParseBlockStatement(); normalBlock != nil {
						normalCase.Block = normalBlock
						return normalCase, false
					} else {
						CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
							"expected a block statement as case handler!", ParsingUnexpected))
					}
				} else {
					if caseBlock := parser.ParseBlockStatement(); caseBlock != nil {
						normalCase.Block = caseBlock
						return normalCase, false
					} else {
						CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
							"expected a block statement as case handler!", ParsingUnexpected))
					}
				}
			}
		} else {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
				"expected an expression as a case!", ParsingUnexpected))
		}
	} else if parser.MatchCurrentTokenType(TokenTypeDefault) {
		parser.PeekNextToken() // 移过 'default'
		return nil, true
	}

	return nil, false
}

func (parser *Parser) ParseSwitchStatement() *SwitchStatement {
	if parser.MatchCurrentTokenType(TokenTypeSwitch) {
		parser.PeekNextToken() // 移过 'switch'

		if entry := parser.ParseExpression(); entry != nil {
			switchStatement := new(SwitchStatement)
			switchStatement.Entry = entry

			if parser.MatchCurrentTokenType(TokenTypeLeftBrace) {
				parser.PeekNextToken() // 移过 '{'

				for _case, isDefault := parser.ParseSwitchCase(); _case != nil || isDefault; _case, isDefault = parser.ParseSwitchCase() {
					if isDefault {
						if defaultBlock := parser.ParseBlockStatement(); defaultBlock != nil {
							switchStatement.Default = defaultBlock
						} else {
							CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
								`expected a block after 'default' keyword in switch statement!`, ParsingUnexpected))
						}
					} else {
						switchStatement.Cases = append(switchStatement.Cases, _case)
					}
				}

				if parser.MatchCurrentTokenType(TokenTypeRightBrace) {
					parser.PeekNextToken() // 移过 '}'
					return switchStatement
				} else {
					CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
						"expected a right brace as ending for switch statement!", ParsingUnexpected))
				}
			}

		} else {
			CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
				"expected an expression as target for switch statement!", ParsingUnexpected))
		}
	}

	return nil
}