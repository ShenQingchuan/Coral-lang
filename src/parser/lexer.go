package parser

import (
	"coral-lang/src/exception"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

const (
	TokenTypeUnknown = iota
	TokenTypeImport
	TokenTypeFrom
	TokenTypeAs
	TokenTypeEnum
	TokenTypeBreak
	TokenTypeContinue
	TokenTypeReturn
	TokenTypeVar
	TokenTypeVal
	TokenTypeIf
	TokenTypeElif
	TokenTypeElse
	TokenTypeSwitch
	TokenTypeDefault
	TokenTypeCase
	TokenTypeWhile
	TokenTypeFor
	TokenTypeEach
	TokenTypeFn
	TokenTypeClass
	TokenTypeInterface
	TokenTypeNew
	TokenTypeNil
	TokenTypeTrue
	TokenTypeFalse
	TokenTypeTry
	TokenTypeCatch
	TokenTypeFinally

	TokenTypeSemi                  // ;
	TokenTypeColon                 // ,
	TokenTypeLeftParen             // (
	TokenTypeRightParen            // )
	TokenTypeLeftBrace             // {
	TokenTypeRightBrace            // }
	TokenTypeLeftBracket           // [
	TokenTypeRightBracket          // ]
	TokenTypeDot                   // .
	TokenTypeEqual                 // =
	TokenTypeDoubleEqual           // ==
	TokenTypeBangEqual             // !=
	TokenTypePlus                  // +
	TokenTypeMinus                 // -
	TokenTypeStar                  // *
	TokenTypeDoubleStar            // **
	TokenTypeSlash                 // /
	TokenTypePercent               // %
	TokenTypeCaret                 // ^
	TokenTypeAmpersand             // &
	TokenTypeBang                  // !
	TokenTypeVertical              // |
	TokenTypeLeftAngle             // <
	TokenTypeRightAngle            // >
	TokenTypeDoubleLeftAngle       // <<
	TokenTypeDoubleRightAngle      // >>
	TokenTypeDoubleAmpersand       // &&
	TokenTypeDoubleVertical        // ||
	TokenTypeLeftAngleEqual        // <=
	TokenTypeRightAngleEqual       // >=
	TokenTypeLeftArrow             // <-
	TokenTypeRightArrow            // ->
	TokenTypeDoublePlus            // ++
	TokenTypeDoubleMinus           // --
	TokenTypePlusEqual             // +=
	TokenTypeMinusEqual            // -=
	TokenTypeStarEqual             // *=
	TokenTypeSlashEqual            // /=
	TokenTypePercentEqual          // %=
	TokenTypeDoubleLeftAngleEqual  // <<=
	TokenTypeDoubleRightAngleEqual // >>=
	TokenTypeAmpersandEqual        // &=
	TokenTypeVerticalEqual         // |=
	TokenTypeCaretEqual            // ^=

	TokenTypeDecimalInteger
	TokenTypeOctalInteger
	TokenTypeHexadecimalInteger
	TokenTypeBinaryInteger
	TokenTypeExponent
	TokenTypeFloat
	TokenTypeRune
	TokenTypeString

	TokenTypeIdentifier
)

type TokenType = int
type Token struct {
	Line, Col int
	Kind      TokenType
	Str       string
}

type UTF8Char struct {
	Rune       rune // utf8.decode 解码出的 utf8 单字符
	ByteLength int  // 对应实际字节数
}
type Lexer struct {
	Content     []byte          // 源代码的buffer
	KeywordMap  *map[string]int // 关键字映射表
	OperatorMap *map[string]int // 关键字映射表

	Line, Col int // 记录行号列号
	BytePos   int // 当前游标位置
}

// 给出路径，打开源代码文件
func OpenSourceFile(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		exception.CoralError("FileSystem", "Can't open source file: "+filePath)
	}
	if file != nil {
		defer file.Close()
	}

	content, err := ioutil.ReadAll(file)
	return content
}

// 初始化词法分析器
func InitLexer(lexer *Lexer, content []byte) {
	lexer.Content = content
	lexer.Line = 1
	lexer.Col = 1
	lexer.BytePos = 0

	lexer.KeywordMap = &map[string]TokenType{
		"import":    TokenTypeImport,
		"from":      TokenTypeFrom,
		"as":        TokenTypeAs,
		"enum":      TokenTypeEnum,
		"break":     TokenTypeBreak,
		"continue":  TokenTypeContinue,
		"return":    TokenTypeReturn,
		"var":       TokenTypeVar,
		"val":       TokenTypeVal,
		"if":        TokenTypeIf,
		"elif":      TokenTypeElif,
		"else":      TokenTypeElse,
		"switch":    TokenTypeSwitch,
		"default":   TokenTypeDefault,
		"case":      TokenTypeCase,
		"while":     TokenTypeWhile,
		"for":       TokenTypeFor,
		"each":      TokenTypeEach,
		"fn":        TokenTypeFn,
		"class":     TokenTypeClass,
		"interface": TokenTypeInterface,
		"new":       TokenTypeNew,
		"nil":       TokenTypeNil,
		"true":      TokenTypeTrue,
		"false":     TokenTypeFalse,
		"try":       TokenTypeTry,
		"catch":     TokenTypeCatch,
		"finally":   TokenTypeFinally,
	}
}

