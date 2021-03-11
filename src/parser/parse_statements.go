package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
	"fmt"
)

func (parser *Parser) ParseStatement() Statement {
	if simpleStmt := parser.ParseSimpleStatement(true); simpleStmt != nil {
		return simpleStmt
	}
	if breakStatement := parser.ParseBreakStatement(); breakStatement != nil {
		return breakStatement
	}
	if continueStatement := parser.ParseContinueStatement(); continueStatement != nil {
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
	if whileStatement := parser.ParseWhileStatement(); whileStatement != nil {
		return whileStatement
	}
	if forStatement := parser.ParseForStatement(); forStatement != nil {
		return forStatement
	}
	if eachStatement := parser.ParseEachStatement(); eachStatement != nil {
		return eachStatement
	}
	if fnStatement := parser.ParseFnStatement(); fnStatement != nil {
		return fnStatement
	}
	if classStatement := parser.ParseClassStatement(); classStatement != nil {
		return classStatement
	}
	if interfaceStatement := parser.ParseInterfaceStatement(); interfaceStatement != nil {
		return interfaceStatement
	}
	if tryCatchStatement := parser.ParseTryCatchStatement(); tryCatchStatement != nil {
		return tryCatchStatement
	}

	return nil
}

func (parser *Parser) ParseSimpleStatement(needSemiEnd bool) SimpleStatement {
	if expression := parser.ParseExpression(); expression != nil {
		if expression, isPrimary := expression.(PrimaryExpression); isPrimary && parser.MatchCurrentTokenType(TokenTypeComma) {
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
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"expected a equal mark for assignment list!", ParsingUnexpected))
				return nil
			}
			assignListStatement.Token = parser.CurrentToken
			parser.PeekNextToken() // 移过 '='

			if valueList := parser.ParseExpressionList(); valueList != nil {
				assignListStatement.Values = valueList
				if !parser.AssertCurrentTokenIs(TokenTypeSemi, "a semicolon",
					"to terminate a assignment list!") {
					return nil
				}
				return assignListStatement
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"expected a list of expression as values for assignment list!", ParsingUnexpected))
				return nil
			}

		} else if parser.MatchCurrentTokenType(TokenTypeDoublePlus) || parser.MatchCurrentTokenType(TokenTypeDoubleMinus) {
			incDecStatement := new(IncDecStatement)
			incDecStatement.Expression = expression
			incDecStatement.Operator = parser.CurrentToken

			parser.PeekNextToken() // 移过 '++'/'--'

			if needSemiEnd {
				if !parser.AssertCurrentTokenIs(TokenTypeSemi, "a semicolon",
					"to terminate increase/decrease statement") {
					return nil
				}
			}
			return incDecStatement
		}

		if needSemiEnd {
			if !parser.AssertCurrentTokenIs(TokenTypeSemi, "a semicolon",
				"to terminate this statement!") {
				return nil
			}
		}
		return expression // 表达式作为语句
	}
	if varDeclStatement := parser.ParseVarDeclStatement(); varDeclStatement != nil {
		if needSemiEnd {
			parser.PeekNextToken() // 移过 ';'
		}
		return varDeclStatement
	}
	return nil
}

