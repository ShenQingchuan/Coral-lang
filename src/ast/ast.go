package ast

import (
	. "coral-lang/src/lexer"
)

// Package ast 实现了抽象语法树，勇于表达语法解析结果
// 上承词法分析，下启语义分析

/*
  !! Golang 基于接口、结构体的类型系统：
  要想理解此套心智模型，需要首先明确结构体是数据的类型、接口是方法的类型
  要想确定一个传统 OOP 意义上的「类」固然需要两者兼备
  所以我们会用 IA 表达一个接口，其实就是一个方法列表
  然后定义结构体 A 去实现这个接口
*/

// Node 为语法树中的所有节点定义了接口
type Node interface {
	NodeType() string
}

// Program 为语法树根节点，每一个 .coral 源代码文件都被视为一整段程序，以一个语句的切片表达
type Program struct {
	Root []Statement
}

// 标识符节点
type Identifier struct {
	Token *Token
}

func (it *Identifier) NodeType() string {
	return "Identifier: " + it.Token.Str
}
