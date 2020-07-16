package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
	. "coral-lang/src/utils"
	"fmt"
	"os"
)

type Parser struct {
	Lexer   *Lexer

	LastToken    *Token
	CurrentToken *Token
}

func CoralErrorCrashHandlerWithPos(parser *Parser, c *CoralError) {
	if parser.LastToken != nil {
		fmt.Print("\n" + Green(fmt.Sprintf("* line %d:%d ", parser.LastToken.Line, parser.LastToken.Col)))
	}
	fmt.Println(c.Err)
	fmt.Println("* " + Cyan(fmt.Sprintf("Error code: %d", c.ErrEnum)))
	os.Exit(c.ErrEnum)
}
func CoralCompileWarningWithPos(parser *Parser, msg string) {
	if parser.LastToken != nil {
		fmt.Print("\n" + Green(fmt.Sprintf("* line %d:%d ", parser.LastToken.Line, parser.LastToken.Col)))
	}
	CoralCompileWarning(msg)
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
		CoralErrorCrashHandlerWithPos(parser, NewCoralError("Syntax",
			fmt.Sprintf("expected %s %s!", expected, situation), ParsingUnexpected))
	}
}
func (parser *Parser) PeekNextToken() {
	token, err := parser.Lexer.GetNextToken(true)
	if err != nil {
		CoralErrorCrashHandler(err)
	}

	parser.LastToken = parser.CurrentToken
	parser.CurrentToken = token
}
func (parser *Parser) PeekNextTokenAvoidAngleConfusing() {
	token, err := parser.Lexer.GetNextToken(false)
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

func (parser *Parser) ParseProgram() *Program {
	program := new(Program)
	for stmt := parser.ParseStatement(); stmt != nil; stmt = parser.ParseStatement() {
		// stmt 为 nil 的情况中其实早已被 CoralErrorCrashHandler 处理并退出了
		program.Root = append(program.Root, stmt)
	}

	if _, isPkgStmt := program.Root[0].(*PackageStatement); !isPkgStmt {
		CoralErrorCrashHandlerWithPos(parser, NewCoralError("Syntax",
			"expected a package name for a source file as the first statement!", NoPackageNameDefinition))
	}

	return program
}