func (parser *Parser) ParseVarDeclElement(mutable bool) *VarDeclElement {
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
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						fmt.Sprintf("expected an expression as initial value for variable '%s'", varNameToken.Str),
						ParsingUnexpected))
					return nil
				}

			} else if parser.MatchCurrentTokenType(TokenTypeComma) || parser.MatchCurrentTokenType(TokenTypeSemi) {
				// 此时即没有给出初始值
				if !mutable {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"no initial value for \"val\" declaration is not allowed!", ParsingUnexpected))
					return nil
				}
				CoralCompileWarningWithPos(parser, fmt.Sprintf(`no initial value for variable: "%s".`, varNameToken.Str))
				// 那么一个变量定义元素可以结束了，不移过逗号 ','、分号';' 而等待外部断言
				return varDeclElement
			}
		}

		// 无类型标注 -> 则当前 token 必须是等号、有初始值
		if parser.MatchCurrentTokenType(TokenTypeEqual) {
			parser.PeekNextToken() // 移过 '='
			if initValue := parser.ParseExpression(); initValue != nil {
				varDeclElement.InitValue = initValue
				return varDeclElement
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					fmt.Sprintf("expected an expression as initial value to type inferring "+
						"for the non-typed variable '%s'!", varNameToken.Str),
					ParsingUnexpected))
				return nil
			}
		} else {
			if varDeclElement.Type == nil && varDeclElement.InitValue == nil {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"a variable initialized with neither a type descriptor nor a initial value is not allowed!",
					ParsingUnexpected))
				return nil
			}

			// 其他不正确的 token
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				fmt.Sprintf("unexpected token '%s' for variabel declaration!", parser.CurrentToken.Str),
				ParsingUnexpected))
			return nil
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
		for varDeclElement := parser.ParseVarDeclElement(varDeclStatement.Mutable); varDeclElement != nil; varDeclElement = parser.ParseVarDeclElement(varDeclStatement.Mutable) {
			varDeclStatement.Declarations = append(varDeclStatement.Declarations, varDeclElement)

			if parser.MatchCurrentTokenType(TokenTypeSemi) {
				// 分号即应该结束此段定义语句，是否取下一个 token 看外部函数是否 needSemiEnd
				return varDeclStatement
			} else {
				if !parser.AssertCurrentTokenIs(TokenTypeComma, "a comma",
					"to separate multiple variable declarations!") {
					return nil
				}
			}
		}
	}

	return nil
}

func (parser *Parser) ParseBreakStatement() *BreakStatement {
	if parser.MatchCurrentTokenType(TokenTypeBreak) {
		breakToken := parser.CurrentToken
		parser.PeekNextToken() // 移过 'break'

		if !parser.AssertCurrentTokenIs(TokenTypeSemi, "a semicolon",
			"to terminate a break statement!") {
			return nil
		}
		return &BreakStatement{Token: breakToken}
	}

	return nil
}

func (parser *Parser) ParseContinueStatement() *ContinueStatement {
	if parser.MatchCurrentTokenType(TokenTypeContinue) {
		continueToken := parser.CurrentToken
		parser.PeekNextToken() // 移过 'continue'

		if !parser.AssertCurrentTokenIs(TokenTypeSemi, "a semicolon",
			"to terminate a continue statement!") {
			return nil
		}
		return &ContinueStatement{Token: continueToken}
	}

	return nil
}

func (parser *Parser) ParseReturnStatement() *ReturnStatement {
	if parser.MatchCurrentTokenType(TokenTypeReturn) {
		returnToken := parser.CurrentToken
		parser.PeekNextToken() // 移过 'return'

		if expressionList := parser.ParseExpressionList(); expressionList != nil {
			if !parser.AssertCurrentTokenIs(TokenTypeSemi, "a semicolon",
				"to terminate a return statement!") {
				return nil
			}
			return &ReturnStatement{
				Token:      returnToken,
				Expression: expressionList,
			}
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected an expression for return statement!", ParsingUnexpected))
			return nil
		}
	}

	return nil
}

