package lexer

import (
	. "coral-lang/src/exception"
	"coral-lang/src/utils"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"unicode/utf8"
)

const (
	TokenTypeUnknown = iota
	TokenTypeImport
	TokenTypePackage
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
	TokenTypeIn
	TokenTypeFn
	TokenTypeClass
	TokenTypeInterface
	TokenTypeThis
	TokenTypeSuper
	TokenTypeStatic
	TokenTypePublic
	TokenTypePrivate
	TokenTypeNew
	TokenTypeNil
	TokenTypeTrue
	TokenTypeFalse
	TokenTypeTry
	TokenTypeCatch
	TokenTypeFinally
	TokenTypeThrows

	TokenTypeSemi                  // ;
	TokenTypeComma                 // ,
	TokenTypeColon                 // :
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
	TokenTypeAlpha                 // @
	TokenTypeWavy                  // ~
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
	TokenTypeEllipsis              // ...
	TokenTypeDoubleDot             // ..

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
	Content     []byte         // 源代码的buffer
	KeywordMap  map[string]int // 关键字映射表
	OperatorMap map[string]int // 关键字映射表

	ParenCount   int
	BracketCount int
	BraceCount   int

	Line, Col int // 记录行号列号
	BytePos   int // 当前游标位置
}

// 给出路径，打开源代码文件
func OpenSourceFile(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		CoralErrorCrashHandler(NewCoralError("FileSystem", "Can'Token open source file: "+filePath, FileSystemOpenFileError))
	}
	if file != nil {
		defer file.Close()
	}

	content, err := ioutil.ReadAll(file)
	return content
}

// 初始化词法分析器
func InitLexerCommonOperations(lexer *Lexer) {
	lexer.Line = 1
	lexer.Col = 1
	lexer.BytePos = 0

	lexer.ParenCount = 0
	lexer.BraceCount = 0
	lexer.BracketCount = 0

	lexer.KeywordMap = map[string]TokenType{
		"import":    TokenTypeImport,
		"package":   TokenTypePackage,
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
		"in":        TokenTypeIn,
		"fn":        TokenTypeFn,
		"class":     TokenTypeClass,
		"interface": TokenTypeInterface,
		"this":      TokenTypeThis,
		"super":     TokenTypeSuper,
		"static":    TokenTypeStatic,
		"public":    TokenTypePublic,
		"private":   TokenTypePrivate,
		"new":       TokenTypeNew,
		"nil":       TokenTypeNil,
		"true":      TokenTypeTrue,
		"false":     TokenTypeFalse,
		"try":       TokenTypeTry,
		"catch":     TokenTypeCatch,
		"finally":   TokenTypeFinally,
		"throws":    TokenTypeThrows,
	}
}
func InitLexerFromString(lexer *Lexer, content string) {
	lexer.Content = []byte(content)
	InitLexerCommonOperations(lexer)
}
func InitLexerFromBytes(lexer *Lexer, content []byte) {
	lexer.Content = content
	InitLexerCommonOperations(lexer)
}

func (lexer *Lexer) ResetBytePos(i int) {
	lexer.BytePos = i
}

