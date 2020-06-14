package analyzer

import (
	. "container/list"
	. "coral-lang/src/parser"
)

// 国防科技大学：《编译原理》教学 [P37 - 符号表与作用域分析]
// https://www.bilibili.com/video/BV1QE411o73W?p=37
type Analyzer struct {
	parser *Parser

	// 变量作用域分析
	display    *List // 显示表
	blockTable *List // 程序体表
	nameTable  *List // 名字表
}

func InitAnalyzer(analyzer *Analyzer) {

}
func InitAnalyzerFromBytes(analyzer *Analyzer, content []byte) {
	analyzer.parser = new(Parser)
	InitParserFromBytes(analyzer.parser, content)
	InitAnalyzer(analyzer)
}
func InitAnalyzerFromString(analyzer *Analyzer, content string) {
	analyzer.parser = new(Parser)
	InitParserFromString(analyzer.parser, content)
	InitAnalyzer(analyzer)
}
