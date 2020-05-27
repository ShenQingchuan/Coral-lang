package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
	"fmt"
	"os"
)

type Parser struct {
	Lexer   *Lexer
	Program *Program

	LastToken    *Token
	CurrentToken *Token
}

func CoralErrorCrashHandlerWithPos(parser *Parser, c *CoralError) {
	if parser.CurrentToken != nil {
		fmt.Printf("\n* line %d:%d ", parser.CurrentToken.Line, parser.CurrentToken.Col)
	} else if parser.LastToken != nil {
		fmt.Printf("\n* line %d:%d ", parser.LastToken.Line, parser.LastToken.Col)
	}
	fmt.Println(c.Err)
	fmt.Printf("* Error code: %d", c.ErrEnum)
	os.Exit(c.ErrEnum)
}

func InitParserFromBytes(parser *Parser, content []byte) {
	parser.Lexer = new(Lexer)
	InitLexerFromBytes(parser.Lexer, content)
	parser.PeekNextToken() // 统一获取到第一个 Token
}
func InitParserFromString(parser *Parser, content string) {
	parser.Lexer = new(Lexer)
	InitLexerFromString(parser.Lexer, content)
	parser.PeekNextToken() // 统一获取到第一个 Token
}
func (parser *Parser) AssertCurrentTokenIs(tokenType TokenType, expected string, situation string) {
	if parser.MatchCurrentTokenType(tokenType) {
		parser.PeekNextToken()
	} else {
		CoralErrorCrashHandlerWithPos(parser, NewCoralError("Compile",
			fmt.Sprintf("expected a %s %s!", expected, situation), ParsingUnexpected))
	}
}
func (parser *Parser) PeekNextToken() {
	token, err := parser.Lexer.GetNextToken()
	if err != nil {
		CoralErrorCrashHandler(err)
	}

	parser.LastToken = parser.CurrentToken
	parser.CurrentToken = token
}
func (parser *Parser) GetCurrentTokenPos() string {
	return fmt.Sprintf("line %d:%d: ", parser.CurrentToken.Line, parser.CurrentToken.Col)
}
func (parser *Parser) MatchCurrentTokenType(tokenType TokenType) bool {
	if parser.CurrentToken != nil {
		return parser.CurrentToken.Kind == tokenType
	}
	return false
}

func (parser *Parser) ParseProgram() {
	for stmt := parser.ParseStatement(); stmt != nil; stmt = parser.ParseStatement() {
		// stmt 为 nil 的情况中其实早已被 CoralErrorCrashHandler 处理并退出了
		parser.Program.Root = append(parser.Program.Root, stmt)
	}
}