func (parser *Parser) ParseImportElement() *ImportElement {
	importElement := new(ImportElement)
	if moduleName := parser.ParseIdentifier(false); moduleName != nil {
		importElement.ModuleName = moduleName
		if parser.MatchCurrentTokenType(TokenTypeAs) {
			parser.PeekNextToken() // 移过 'as'
			if asName := parser.ParseIdentifier(false); asName != nil {
				importElement.As = asName
				return importElement
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"expected an identifier as target's another name for import statement!", ParsingUnexpected))
				return nil
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
		if from, isStringLit := parser.ParseLiteral().(*StringLit); isStringLit && from != nil {
			if parser.MatchCurrentTokenType(TokenTypeImport) {
				parser.PeekNextToken() // 移过 'import'
				// 进入分支判断：
				if parser.MatchCurrentTokenType(TokenTypeLeftBrace) {
					// 当前的如果是大括号说明是 群组引入
					parser.PeekNextToken() // 移过 '{'
					if elementList := parser.ParseImportElementList(); elementList != nil {
						listImportStatement := &ListImportStatement{
							From:     from.Value.Str,
							Elements: elementList,
						}

						if parser.MatchCurrentTokenType(TokenTypeRightBrace) {
							parser.PeekNextToken() // 移过 '}'
							return listImportStatement
						} else {
							CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
								"expected a right brace as ending for a block import statement!", ParsingUnexpected))
							return nil
						}
					} else {
						CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
							"expected at least two import element for a block import statement!", ParsingUnexpected))
						return nil
					}
				} else if importElement := parser.ParseImportElement(); importElement != nil {
					singleImportStatement := new(SingleFromImportStatement)
					singleImportStatement.From = from.Value.Str
					singleImportStatement.Element = importElement
					if parser.MatchCurrentTokenType(TokenTypeSemi) {
						parser.PeekNextToken() // 移过 ';'
						return singleImportStatement
					} else {
						CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
							"expected a semicolon as ending for import statement!", ParsingUnexpected))
						return nil
					}
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected a module name as target for import statement!", ParsingUnexpected))
					return nil
				}
			}
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a module name as source for import statement!", ParsingUnexpected))
			return nil
		}
	} else if parser.MatchCurrentTokenType(TokenTypeImport) {
		parser.PeekNextToken() // 移过 'import'
		if path, isStringLit := parser.ParseLiteral().(*StringLit); isStringLit && path != nil {
			singleGlobalImport := &SingleGlobalImportStatement{Path: path.Value.Str}

			if parser.MatchCurrentTokenType(TokenTypeAs) {
				parser.PeekNextToken() // 移过 'as'

				if as := parser.ParseIdentifier(false); as != nil {
					singleGlobalImport.As = as
				}
			} // 可能有 as xxx 令别名

			if !parser.AssertCurrentTokenIs(TokenTypeSemi, "a semicolon",
				"to terminate a single global import statement!") {
				return nil
			}
			return singleGlobalImport
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a source file path for importee!", ParsingUnexpected))
			return nil
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
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"expected a decimal literal as the enum element's value!", ParsingUnexpected))
				return nil
			}
		}

		if !parser.MatchCurrentTokenType(TokenTypeComma) && !parser.MatchCurrentTokenType(TokenTypeRightBrace) {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a comma to separate multiple enum elements!", ParsingUnexpected))
			return nil
		}
		return enumElement
	}

	return nil
}

func (parser *Parser) ParseEnumStatement() *EnumStatement {
	if parser.MatchCurrentTokenType(TokenTypeEnum) {
		parser.PeekNextToken() // 移过 'enum'
		if enumName := parser.ParseIdentifier(false); enumName != nil {
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
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected a right brace as ending for enum definition!", ParsingUnexpected))
					return nil
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
		}

		if parser.MatchCurrentTokenType(TokenTypeRightBrace) {
			parser.PeekNextToken()
			return blockStatement
		} else {
			// 没有正常解析到右括号
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a right brace as ending for block statement!", ParsingUnexpected))
			return nil
		}
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
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				`expected a block as a "if" block for "if" statement!`, ParsingUnexpected))
			return nil
		}
	} else {
		CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
			`expected an expression as a condition for "if" statement!`, ParsingUnexpected))
		return nil
	}
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
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						`expected a block statement for "else" statement!`, ParsingUnexpected))
					return nil
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
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					`expected an condition and block statement for "elif" statement!`, ParsingUnexpected))
				return nil
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
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected a block statement as case handler!", ParsingUnexpected))
					return nil, false
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
						CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
							"expected a block statement as case handler!", ParsingUnexpected))
						return nil, false
					}
				} else {
					if caseBlock := parser.ParseBlockStatement(); caseBlock != nil {
						normalCase.Block = caseBlock
						return normalCase, false
					} else {
						CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
							"expected a block statement as case handler!", ParsingUnexpected))
						return nil, false
					}
				}
			}
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected an expression as a case!", ParsingUnexpected))
			return nil, false
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
							CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
								`expected a block after 'default' keyword in switch statement!`, ParsingUnexpected))
							return nil
						}
					} else {
						switchStatement.Cases = append(switchStatement.Cases, _case)
					}
				}

				if parser.MatchCurrentTokenType(TokenTypeRightBrace) {
					parser.PeekNextToken() // 移过 '}'
					return switchStatement
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected a right brace as ending for switch statement!", ParsingUnexpected))
					return nil
				}
			}

		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected an expression as target for switch statement!", ParsingUnexpected))
			return nil
		}
	}

	return nil
}

