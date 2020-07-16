package analyzer

import (
	. "container/list"
	. "coral-lang/src/ast"
)

// 国防科技大学：《编译原理》教学 [P37 - 符号表与作用域分析]
// https://www.bilibili.com/video/BV1QE411o73W?p=37
type Analyzer struct {
	// 变量作用域分析
	Display    *List // 显示表
	BlockTable *List // 程序体表
	NameTable  *List // 名字表
}

func InitAnalyzer(analyzer *Analyzer) {
	analyzer.BuiltInTypeNameInjection() // 注入内建类型名称
	analyzer.AppendBlockTable()         // 初始化第一张区块表
	analyzer.AppendDisplay()            // 初始化第一张显示表
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
	}) // 初始化 "名字表" 时直接插入一条新记录
}
func (analyzer *Analyzer) AppendDisplay() {
	// 初始化 "显示层表" 时：立即记录刚加入的第一个 BlockTableItem
	analyzer.Display.PushBack(analyzer.BlockTable.Front())
}

