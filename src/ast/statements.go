package ast

import (
	. "coral-lang/src/lexer"
)

// 定义简单语句的种类来区分
const (
	SimpleStmtTypeExpression = iota
	SimpleStmtIncDecStmt
	SimpleStmtTypeVariableDecl
	SimpleStmtTypeBreak
	SimpleStmtTypeContinue
	SimpleStmtTypeReturn
	SimpleStmtTypeAssign
)

// 单个变量定义的赋值部分
type VarDeclElement struct {
	Variable  *Identifier // 定义的变量标识符
	Type      *TypeDescription
	InitValue *Expression // 赋予的初始值（是个表达式）
}

// 返回语句节点
type ReturnStatement struct {
	Token      *Token
	Expression *Expression
}

func (it *ReturnStatement) NodeType() string {
	return "Simple_Statement_Return"
}
func (it *ReturnStatement) SimpleStatementNodeType() int {
	return SimpleStmtTypeReturn
}

// 循环中断语句节点
type BreakStatement struct {
	Token *Token
}

func (it *BreakStatement) NodeType() string {
	return "Simple_Statement_Break"
}
func (it *BreakStatement) SimpleStatementNodeType() int {
	return SimpleStmtTypeBreak
}

// 循环继续语句节点
type ContinueStatement struct {
	Token *Token
}

func (it *ContinueStatement) NodeType() string {
	return "Simple_Statement_Continue"
}
func (it *ContinueStatement) SimpleStatementNodeType() int {
	return SimpleStmtTypeContinue
}

// 自增或自减语句节点
type IncDecStatement struct {
	Expression *Expression
	Operator   *Token
}

func (it *IncDecStatement) NodeType() string {
	var stmtType string
	if it.Operator.Kind == TokenTypeDoublePlus {
		stmtType = "Simple_Statement_Self_Increase"
	} else {
		stmtType = "Simple_Statement_Self_Decrease"
	}
	return stmtType
}
func (it *IncDecStatement) SimpleStatementNodeType() int {
	return SimpleStmtIncDecStmt
}

// 变量定义语句节点
type VarDeclStatement struct {
	Mutable      bool              // 用于区分 var 和 val
	declarations []*VarDeclElement // 可能有多个变量定义
}

func (it *VarDeclStatement) NodeType() string {
	var declType string
	if it.Mutable {
		declType = "kind: var"
	} else {
		declType = "kind: val"
	}
	return "Simple_Statement_Variable_Declaration, " + declType
}
func (it *VarDeclStatement) SimpleStatementNodeType() int {
	return SimpleStmtTypeVariableDecl
}

// 表达式语句节点
type ExpressionStatement struct {
	Expression *Expression
}

func (it *ExpressionStatement) NodeType() string {
	return "Simple_Statement_Expression"
}
func (it *ExpressionStatement) SimpleStatementNodeType() int {
	return SimpleStmtTypeExpression
}

// 简单表达式定义
type SimpleStatement interface {
	Node
	SimpleStatementNodeType() int
}