// 拾取当前游标所在位置的字符
func (lexer *Lexer) PeekChar() *UTF8Char {
	r, byteLength := utf8.DecodeRune(lexer.Content[lexer.BytePos:])
	return &UTF8Char{
		Rune:       r,
		ByteLength: byteLength,
	}
}

// 拾取游标处的下一个字符
func (lexer *Lexer) PeekNextChar(currentLength int) *UTF8Char {
	r, byteLength := utf8.DecodeRune(lexer.Content[lexer.BytePos+currentLength:])
	return &UTF8Char{
		Rune:       r,
		ByteLength: byteLength,
	}
}

// 拾取游标处 + 步数位置的字符
func (lexer *Lexer) PeekNextCharByStep(currentLength int, step int) *UTF8Char {
	forwardSummaryLength := currentLength
	for i := 1; i < step; i++ {
		_, forwardLength := utf8.DecodeRune(lexer.Content[lexer.BytePos+forwardSummaryLength:])
		forwardSummaryLength += forwardLength
	}
	r, byteLength := utf8.DecodeRune(lexer.Content[lexer.BytePos+forwardSummaryLength:])
	return &UTF8Char{
		Rune:       r,
		ByteLength: byteLength,
	}
}

// 游标向前移动一个单位
func (lexer *Lexer) GoNextChar() {
	lexer.BytePos += lexer.PeekChar().ByteLength
}

// 游标向前移动多个单位
func (lexer *Lexer) GoNextCharByStep(step int) {
	forwardSummaryLength := lexer.PeekChar().ByteLength
	for i := 1; i < step; i++ {
		_, forwardLength := utf8.DecodeRune(lexer.Content[lexer.BytePos+forwardSummaryLength:])
		forwardSummaryLength += forwardLength
	}
	lexer.BytePos += forwardSummaryLength
}

// 匹配字符是否为给予的
func (uchar *UTF8Char) MatchRune(r rune) bool {
	return uchar.Rune == r
}

// 字面值是否为合法的十进制数字
func (uchar *UTF8Char) IsLegalDecimal() bool {
	return uchar.Rune >= '0' && uchar.Rune <= '9'
}

// 字面值是否为合法的十六进制数字
func (uchar *UTF8Char) IsLegalHexadecimal() bool {
	return (uchar.Rune >= '0' && uchar.Rune <= '9') || (uchar.Rune >= 'A' && uchar.Rune <= 'F') || (uchar.Rune >= 'a' && uchar.Rune <= 'f')
}

// 字面值是否为合法的八进制数字
func (uchar *UTF8Char) IsLegalOctal() bool {
	return uchar.Rune >= '0' && uchar.Rune <= '7'
}

// 字面值是否为合法的二进制数字
func (uchar *UTF8Char) IsLegalBinary() bool {
	return uchar.Rune == '0' || uchar.Rune == '1'
}

// 读出一个十六进制数的 token
func (lexer *Lexer) ReadHexadecimal() *Token {
	lexer.GoNextCharByStep(2) // skip '0x'
	str := "0x"

	for lexer.PeekChar().IsLegalHexadecimal() {
		str += string(lexer.PeekChar().Rune)
		lexer.GoNextChar()
	}
	return lexer.makeToken(TokenTypeHexadecimalInteger, str)
}

// 读出一个八进制数的 token
func (lexer *Lexer) ReadOctal() *Token {
	lexer.GoNextCharByStep(2) // 跳过 '0o'
	str := "0o"
	for lexer.PeekChar().IsLegalOctal() {
		str += string(lexer.PeekChar().Rune)
		lexer.GoNextChar()
	}
	return lexer.makeToken(TokenTypeOctalInteger, str)
}

// 读出一个二进制数的 token
func (lexer *Lexer) ReadBinary() *Token {
	lexer.GoNextCharByStep(2) // 跳过 '0o'
	str := "0b"
	for lexer.PeekChar().IsLegalBinary() {
		str += string(lexer.PeekChar().Rune)
		lexer.GoNextChar()
	}
	return lexer.makeToken(TokenTypeBinaryInteger, str)
}

// 读出一个十进制数的 token
func (lexer *Lexer) ReadDecimal(startFromZero bool) *Token {
	var str string
	if startFromZero {
		lexer.GoNextChar() // 读入 '0'
		str = "0"
	}

	for lexer.PeekChar().IsLegalDecimal() {
		str += string(lexer.PeekChar().Rune)
		lexer.GoNextChar()
	}
	return lexer.makeToken(TokenTypeDecimalInteger, str)
}

