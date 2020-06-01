package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
	"fmt"
)

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
						parser.AssertCurrentTokenIs(TokenTypeComma, "a comma", fmt.Sprintf(
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