// Token 的 ToString() 方法
func (token *Token) ToString() string {
	return fmt.Sprintf("Line %d:%d  Type: %d, Str: %s", token.Line, token.Col, token.Kind, token.Str)
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

// 读出一个十六进制整数的 Token
func (lexer *Lexer) ReadHexadecimal() (*Token, *CoralError) {
	lexer.GoNextCharByStep(2) // skip '0x'
	str := "0x"

	for lexer.PeekChar().IsLegalHexadecimal() {
		str += string(lexer.PeekChar().Rune)
		lexer.GoNextChar()
	}
	return lexer.makeToken(TokenTypeHexadecimalInteger, str), nil
}

// 读出一个八进制整数的 Token
func (lexer *Lexer) ReadOctal() (*Token, *CoralError) {
	lexer.GoNextCharByStep(2) // 跳过 '0o'
	str := "0o"
	for lexer.PeekChar().IsLegalOctal() {
		str += string(lexer.PeekChar().Rune)
		lexer.GoNextChar()
	}
	return lexer.makeToken(TokenTypeOctalInteger, str), nil
}

// 读出一个二进制整数的 Token
func (lexer *Lexer) ReadBinary() (*Token, *CoralError) {
	lexer.GoNextCharByStep(2) // 跳过 '0o'
	str := "0b"
	for lexer.PeekChar().IsLegalBinary() {
		str += string(lexer.PeekChar().Rune)
		lexer.GoNextChar()
	}
	return lexer.makeToken(TokenTypeBinaryInteger, str), nil
}

// 读出一个十进制整数 或 小数/科学记数法 Token
func (lexer *Lexer) ReadDecimal(startFromZero bool) (*Token, *CoralError) {
	var str string
	hadPoint := false
	hadETag := false
	resultType := TokenTypeDecimalInteger

	if startFromZero {
		str += "0" // 记录一个 0 而抛弃所有其他无用的零
		lexer.GoNextChar()

		// 由十进制的特点，把前置的其他无用 '0' 全部抛弃
		for lexer.PeekChar().MatchRune('0') {
			lexer.GoNextChar()
		}
	}

	for {
		if lexer.PeekChar().IsLegalDecimal() {
			str += string(lexer.PeekChar().Rune)
			lexer.GoNextChar()
		} else if lexer.PeekChar().MatchRune('.') {
			// 如果是两个点连着，视为区间运算符
			if lexer.PeekNextChar(lexer.PeekChar().ByteLength).MatchRune('.') {
				break // 从这个点 此处断开，用已经获得的字符串 str 组成一个数值 token
			}

			if !hadPoint && !hadETag { // 读入小数点
				hadPoint = true
				resultType = TokenTypeFloat
				str += string(lexer.PeekChar().Rune)
				lexer.GoNextChar()
			} else {
				// 已经有了小数点，报错小数点重复
				return nil, NewCoralError("Syntax", "multiple decimal point!", LexFloatFormatError)
			}
		} else if lexer.PeekChar().MatchRune('e') {
			if !hadETag { // 读入 e 符号
				hadETag = true
				resultType = TokenTypeExponent
				str += string(lexer.PeekChar().Rune)
				lexer.GoNextChar()

				// 此时已经移过 'e'
				// 如果后方有 +/- 也一并读入
				if lexer.PeekChar().MatchRune('+') || lexer.PeekChar().MatchRune('-') {
					str += string(lexer.PeekChar().Rune)
					lexer.GoNextChar()
				}
			} else {
				// 科学记数法格式错误
				return nil, NewCoralError("Syntax", "incorrect format for scientific notation!", LexExponentFormatError)
			}
		} else {
			break // 不符合十进制整数、小数和科学记数法的格式条件
		}
	}

	// 如果科学记数法是 '0e' 开头，认为其无意义，抛出报错
	if len(str) >= 2 && str[0:2] == "0e" {
		return nil, NewCoralError("Syntax",
			"incorrect format for scientific notation! \nTips: Exponent starts from '0e' is meaningless.",
			LexExponentFormatError)
	}
	// 如果 str 以 '0' 起头 (只是 "0" 则不管) -> 要考虑去掉头部无用的 '0'
	if str != "0" && (str[0] == '0' && str[1] != '.') { // 不是 0. 起头的小数
		str = str[1:]
	}
	// 如果 str 最后一个字符是 'e' 也说明有问题
	if str[len(str)-1] == 'e' {
		return nil, NewCoralError("Syntax",
			"incorrect format for scientific notation! \nTips: Exponent can't just end with 'e'.",
			LexExponentFormatError)
	}

	return lexer.makeToken(resultType, str), nil
}

// 读出一个字符串，含转义字符的处理
func (lexer *Lexer) ReadString() (*Token, *CoralError) {
	var str string
	lexer.GoNextChar() // 移过当前的 '"' 双引号

	for !lexer.PeekChar().MatchRune('"') {
		if lexer.PeekChar().MatchRune('\\') { // 可能遇到转义字符
			switch lexer.PeekNextChar(lexer.PeekChar().ByteLength).Rune {
			case 'a':
				str += "\a"
				lexer.GoNextCharByStep(2)
			case 'b':
				str += "\b"
				lexer.GoNextCharByStep(2)
			case 't':
				str += "\t"
				lexer.GoNextCharByStep(2)
			case 'v':
				str += "\v"
				lexer.GoNextCharByStep(2)
			case 'n':
				str += "\n"
				lexer.GoNextCharByStep(2)
			case 'r':
				str += "\r"
				lexer.GoNextCharByStep(2)
			case 'f':
				str += "\f"
				lexer.GoNextCharByStep(2)
			case '"':
				str += "\""
				lexer.GoNextCharByStep(2)
			case 'u', 'U':
				// Unicode 需要是：\uXXXX 格式：
				lexer.GoNextCharByStep(2) // 移过当前的 '\u'
				unicodeBitCount := 0
				sUnicode := ""
				for lexer.PeekChar().IsLegalHexadecimal() {
					unicodeBitCount++
					sUnicode += string(lexer.PeekChar().Rune)
					lexer.GoNextChar()
				}
				if unicodeBitCount != 4 {
					// 说明不满 4 位，解码出错
					return nil, NewCoralError("Syntax", "(unicode error) 'unicodeEscape' codec can'Token decode bytes in position 0-3: truncated \\uXXXX escape", LexUnicodeEscapeFormatError)
				}
				gotUTF8Decoded := utils.UnicodeToUTF8(sUnicode, 4)
				str += gotUTF8Decoded
			case 'x':
				// Unicode 需要是：\xXX 格式：
				lexer.GoNextCharByStep(2) // 移过当前的 '\x'
				unicodeBitCount := 0
				sUnicode := ""
				for lexer.PeekChar().IsLegalHexadecimal() {
					unicodeBitCount++
					sUnicode += string(lexer.PeekChar().Rune)
					lexer.GoNextChar()
				}
				gotUTF8Decoded := utils.UnicodeToUTF8(sUnicode, 2)
				str += gotUTF8Decoded
			}
		} else {
			// 正常添加字符
			str += string(lexer.PeekChar().Rune)
			lexer.GoNextChar()
		}
	}

	lexer.GoNextChar() // 移过尾部的 '"' 双引号
	return lexer.makeToken(TokenTypeString, str), nil
}

// 读出一个字符，含转义字符的处理
func (lexer *Lexer) ReadRune() (*Token, *CoralError) {
	var str string
	lexer.GoNextChar() // 移过当前的 ' 双引号

	for !lexer.PeekChar().MatchRune('\'') {
		if lexer.PeekChar().MatchRune('\\') { // 可能遇到转义字符
			switch lexer.PeekNextChar(lexer.PeekChar().ByteLength).Rune {
			case 'a':
				str += "\a"
				lexer.GoNextCharByStep(2)
			case 'b':
				str += "\b"
				lexer.GoNextCharByStep(2)
			case 't':
				str += "\t"
				lexer.GoNextCharByStep(2)
			case 'v':
				str += "\v"
				lexer.GoNextCharByStep(2)
			case 'n':
				str += "\n"
				lexer.GoNextCharByStep(2)
			case 'r':
				str += "\r"
				lexer.GoNextCharByStep(2)
			case 'f':
				str += "\f"
				lexer.GoNextCharByStep(2)
			case '"':
				str += "\""
				lexer.GoNextCharByStep(2)
			case 'u', 'U':
				// Unicode 需要是：\uXXXX 格式：
				lexer.GoNextCharByStep(2) // 移过当前的 '\u'
				unicodeBitCount := 0
				sUnicode := ""
				for lexer.PeekChar().IsLegalHexadecimal() {
					unicodeBitCount++
					sUnicode += string(lexer.PeekChar().Rune)
					lexer.GoNextChar()
				}
				if unicodeBitCount != 4 {
					// 说明不满 4 位，解码出错
					return nil, NewCoralError("Syntax", "(unicode error) 'unicodeEscape' codec can'Token decode bytes in position 0-3: truncated \\uXXXX escape", LexUnicodeEscapeFormatError)
				}
				gotUTF8Decoded := utils.UnicodeToUTF8(sUnicode, 4)
				str += gotUTF8Decoded
			case 'x':
				// Unicode 需要是：\xXX 格式：
				lexer.GoNextCharByStep(2) // 移过当前的 '\x'
				unicodeBitCount := 0
				sUnicode := ""
				for lexer.PeekChar().IsLegalHexadecimal() {
					unicodeBitCount++
					sUnicode += string(lexer.PeekChar().Rune)
					lexer.GoNextChar()
				}
				gotUTF8Decoded := utils.UnicodeToUTF8(sUnicode, 2)
				str += gotUTF8Decoded
			}

		} else {
			// 正常添加字符
			str += string(lexer.PeekChar().Rune)
			lexer.GoNextChar()
		}
	}

	lexer.GoNextChar() // 移过末尾的单引号
	return lexer.makeToken(TokenTypeRune, str), nil
}

func (lexer *Lexer) ReadIdentifier() (*Token, *CoralError) {
	// 保证第一位不为数字
	firstRuneMatcher := regexp.MustCompile(`[0-9]`) // 第一个字符一定不会是 switch 条件上的操作符、空白符等
	restRuneMatcher := regexp.MustCompile(`[ \t\n;:,(){}\[\].=!*/%^|&><+\-'"]`)
	if firstRuneMatcher == nil {
		return nil, NewCoralError("Compiler", "RegexExp creating error!", CompilerRegexExpCreatingFailed)
	}
	// 读入第一个字符
	str := string(lexer.PeekChar().Rune)
	if firstRuneMatcher.MatchString(str) {
		return nil, NewCoralError("Syntax", "Digit can'Token be used for the first character of an identifier!", LexIdentifierFirstRuneCanNotBeDigit)
	}
	lexer.GoNextChar()
	for lexer.BytePos < len(lexer.Content) &&
		restRuneMatcher != nil &&
		!restRuneMatcher.MatchString(string(lexer.PeekChar().Rune)) {
		str += string(lexer.PeekChar().Rune)
		lexer.GoNextChar()
	}

	if keywordType, isKeyword := lexer.KeywordMap[str]; isKeyword {
		return lexer.makeToken(keywordType, str), nil
	} // 如果是关键字 则 返回对应关键字的 Token 类型
	return lexer.makeToken(TokenTypeIdentifier, str), nil
}

// 跳过块注释
func (lexer *Lexer) SkipBlockComment() *CoralError {
	lexer.GoNextCharByStep(2) // 跳过 "/*"
	nested := 1               // 初始嵌套层次为 1

	current := lexer.PeekChar()
	next := lexer.PeekNextChar(current.ByteLength)
	for {
		lexer.GoNextChar() // 首要任务 跳过当前注释内容

		// 然后更新 current 和 next
		current = lexer.PeekChar()
		next = lexer.PeekNextChar(current.ByteLength)

		if current.MatchRune('/') && next.MatchRune('*') {
			nested++
			if nested > 5 {
				return NewCoralError("Syntax", "too many nested levels in a block comment!", LexBlockCommentTooNested)
			}
		}
		if current.MatchRune('*') && next.MatchRune('/') {
			nested--
			lexer.GoNextCharByStep(2) // 移动 2 位移过 "*/"
		}
		if nested == 0 {
			break
		}
	}

	return nil
}

// 跳过行注释
func (lexer *Lexer) SkipLineComment() {
	// 直到换行符
	for !lexer.PeekChar().MatchRune('\n') {
		lexer.GoNextChar()
	}
}

// 产出 Token，词法分析器的行号也移动字面值 s 的长度
func (lexer *Lexer) makeToken(t TokenType, s string) *Token {
	lexer.Col += utf8.RuneCountInString(s)
	// s 这个字符串的长度就是其中 UTF8 字符个数的长度
	return &Token{
		Line: lexer.Line,
		Col:  lexer.Col,
		Kind: t,
		Str:  s,
	}
}

// 词法分析器获取下一个 Token
func (lexer *Lexer) GetNextToken(avoidAngleConfusing bool) (*Token, *CoralError) {
	for lexer.BytePos < len(lexer.Content) {
		c := lexer.PeekChar()
		switch c.Rune {
		default:
			if c.IsLegalDecimal() {
				return lexer.ReadDecimal(false)
			}
			return lexer.ReadIdentifier()
		case '\t', ' ':
			lexer.Col += 1
			lexer.GoNextChar() // skip whitespace
		case '\n':
			lexer.Line++
			lexer.Col = 1
			lexer.GoNextChar() // skip
		case ';':
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeSemi, ";"), nil
		case ',':
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeComma, ","), nil
		case ':':
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeColon, ":"), nil
		case '(':
			lexer.GoNextChar()
			lexer.ParenCount++
			return lexer.makeToken(TokenTypeLeftParen, "("), nil
		case ')':
			lexer.GoNextChar()
			lexer.ParenCount--
			return lexer.makeToken(TokenTypeRightParen, ")"), nil
		case '{':
			lexer.GoNextChar()
			lexer.BraceCount++
			return lexer.makeToken(TokenTypeLeftBrace, "{"), nil
		case '}':
			lexer.GoNextChar()
			lexer.BraceCount--
			return lexer.makeToken(TokenTypeRightBrace, "}"), nil
		case '[':
			lexer.GoNextChar()
			lexer.BracketCount++
			return lexer.makeToken(TokenTypeLeftBracket, "["), nil
		case ']':
			lexer.GoNextChar()
			lexer.BracketCount--
			return lexer.makeToken(TokenTypeRightBracket, "]"), nil
		case '.':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('.') {
				if lexer.PeekNextCharByStep(c.ByteLength, 2).MatchRune('.') {
					lexer.GoNextCharByStep(3)
					return lexer.makeToken(TokenTypeEllipsis, "..."), nil
				}
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeDoubleDot, ".."), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeDot, "."), nil
		case '~':
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeWavy, "~"), nil
		case '@':
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeAlpha, "@"), nil
		case '=':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeDoubleEqual, "=="), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeEqual, "="), nil
		case '!':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeBangEqual, "!="), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeBang, "!"), nil
		case '*':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('*') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeDoubleStar, "**"), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeStarEqual, "*="), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeStar, "*"), nil
		case '/':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeSlashEqual, "/="), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('*') {
				err := lexer.SkipBlockComment()
				if err != nil {
					return nil, err // 可能的块注释略过时出错
				}
				continue
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('/') {
				lexer.SkipLineComment()
				continue
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeSlash, "/"), nil
		case '%':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypePercentEqual, "%="), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypePercent, "%"), nil
		case '^':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeCaretEqual, "^="), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeCaret, "^"), nil
		case '&':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('&') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeDoubleAmpersand, "&&"), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeAmpersandEqual, "&="), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeAmpersand, "&"), nil
		case '|':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('|') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeDoubleVertical, "||"), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeVerticalEqual, "|="), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeVertical, "|"), nil
		case '<':
			if !avoidAngleConfusing && lexer.PeekNextChar(c.ByteLength).MatchRune('<') {
				if lexer.PeekNextCharByStep(c.ByteLength, 2).MatchRune('=') {
					lexer.GoNextCharByStep(3)
					return lexer.makeToken(TokenTypeDoubleLeftAngleEqual, "<<="), nil
				}
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeDoubleLeftAngle, "<<"), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeLeftAngleEqual, "<="), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('-') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeLeftArrow, "<-"), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeLeftAngle, "<"), nil
		case '>':
			if !avoidAngleConfusing && lexer.PeekNextChar(c.ByteLength).MatchRune('>') {
				if lexer.PeekNextCharByStep(c.ByteLength, 2).MatchRune('=') {
					lexer.GoNextCharByStep(3)
					return lexer.makeToken(TokenTypeDoubleRightAngleEqual, ">>="), nil
				}
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeDoubleRightAngle, ">>"), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeRightAngleEqual, ">="), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeRightAngle, ">"), nil
		case '+':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('+') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeDoublePlus, "++"), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypePlusEqual, "+="), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypePlus, "+"), nil
		case '-':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('-') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeDoubleMinus, "--"), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('=') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeMinusEqual, "-="), nil
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('>') {
				lexer.GoNextCharByStep(2)
				return lexer.makeToken(TokenTypeRightArrow, "->"), nil
			}
			lexer.GoNextChar()
			return lexer.makeToken(TokenTypeMinus, "-"), nil
		case '0':
			if lexer.PeekNextChar(c.ByteLength).MatchRune('x') {
				return lexer.ReadHexadecimal()
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('o') {
				return lexer.ReadOctal()
			} else if lexer.PeekNextChar(c.ByteLength).MatchRune('b') {
				return lexer.ReadBinary()
			}
			return lexer.ReadDecimal(true)
		case '"':
			return lexer.ReadString()
		case '\'':
			return lexer.ReadRune()
		}
	}

	if lexer.ParenCount > 0 {
		return nil, NewCoralError("Syntax", "Unclosed parentheses '(' !", LexParenthesesUnclosed)
	}
	if lexer.BracketCount > 0 {
		return nil, NewCoralError("Syntax", "Unclosed bracket '[' !", LexBracketUnclosed)
	}
	if lexer.ParenCount > 0 {
		return nil, NewCoralError("Syntax", "Unclosed brace '{' !", LexBraceUnclosed)
	}

	return nil, nil
}