func (parser *Parser) ParseWhileStatement() *WhileStatement {
	if parser.MatchCurrentTokenType(TokenTypeWhile) {
		parser.PeekNextToken() // 移过 'while'

		if condition := parser.ParseExpression(); condition != nil {
			whileStatement := new(WhileStatement)
			whileStatement.Condition = condition

			if whileBlock := parser.ParseBlockStatement(); whileBlock != nil {
				whileStatement.Block = whileBlock
				return whileStatement
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"expected a block statement in while statement!", ParsingUnexpected))
				return nil
			}

		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected an expression as condition for while statement!", ParsingUnexpected))
			return nil
		}
	}

	return nil
}

func (parser *Parser) ParseForStatement() *ForStatement {
	if parser.MatchCurrentTokenType(TokenTypeFor) {
		parser.PeekNextToken() // 移过 'for'
		forStatement := new(ForStatement)

		if parser.MatchCurrentTokenType(TokenTypeSemi) {
			// 可能没有初始化操作的语句
		} else if initial := parser.ParseSimpleStatement(false); initial != nil {
			forStatement.Initial = initial
		}

		// 但总之需要一个分号
		if !parser.AssertCurrentTokenIs(TokenTypeSemi, "the first semicolon", "in for clause!") {
			return nil
		}

		if condition := parser.ParseExpression(); condition != nil {
			forStatement.Condition = condition
			if !parser.AssertCurrentTokenIs(TokenTypeSemi, "the second semicolon", "in for clause!") {
				return nil
			}

			if parser.MatchCurrentTokenType(TokenTypeLeftBrace) {
				CoralCompileWarningWithPos(parser, `a "for" loop only defined with condition, consider using 
while condition { ... } instead.`)
			} else {
				for appendix := parser.ParseSimpleStatement(false); appendix != nil; appendix = parser.ParseSimpleStatement(false) {
					forStatement.Appendix = append(forStatement.Appendix, appendix)
					if parser.MatchCurrentTokenType(TokenTypeComma) {
						parser.PeekNextToken() // 移过 ','
					} else {
						break
					}
				}
			}

			if forBlock := parser.ParseBlockStatement(); forBlock != nil {
				forStatement.Block = forBlock
				return forStatement
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					`expected a block statement in "for" statement!`, ParsingUnexpected))
				return nil
			}
		} else {
			// 不允许没有 for 循环的条件，如果需要一个无限循环，提示建议用 while true
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected an expression as condition in \"for\" statement!\n  "+
					"Tips: If you need a infinite loop, please use 'while true { ... }'", ParsingUnexpected))
			return nil
		}
	}

	return nil
}

func (parser *Parser) ParseEachStatement() *EachStatement {
	if parser.MatchCurrentTokenType(TokenTypeEach) {
		parser.PeekNextToken() // 移过 'each'
		eachStatement := new(EachStatement)

		if elementId := parser.ParseIdentifier(false); elementId != nil {
			eachStatement.Element = elementId

			if parser.MatchCurrentTokenType(TokenTypeComma) {
				parser.PeekNextToken() // 移过 ','
				if keyId := parser.ParseIdentifier(false); keyId != nil {
					eachStatement.Key = keyId
				}
			} // 没有 key Identifier 也不算错

			if parser.MatchCurrentTokenType(TokenTypeIn) {
				parser.PeekNextToken() // 移过 'in'

				if iterateTarget := parser.ParseExpression(); iterateTarget != nil {
					eachStatement.Target = iterateTarget

					if block := parser.ParseBlockStatement(); block != nil {
						eachStatement.Block = block
						return eachStatement
					} else {
						CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
							"expected a block statement for \"each\" iteration loop!", ParsingUnexpected))
						return nil
					}
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected an expression as a target for \"each\" iteration loop!", ParsingUnexpected))
					return nil
				}
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"expected a \"in\" keyword for \"each\" iteration loop!", ParsingUnexpected))
				return nil
			}
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected at least one identifier for \"each\" iteration loop!", ParsingUnexpected))
			return nil
		}
	}

	return nil
}

