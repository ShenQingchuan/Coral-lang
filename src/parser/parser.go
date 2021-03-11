package parser

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/exception"
	. "coral-lang/src/lexer"
	. "coral-lang/src/utils"
	"fmt"
	"strings"
)

type Parser struct {
	Lexer *Lexer

	LastToken    *Token
	CurrentToken *Token

	ErrCount  int
	WarnCount int
}

func CoralCompileErrorWithPos(parser *Parser, c *CoralCompileError) {
	if parser.LastToken != nil {
		fmt.Print("\n" + Bold(Green(fmt.Sprintf("* line %d:%d ", parser.LastToken.Line, parser.LastToken.Col))))
	}
	fmt.Println(c.Err)

	// 打印错误代码所在行以及附近两行
	lines := strings.Split(string(parser.Lexer.Content), "\n")
	var startLineIndex int
	if parser.LastToken.Line == 1 {
		startLineIndex = 0
	} else {
		startLineIndex = parser.LastToken.Line - 2
	}
	for i := 0; i < 3 && (startLineIndex+i) < len(lines); i++ {
		fmt.Print(Yellow(fmt.Sprintf("%4d", startLineIndex+i+1)))
		fmt.Printf("| %s\n", lines[startLineIndex+i])
		if startLineIndex+i == parser.LastToken.Line-1 {
			trimmed := false
			for k := 0; k < 6; k++ {
				fmt.Print(" ")
			}
			for j := 0; j < parser.LastToken.Col; j++ {
				if !trimmed {
					if lines[startLineIndex+i][j] == ' ' {
						fmt.Print(" ")
						continue
					} else if lines[startLineIndex+i][j] == '\t' {
						fmt.Print("  ")
						continue
					} else {
						trimmed = true
					}
				}

				fmt.Print(Yellow("∼"))
			}
			fmt.Print(Red("^") + "\n")
		}
	}

	parser.ErrCount++
}
func CoralCompileWarningWithPos(parser *Parser, msg string) {
	if parser.LastToken != nil {
		fmt.Print("\n" + Green(fmt.Sprintf("* line %d:%d ", parser.LastToken.Line, parser.LastToken.Col)))
	}
	CoralCompileWarning(msg)
	parser.WarnCount++
}

func (parser *Parser) InitFromBytes(content []byte) {
	parser.Lexer = new(Lexer)
	parser.Lexer.InitFromBytes(content)
	parser.PeekNextToken() // 统一获取到第一个 Token
}
func (parser *Parser) InitFromString(content string) {
	parser.Lexer = new(Lexer)
	parser.Lexer.InitFromString(content)
	parser.PeekNextToken() // 统一获取到第一个 Token
}
func (parser *Parser) AssertCurrentTokenIs(tokenType TokenType, expected string, situation string) bool {
	if !parser.MatchCurrentTokenType(tokenType) {
		CoralCompileErrorWithPos(parser, NewCoralError("Syntax",
			fmt.Sprintf("expected %s %s!", expected, situation), ParsingUnexpected))
		return false
	}
	parser.PeekNextToken()
	return true
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

	fmt.Println("\n" + Yellow(fmt.Sprintf("(Parser: %d error, %d warning)", parser.ErrCount, parser.WarnCount)))
	return program
}
