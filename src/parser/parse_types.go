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
						CoralErrorCrashHandlerWithPos(parser, NewCoralError("Syntax",
							"expected a decimal number as array length declaration!", ParsingUnexpected))
					}
					arrayLit.ArrayLength = arrLen
				}

				parser.AssertCurrentTokenIs(TokenTypeRightBracket, "a right bracket",
					"to terminate a array type descriptor!")
				return arrayLit
			} else {
				// 否则就将 typeName 返回作为该 typeDescription
				return typeName
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