func (parser *Parser) ParseArgument() *Argument {
	if argName := parser.ParseIdentifier(false); argName != nil {
		argument := new(Argument)
		argument.Name = argName

		if argType := parser.ParseTypeDescription(); argType != nil {
			argument.Type = argType
		}

		return argument
	}

	return nil
}

func (parser *Parser) ParseArgumentList(allowIgnoreTyping bool) []*Argument {
	var argList []*Argument
	var noTypeDescriptorList []*Argument
	currentInShorthand := false

	for arg := parser.ParseArgument(); arg != nil; arg = parser.ParseArgument() {
		if arg.Type == nil {
			// 监测到一个没有类型声明的形参
			noTypeDescriptorList = append(noTypeDescriptorList, arg) // 记录入队
			currentInShorthand = true                                // 亮起标志位，之后如果遇到有类型，则清空队列
		} else {
			if currentInShorthand {
				for _, noTypingArg := range noTypeDescriptorList {
					noTypingArg.Type = arg.Type
					argList = append(argList, noTypingArg)
				}
				noTypeDescriptorList = make([]*Argument, 0) // 让 GC 回收原队列切片内存
				currentInShorthand = false // 重置标志
			}
		}
		argList = append(argList, arg)

		if parser.MatchCurrentTokenType(TokenTypeComma) {
			parser.PeekNextToken() // 移过 ','
		} else {
			break
		}
	}

	// 此时的 parse 进行时不允许省略类型标注，最后又仍是所有变量没类型：
	if !allowIgnoreTyping && currentInShorthand {
		CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
			"expected at least one type for all arguments!!", ParsingUnexpected))
		return nil
	}
	return argList
}

func (parser *Parser) ParseReturnList() []TypeDescription {
	var returnList []TypeDescription
	for returnType := parser.ParseTypeDescription(); returnType != nil; returnType = parser.ParseTypeDescription() {
		returnList = append(returnList, returnType)
		if parser.MatchCurrentTokenType(TokenTypeComma) {
			parser.PeekNextToken() // 移过 ','
		} else {
			break
		}
	}
	return returnList
}

func (parser *Parser) ParseSignature(allowMismatched bool, allowIgnoreTyping bool) *Signature {
	if parser.MatchCurrentTokenType(TokenTypeLeftParen) {
		parser.PeekNextToken() // 移过左括号
		signature := new(Signature)
		signature.Arguments = parser.ParseArgumentList(allowIgnoreTyping)
		if !parser.MatchCurrentTokenType(TokenTypeRightParen) {
			if allowMismatched {
				return nil
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"expected a right parenthesis in the function signature!", ParsingUnexpected))
				return nil
			}
		}
		parser.PeekNextToken() // 移过右括号
		signature.Returns = parser.ParseReturnList()

		if parser.MatchCurrentTokenType(TokenTypeThrows) {
			parser.PeekNextTokenAvoidAngleConfusing() // 移过 'throws'

			for exceptionType := parser.ParseTypeDescription(); exceptionType != nil; exceptionType = parser.ParseTypeDescription() {
				signature.Throws = append(signature.Throws, exceptionType)
				if parser.MatchCurrentTokenType(TokenTypeComma) {
					parser.PeekNextTokenAvoidAngleConfusing() // 移过 ','
				} else {
					break
				}
			}
			if len(signature.Throws) == 0 {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"expected the exceptions' type after keyword \"throws\"!", ParsingUnexpected))
				return nil
			}
		}

		return signature
	}

	return nil
}

func (parser *Parser) ParseGenericsArgElement() *GenericsArgElement {
	if argName := parser.ParseIdentifier(true); argName != nil {
		argElement := new(GenericsArgElement)
		argElement.ArgName = argName

		if argGenerics := parser.ParseGenericsArgs(); argGenerics != nil {
			argElement.Generics = argGenerics
		} // 也可能只是通配符 而不是其他泛型类

		return argElement
	}

	return nil
}

