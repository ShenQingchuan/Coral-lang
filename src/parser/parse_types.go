package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
	"fmt"
	"strconv"
)

func (parser *Parser) ParseTypeDescription() TypeDescription {
	// 如果是 identifier 说明可能是 GenericsLit
	if parser.MatchCurrentTokenType(TokenTypeIdentifier) {
		typeName := parser.ParseTypeName()
		if typeName != nil {
			// 结束 typeName 解析后如果是一个 '<' 则进入 Generics 解析
			if parser.MatchCurrentTokenType(TokenTypeLeftAngle) {
				genericsTypeLit := new(GenericsTypeLit)
				genericsTypeLit.BasicType = typeName
				parser.PeekNextTokenAvoidAngleConfusing() // 越过 '<'

				for {
					genericsLitElement := parser.ParseTypeDescription()
					if genericsLitElement != nil {
						genericsTypeLit.GenericsArgs = append(genericsTypeLit.GenericsArgs, genericsLitElement)
					}
					if parser.MatchCurrentTokenType(TokenTypeRightAngle) {
						parser.PeekNextTokenAvoidAngleConfusing() // 移过 '>'
						return genericsTypeLit                    // 结束泛型参数解析
					} else {
						parser.AssertCurrentTokenIs(TokenTypeComma, "a comma", fmt.Sprintf(
							"to seperate several generics arguments but got '%s'",
							parser.CurrentToken.Str))
					}
				}
			} else if parser.MatchCurrentTokenType(TokenTypeLeftBracket) {
				parser.PeekNextToken() // 移过左中括号
				arrayLit := new(ArrayTypeLit)
				arrayLit.ElementType = typeName

				if arrLenLiteral, isDecimal := parser.ParseLiteral().(*DecimalLit); isDecimal && arrLenLiteral != nil {
					arrLen, convertErr := strconv.Atoi(arrLenLiteral.Value.Str)
					if convertErr != nil {
						CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
							"expected a decimal number as array length declaration!", ParsingUnexpected))
						return nil
					}
					arrayLit.ArrayLength = arrLen
				}

				if !parser.AssertCurrentTokenIs(TokenTypeRightBracket, "a right bracket",
					"to terminate a array type descriptor!") {
					return nil
				}
				return arrayLit
			} else {
				// 否则就将 typeName 返回作为该 typeDescription
				return typeName
			}
		}
	} else if parser.MatchCurrentTokenType(TokenTypeLeftParen) {
		parser.PeekNextToken() // 移过左圆括号
		funcType := new(FuncType)
		for {
			if argType := parser.ParseTypeDescription(); argType != nil {
				funcType.ArgTypes = append(funcType.ArgTypes, argType)
				if parser.MatchCurrentTokenType(TokenTypeComma) {
					parser.PeekNextToken() // 移过逗号
				} else if parser.MatchCurrentTokenType(TokenTypeRightParen) {
					parser.PeekNextToken() // 移过右括号
					break                  // 结束函数类型参数部分解析
				} else {
					CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
						"expected a comma to separate arguments' type or right parenthesis to terminate in function type!",
						ParsingUnexpected))
					return nil
				}
			}
		}
		parser.AssertCurrentTokenIs(TokenTypeRightArrow, "a right arrow",
			"in the function type declaration!")
		for {
			if returnType := parser.ParseTypeDescription(); returnType != nil {
				funcType.ReturnTypes = append(funcType.ReturnTypes, returnType)
				if parser.MatchCurrentTokenType(TokenTypeComma) {
					parser.PeekNextToken() // 移过逗号
				} else {
					return funcType
				}
			}
		}
	}
	return nil
}

func (parser *Parser) ParseTypeName() *TypeName {
	if typeNameId := parser.ParseIdentifier(false); typeNameId != nil {
		typeName := &TypeName{Identifier: typeNameId}
		return typeName
	}

	return nil
}