// 产出 token，词法分析器的行号也移动字面值 s 的长度
func (lexer *Lexer) makeToken(t TokenType, s string) *Token {
	lexer.Col += len(s)
	return &Token{
		Line: lexer.Line,
		Col:  lexer.Col,
		Kind: t,
		Str:  s,
	}
}

// 词法分析器获取下一个 token
func (lexer *Lexer) GetNextToken() *Token {
	for lexer.BytePos < len(lexer.Content) {
		c := lexer.PeekChar()
		switch c.Rune {
		default:
			if c.IsLegalDecimal() {
				return lexer.ReadDecimal(false)
			}
		case '\t', ' ':
			lexer.GoNextChar() // skip
		case '\n':
			lexer.Line++
			lexer.Col = 1
			lexer.GoNextChar() // skip
		case ';':
			return lexer.makeToken(TokenTypeSemi, ";")
		case ',':
			return lexer.makeToken(TokenTypeColon, ",")
		case '(':
			return lexer.makeToken(TokenTypeLeftParen, "(")
		case ')':
			return lexer.makeToken(TokenTypeRightParen, ")")
		case '{':
			return lexer.makeToken(TokenTypeLeftBrace, "{")
		case '}':
			return lexer.makeToken(TokenTypeRightBrace, "}")
		case '[':
			return lexer.makeToken(TokenTypeLeftBracket, "[")
		case ']':
			return lexer.makeToken(TokenTypeRightBracket, "]")
		case '.':
			return lexer.makeToken(TokenTypeDot, ".")
		case '=':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				return lexer.makeToken(TokenTypeDoubleEqual, "==")
			}
			return lexer.makeToken(TokenTypeEqual, "=")
		case '!':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				return lexer.makeToken(TokenTypeBangEqual, "!=")
			}
			return lexer.makeToken(TokenTypeBang, "!")
		case '*':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('*') {
				lexer.makeToken(TokenTypeDoubleStar, "**")
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.makeToken(TokenTypeStarEqual, "*=")
			}
			return lexer.makeToken(TokenTypeStar, "*")
		case '/':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				return lexer.makeToken(TokenTypeSlashEqual, "/=")
			}
			return lexer.makeToken(TokenTypeSlash, "/")
		case '%':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				return lexer.makeToken(TokenTypePercentEqual, "%=")
			}
			return lexer.makeToken(TokenTypePercent, "%")
		case '^':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				return lexer.makeToken(TokenTypeCaretEqual, "^=")
			}
			return lexer.makeToken(TokenTypeCaret, "^")
		case '&':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('&') {
				lexer.makeToken(TokenTypeDoubleAmpersand, "&&")
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.makeToken(TokenTypeAmpersandEqual, "&=")
			}
			return lexer.makeToken(TokenTypeAmpersand, "&")
		case '|':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('|') {
				lexer.makeToken(TokenTypeDoubleVertical, "||")
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.makeToken(TokenTypeVerticalEqual, "|=")
			}
			return lexer.makeToken(TokenTypeVertical, "|")
		case '<':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('<') {
				if lexer.PeekNextCharByStep(c.ByteLength, 2).MatchRune('=') {
					return lexer.makeToken(TokenTypeDoubleLeftAngleEqual, "<<=")
				}
				return lexer.makeToken(TokenTypeDoubleLeftAngle, "<<")
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				return lexer.makeToken(TokenTypeLeftAngleEqual, "<=")
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('-') {
				return lexer.makeToken(TokenTypeLeftArrow, "<-")
			}
			return lexer.makeToken(TokenTypeLeftAngle, "<")
		case '>':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('>') {
				if lexer.PeekNextCharByStep(c.ByteLength, 2).MatchRune('=') {
					return lexer.makeToken(TokenTypeDoubleRightAngleEqual, ">>=")
				}
				return lexer.makeToken(TokenTypeDoubleRightAngle, ">>")
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				return lexer.makeToken(TokenTypeRightAngleEqual, ">=")
			}
			return lexer.makeToken(TokenTypeRightAngle, ">")
		case '+':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('+') {
				return lexer.makeToken(TokenTypeDoublePlus, "++")
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				return lexer.makeToken(TokenTypePlusEqual, "+=")
			}
			return lexer.makeToken(TokenTypePlus, "+")
		case '-':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('-') {
				return lexer.makeToken(TokenTypeDoubleMinus, "--")
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				return lexer.makeToken(TokenTypeMinusEqual, "-=")
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('>') {
				return lexer.makeToken(TokenTypeRightArrow, "->")
			}
			return lexer.makeToken(TokenTypeMinus, "-")
		case '0':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('x') {
				return lexer.ReadHexadecimal()
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('o') {
				// return lexer.ReadOctal()
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('b') {
				// return lexer.ReadBinary()
			}
			return lexer.ReadDecimal(true)
			// TODO: 字符串等其他情况...
		}
	}

	return nil // TODO: 未知情况
}
