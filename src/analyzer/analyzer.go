package analyzer

import (
	. "container/list"
	. "coral-lang/src/ast"
	. "coral-lang/src/parser"
)

// 国防科技大学：《编译原理》教学 [P37 - 符号表与作用域分析]
// https://www.bilibili.com/video/BV1QE411o73W?p=37
type Analyzer struct {
	Parser *Parser

	// 变量作用域分析
	Display    *List // 显示表
	BlockTable *List // 程序体表
	NameTable  *List // 名字表
}

func InitAnalyzer(analyzer *Analyzer) {
	analyzer.Parser.ParseProgram()      // 执行语法解析得到 AST
	analyzer.BuiltInTypeNameInjection() // 注入内建类型名称
	analyzer.AppendBlockTable()         // 初始化第一张区块表
	analyzer.AppendDisplay()            // 初始化第一张显示表

	// 开始执行对 AST 的语法分析
	analyzer.AnalyzeTree()
}
func InitAnalyzerFromBytes(analyzer *Analyzer, content []byte) {
	analyzer.Parser = new(Parser)
	InitParserFromBytes(analyzer.Parser, content)
	InitAnalyzer(analyzer)
}
func InitAnalyzerFromString(analyzer *Analyzer, content string) {
	analyzer.Parser = new(Parser)
	InitParserFromString(analyzer.Parser, content)
	InitAnalyzer(analyzer)
}
func (analyzer *Analyzer) BuiltInTypeNameInjection() {
	builtInTypes := []string{
		"int", "int8", "int16", "int64",
		"uint", "uint8", "uint16", "uint64",
		"float", "double",
		"bool", "rune", "String",
	}
	var prev *NameTableItem = nil
	for _, typeName := range builtInTypes {
		item := &NameTableItem{
			Name:  typeName,
			Kind:  NameKindBuiltInType,
			Level: 0,
			Link:  prev,
		}
		analyzer.NameTable.PushBack(item)
		prev = item // 将会作为下一元素的前驱
	}
}
func (analyzer *Analyzer) AppendBlockTable() {
	analyzer.BlockTable.PushBack(&BlockTableItem{
		First: analyzer.NameTable.Front().Value.(*NameTableItem),
		Last:  analyzer.NameTable.Back().Value.(*NameTableItem),
	}) // 初始化名字表后直接记录
}
func (analyzer *Analyzer) AppendDisplay() {
	// 初始化显示表是：立即记录刚加入的第一个 BlockTableItem
	analyzer.Display.PushBack(analyzer.BlockTable.Front())
}

func (analyzer *Analyzer) AnalyzeTree() {
	for _, stmtNode := range analyzer.Parser.Program.Root {
		switch stmtNode.StatementNodeType() {
		case StatementTypeSimple:
		case StatementTypePackage:
		case StatementTypeImport:
			importStmtNode := stmtNode.(ImportStatement)
			analyzer.AnalyzeImportStatement(importStmtNode) // 针对 import 语句的三种情况进行分析
		case StatementTypeEnum:
		case StatementTypeBlock:
		case StatementTypeTryCatch:
		case StatementTypeIf:
		case StatementTypeSwitch:
		case StatementTypeWhile:
		case StatementTypeFor:
		case StatementTypeEach:
		case StatementTypeFunctionDecl:
		case StatementTypeClassDecl:
		case StatementTypeInterfaceDecl:
		}
	}
}

func (analyzer *Analyzer) AnalyzeImportStatement(node ImportStatement) {
	switch node.ImportStatementNodeType() {
	case ImportStatementTypeSingleGlobal:
		singleGlobalImportStmtNode := node.(*SingleGlobalImportStatement)
		packageName := singleGlobalImportStmtNode.Element.ModuleName.GetName()
		if isStandardLibrary(packageName) {
			analyzer.LoadStandardLibrary(packageName)
		} else {

		}
	case ImportStatementTypeSingleFrom:
		singleFromStmtNode := node.(*SingleFromImportStatement)
		packageName := singleFromStmtNode.From.GetName()
		moduleElement := singleFromStmtNode.Element
		if isStandardLibrary(packageName) {
			analyzer.LoadNameFromStandardLibrary(packageName, moduleElement)
		} else {

		}
	case ImportStatementTypeList:
		// 同上，需遍历列表中的引入模块名字
		importList := node.(*ListImportStatement)
		packageName := importList.From.GetName()
		if isStandardLibrary(packageName) {
			for _, moduleElement := range importList.Elements {
				analyzer.LoadNameFromStandardLibrary(packageName, moduleElement)
			}
		} else {

		}
	}
}

func (analyzer *Analyzer) LoadStandardLibrary(packageName string) {
	// TODO: 将标准库中的名字引入
}

func (analyzer *Analyzer) LoadNameFromStandardLibrary(packageName string, moduleElement *ImportElement) {
	// TODO: 如果在该标准库中有该名字的模块则录入到名字表，若没有就报错
}
