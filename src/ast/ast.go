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
type INode interface {
	NodeType() string
}

// Statement 为所有语句节点定义了接口
type IStatement interface {
	INode
	StatementNode()
}

// Expression 为所有表达式节点定义了接口
type IExpression interface {
	INode
	ExpressionNode()
}

// Program 为语法树根节点，每一个 .coral 源代码文件都被视为一整段程序，以一个语句的切片表达
type Program struct {
	Root []IStatement
}

// 标识符节点
type Identifier struct {
	Name *Token
}

func (it *Identifier) NodeType() string {
	return "Identifier: " + it.Name.Value
}

// 类型的名称节点
type TypeName struct {
	Link []*Identifier
}

func (it *TypeName) NodeType() string {
	var typeName string
	for i, id := range it.Link {
		typeName += id.Name.Value
		if i != len(it.Link)-1 {
			typeName += "."
		}
	} // 拼接 x.y.z 形状的类型名
	return "Type_Name: " + typeName
}

// 类型字面量节点
type TypeLit interface {
	// TODO: typeLit ::= ('[' typeDescription ']') | (typeName '<' typeName (',' typeName)* '>')
}

// 类型标注节点
type TypeDescription struct {
	// TODO: typeDescription ::= typeName | typeLit
}

// 单个变量定义的赋值部分
type VarDeclElement struct {
	Variable  *Identifier // 定义的变量标识符
	InitValue IExpression // 赋予的初始值（是个表达式）
}

// 变量定义语句
type VarDeclStatement struct {
	Mutable      bool             // 用于区分 var 和 val
	declarations []VarDeclElement // 可能有多个变量定义
}

func (it *VarDeclStatement) NodeType() string {
	var declType string
	if it.Mutable {
		declType = "kind: var"
	} else {
		declType = "kind: val"
	}
	return "Variable_Declaration_Statement, " + declType
}