func (parser *Parser) ParseGenericsArgs() *GenericArgs {
	if parser.MatchCurrentTokenType(TokenTypeLeftAngle) {
		parser.PeekNextTokenAvoidAngleConfusing() // 移过 '<'
		genericsArg := new(GenericArgs)

		for element := parser.ParseGenericsArgElement(); element != nil; element = parser.ParseGenericsArgElement() {
			genericsArg.Args = append(genericsArg.Args, element)

			if parser.MatchCurrentTokenType(TokenTypeComma) {
				parser.PeekNextTokenAvoidAngleConfusing() // 移过 ','
				if parser.MatchCurrentTokenType(TokenTypeRightAngle) {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"generics arguments can't be end with a comma!", ParsingUnexpected))
					return nil
				} // 不允许出现 <T,K,>
			} else {
				break
			}
		}

		if !parser.AssertCurrentTokenIs(TokenTypeRightAngle, "a right angle",
			"to terminate a generics arguments!") {
			return nil
		}
		return genericsArg
	}

	return nil
}

func (parser *Parser) ParseFnStatement() *FunctionDeclarationStatement {
	if parser.MatchCurrentTokenType(TokenTypeFn) {
		parser.PeekNextToken() // 移过 'fn'
		fnStmt := new(FunctionDeclarationStatement)

		if fnName := parser.ParseIdentifier(true); fnName != nil {
			// 取 Identifier 结束后，GetNextToken 时避免读取 << 导致词法解析错误
			// avoidAngleConfusing 这个项不会影响到其他类型 Token 的解析，只是于尖括号的解析相关
			fnStmt.Name = fnName

			if fnGenerics := parser.ParseGenericsArgs(); fnGenerics != nil {
				fnStmt.Generics = fnGenerics
			} // 函数也可能没有泛型参数

			if signature := parser.ParseSignature(false, false); signature != nil {
				fnStmt.Signature = signature

				if fnBlock := parser.ParseBlockStatement(); fnBlock != nil {
					fnStmt.Block = fnBlock
					return fnStmt
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected a block when defining a function statement!", ParsingUnexpected))
					return nil
				}
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					"expected a signature when defining a function statement!", ParsingUnexpected))
				return nil
			}
		}
	}

	return nil
}

func (parser *Parser) ParseClassIdentifier() *ClassIdentifier {
	if className := parser.ParseIdentifier(true); className != nil {
		classIdentifier := new(ClassIdentifier)
		classIdentifier.Name = className

		if genericsArgs := parser.ParseGenericsArgs(); genericsArgs != nil {
			classIdentifier.Generics = genericsArgs
		} // 也可能没有泛型参数

		return classIdentifier
	}

	return nil
}

func (parser *Parser) ParseClassMember() ClassMember {
	var scopeType ClassMemberScopeType = ClassMemberScopePrivate
	if parser.MatchCurrentTokenType(TokenTypePublic) {
		scopeType = ClassMemberScopePublic
		parser.PeekNextToken()
	} else if parser.MatchCurrentTokenType(TokenTypePrivate) {
		parser.PeekNextToken()
	}
	if memberVarDecl := parser.ParseVarDeclStatement(); memberVarDecl != nil {
		if !parser.AssertCurrentTokenIs(TokenTypeSemi, "a semicolon",
			"to terminate a class member variable declaration!") {
			return nil
		}

		classMemberVar := new(ClassMemberVar)
		classMemberVar.Scope = scopeType
		classMemberVar.VarDecl = memberVarDecl
		return classMemberVar
	} else if memberMethodDecl := parser.ParseFnStatement(); memberMethodDecl != nil {
		classMemberMethod := new(ClassMemberMethod)
		classMemberMethod.Scope = scopeType
		classMemberMethod.MethodDecl = memberMethodDecl
		return classMemberMethod
	}

	return nil
}

