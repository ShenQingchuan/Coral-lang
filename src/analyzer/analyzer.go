package analyzer

import (
	. "coral-lang/src/ast"
	. "coral-lang/src/lexer"
	. "coral-lang/src/parser"
)

// 类型符号表
type TypeSymbol struct {
	IsFn        bool            // @private 是否为函数
	Signature   *Signature      // @private 函数签名
	Description TypeDescription // 类型描述
	DescType    int             // 类型描述的枚举
}

// 标识符号表
type IdSymbol struct {
	Token *Token      // 符号相应 token
	Type  *TypeSymbol // 符号所属的类型
}

type BlockScope struct {
	OuterScope    *BlockScope // 外层区块
	IdSymbolMap   map[string]IdSymbol
	TypeSymbolMap map[string]TypeSymbol
}

type Analyzer struct {
	parser *Parser // @private 语法解析器

	RootScope    *BlockScope // 顶层区块
	CurrentScope *BlockScope // 遍历区块层级时的指针
	Ast          *Program    // AST
}

func (analyzer *Analyzer) InitAnalyzerCommon() {
	analyzer.Ast = analyzer.parser.ParseProgram() // 获取抽象语法树
	analyzer.RootScope = new(BlockScope)
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
