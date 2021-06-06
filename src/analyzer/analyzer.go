package analyzer

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/lexer"
	. "coral-lang/src/parser"
)

const (
	IdentifierSymbolKind = iota
	TypeSymbolKind
	EnumSymbolKind
)

type Symbol struct {
	Token *Token // 符号相应 token
}
type ISymbol interface {
	GetToken() *Token
	GetKind() int
}

// 标识符号表
type IdSymbol struct {
	*Symbol
	Type *TypeSymbol // 符号所属的类型
}

func (idSymbol *IdSymbol) GetToken() *Token {
	return idSymbol.Symbol.Token
}
func (idSymbol *IdSymbol) GetKind() int {
	return IdentifierSymbolKind
}

// 类型符号表
type TypeSymbol struct {
	*Symbol
	IsFn        bool            // @private 是否为函数
	Signature   *Signature      // @private 函数签名
	Description TypeDescription // 类型描述
	DescType    int             // 类型描述的枚举
}

func (typeSymbol *TypeSymbol) GetToken() *Token {
	return typeSymbol.Symbol.Token
}
func (typeSymbol *TypeSymbol) GetKind() int {
	return TypeSymbolKind
}

// 枚举符号
type EnumSymbol struct {
	*Symbol
	CollectionName string
	ElementsMap    map[string]*EnumElement
}

func (enumSymbol *EnumSymbol) GetToken() *Token {
	return enumSymbol.Symbol.Token
}
func (enumSymbol *EnumSymbol) GetKind() int {
	return EnumSymbolKind
}

type BlockScope struct {
	OuterScope *BlockScope // 外层区块

	/* Q: 为什么这里要把多种符号抽象成一个统一接口来继承
	 * A: 因为符号名称要尽可能地保持无重复、冲突 */
	SymbolMap map[string]ISymbol
}

type Analyzer struct {
	parser *Parser // @private 语法解析器

	RootScope    *BlockScope // 顶层区块
	CurrentScope *BlockScope // 遍历区块层级时的指针
	Ast          *Program    // AST
}

func (analyzer *Analyzer) InitAnalyzerCommon() {
	analyzer.Ast = analyzer.parser.ParseProgram() // 获取抽象语法树
	rootScope := new(BlockScope)
	rootScope.SymbolMap = make(map[string]ISymbol)

	analyzer.RootScope = rootScope
	analyzer.CurrentScope = analyzer.RootScope
}
func (analyzer *Analyzer) InitAnalyzerFromString(content string) {
	parser := new(Parser)
	parser.InitFromString(content)
	analyzer.parser = parser
	analyzer.InitAnalyzerCommon()
}
func (analyzer *Analyzer) InitAnalyzerFromBytes(content []byte) {
	parser := new(Parser)
	parser.InitFromBytes(content)
	analyzer.parser = parser
	analyzer.InitAnalyzerCommon()
}
func (analyzer *Analyzer) EnterNewBlockScope() {
	newScope := new(BlockScope)
	newScope.OuterScope = analyzer.CurrentScope
	analyzer.CurrentScope = newScope
}
func (analyzer *Analyzer) LeaveCurrentBlockScope() {
	analyzer.CurrentScope = analyzer.CurrentScope.OuterScope
}