func (parser *Parser) ParseClassStatement() *ClassDeclarationStatement {
	if parser.MatchCurrentTokenType(TokenTypeClass) {
		parser.PeekNextTokenAvoidAngleConfusing() // 移过 'class'
		classStmt := new(ClassDeclarationStatement)

		if classId := parser.ParseClassIdentifier(); classId != nil {
			classStmt.Definition = classId

			if parser.MatchCurrentTokenType(TokenTypeColon) {
				parser.PeekNextTokenAvoidAngleConfusing() // 移过 ':'
				if extends := parser.ParseClassIdentifier(); extends != nil {
					classStmt.Extends = extends
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected an class identifier for extended class name!", ParsingUnexpected))
					return nil
				}
			} // 也可能没有继承

			if parser.MatchCurrentTokenType(TokenTypeLeftArrow) {
				parser.PeekNextToken() // 移过 左箭头
				for impl := parser.ParseClassIdentifier(); impl != nil; impl = parser.ParseClassIdentifier() {
					classStmt.Implements = append(classStmt.Implements, impl)

					if parser.MatchCurrentTokenType(TokenTypeComma) {
						parser.PeekNextTokenAvoidAngleConfusing()
					} else {
						break
					}
				}
			}

			if !parser.AssertCurrentTokenIs(TokenTypeLeftBrace, "a left brace",
				"to start the class statement definition body!") {
				return nil
			}

			hasInitMethod := false
			for member := parser.ParseClassMember(); member != nil; member = parser.ParseClassMember() {
				classStmt.Members = append(classStmt.Members, member)

				if method, isMethod := member.(*ClassMemberMethod); !hasInitMethod && isMethod && method.MethodDecl.Name.Token.Str == classId.Name.Token.Str {
					hasInitMethod = true
					method.Scope = ClassMemberScopePublic // 构造方法默认 public
					if len(method.MethodDecl.Signature.Returns) > 0 {
						CoralCompileErrorWithPos(parser, NewCoralError("Compile",
							fmt.Sprintf("Any returns by constructor method of class \"%s\" are not allowed!",
								classId.Name.Token.Str),
							NoConstructorMethod))
						return nil
					}
				}
				if parser.MatchCurrentTokenType(TokenTypeRightBrace) {
					break // 等待外部断言
				}
			}

			if !hasInitMethod {
				CoralCompileErrorWithPos(parser, NewCoralError("Compile",
					fmt.Sprintf("expected a constructor for class \"%s\"!", classId.Name.Token.Str),
					NoConstructorMethod))
				return nil
			} // <- 没有构造函数的报错

			if !parser.AssertCurrentTokenIs(TokenTypeRightBrace, "a right brace",
				"to terminate the class statement definition body!") {
				return nil
			}

			return classStmt

		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected an class identifier for class name!", ParsingUnexpected))
			return nil
		}
	}

	return nil
}

func (parser *Parser) ParseInterfaceMethodDecl() *InterfaceMethodDeclaration {
	var scopeType ClassMemberScopeType = ClassMemberScopePrivate
	if parser.MatchCurrentTokenType(TokenTypePublic) {
		scopeType = ClassMemberScopePublic
		parser.PeekNextToken()
	} else if parser.MatchCurrentTokenType(TokenTypePrivate) {
		parser.PeekNextToken()
	}

	if !parser.AssertCurrentTokenIs(TokenTypeFn, "keyword \"fn\"",
		"to start the announcement of the method in interface declaration statement!") {
		return nil
	}

	methodDecl := new(InterfaceMethodDeclaration)
	methodDecl.Scope = scopeType
	if interfaceName := parser.ParseIdentifier(true); interfaceName != nil {
		methodDecl.Name = interfaceName

		if methodGenerics := parser.ParseGenericsArgs(); methodGenerics != nil {
			methodDecl.Generics = methodGenerics
		} // 也可能没有泛型参数

		if signature := parser.ParseSignature(false, false); signature != nil {
			methodDecl.Signature = signature

			if !parser.AssertCurrentTokenIs(TokenTypeSemi, "a semicolon",
				"to terminate a interface method declaration!") {
				return nil
			}

			return methodDecl
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				fmt.Sprintf("expected a function signature for method \"%s\"!", interfaceName.Token.Str),
				ParsingUnexpected))
			return nil
		}
	}

	return nil
}

func (parser *Parser) ParseInterfaceStatement() *InterfaceDeclarationStatement {
	if parser.MatchCurrentTokenType(TokenTypeInterface) {
		parser.PeekNextTokenAvoidAngleConfusing() // 移过 'interface'
		interfaceStmt := new(InterfaceDeclarationStatement)

		if interfaceId := parser.ParseClassIdentifier(); interfaceId != nil {
			interfaceStmt.Definition = interfaceId

			if parser.MatchCurrentTokenType(TokenTypeColon) {
				parser.PeekNextTokenAvoidAngleConfusing() // 移过 ':'
				if extends := parser.ParseClassIdentifier(); extends != nil {
					interfaceStmt.Extends = extends
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected an class identifier for extended class name!", ParsingUnexpected))
					return nil
				}
			} // 也可能没有继承

			if !parser.AssertCurrentTokenIs(TokenTypeLeftBrace, "a left brace",
				"to start the interface statement definition body!") {
				return nil
			}

			for method := parser.ParseInterfaceMethodDecl(); method != nil; method = parser.ParseInterfaceMethodDecl() {
				interfaceStmt.Methods = append(interfaceStmt.Methods, method)

				if method.Name.Token.Str == interfaceId.Name.Token.Str {
					CoralCompileErrorWithPos(parser, NewCoralError("Compile",
						fmt.Sprintf("method name being the same with interface name \"%s\"", interfaceId.Name.Token.Str),
						MethodNameSameWithInterfaceName))
					return nil
				} // 方法名不能与接口名相同！

				if parser.MatchCurrentTokenType(TokenTypeRightBrace) {
					break // 等待外部断言
				}
			}

			if len(interfaceStmt.Methods) == 0 {
				CoralCompileErrorWithPos(parser, NewCoralError("Compile",
					fmt.Sprintf("expected at least one method for interface \"%s\"!", interfaceId.Name.Token.Str),
					EmptyInterfaceDeclaration))
				return nil
			} // <- 空接口

			if !parser.AssertCurrentTokenIs(TokenTypeRightBrace, "a right brace",
				"to terminate the interface statement definition body!") {
				return nil
			}

			return interfaceStmt

		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected an class identifier for interface name!", ParsingUnexpected))
			return nil
		}
	}

	return nil
}

func (parser *Parser) ParseErrorCatchHandler() *ErrorCatchHandler {
	if parser.MatchCurrentTokenType(TokenTypeCatch) {
		parser.PeekNextToken() // 移过 'catch'
		errHandler := new(ErrorCatchHandler)

		if errId := parser.ParseIdentifier(false); errId != nil {
			errHandler.Name = errId

			if errType := parser.ParseTypeDescription(); errType != nil {
				errHandler.ErrorType = errType

				if handleBlock := parser.ParseBlockStatement(); handleBlock != nil {
					errHandler.Handler = handleBlock

					return errHandler
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						fmt.Sprintf("expected a block as handler for exception \"%s\"!", errId.Token.Str),
						ParsingUnexpected))
					return nil
				}
			} else {
				CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
					fmt.Sprintf("expected a type descriptor for exception \"%s\"!", errId.Token.Str),
					ParsingUnexpected))
				return nil
			}
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected an identifier for exception name after keyword \"catch\"!", ParsingUnexpected))
			return nil
		}
	}

	return nil
}

func (parser *Parser) ParseTryCatchStatement() *TryCatchStatement {
	if parser.MatchCurrentTokenType(TokenTypeTry) {
		parser.PeekNextToken() // 移过 'try'
		tryCatchStmt := new(TryCatchStatement)

		if tryBlock := parser.ParseBlockStatement(); tryBlock != nil {
			tryCatchStmt.TryBlock = tryBlock

			for errHandler := parser.ParseErrorCatchHandler(); errHandler != nil; errHandler = parser.ParseErrorCatchHandler() {
				tryCatchStmt.Handlers = append(tryCatchStmt.Handlers, errHandler)
			}

			if parser.MatchCurrentTokenType(TokenTypeFinally) {
				parser.PeekNextToken() // 移过 'finally'
				if finallyBlock := parser.ParseBlockStatement(); finallyBlock != nil {
					tryCatchStmt.Finally = finallyBlock

					return tryCatchStmt
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						fmt.Sprintf("expected a block after keyword \"finally\"!"),
						ParsingUnexpected))
					return nil
				}
			} // 也可能无 finally

			return tryCatchStmt
		} else {
			CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
				"expected a block statement after keyword \"try\"!", ParsingUnexpected))
			return nil
		}
	}

	return nil
}
