package ast

import (
	. "coral-lang/src/lexer"
)

// 定义所有语句的种类来区分
const (
	StatementTypeSimple = iota
	StatementTypeImport
	StatementTypeEnum
	StatementTypeBlock
	StatementTypeTryCatch
	StatementTypeIf
	StatementTypeSwitch
	StatementTypeWhile
	StatementTypeFor
	StatementTypeEach
	StatementTypeFunctionDecl
	StatementTypeClassDecl
)

// 定义引入外部模块语句的种类来区分
const (
	ImportStatementTypeSingle = iota
	ImportStatementTypeList
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

// 定义赋值语句的种类来区分
const (
	AssignStmtTypeList = iota
	AssignStmtTypeMix
)

// 定义条件语句匹配项的种类来区分
const (
	SwitchStatementTypeNormal = iota
	SwitchStatementTypeRange
)

// 定义类成员的种类来区分
const (
	ClassMemberTypeVar = iota
	ClassMemberTypeMethod
)

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
func (it *ReturnStatement) StatementNodeType() int {
	return StatementTypeSimple
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
func (it *BreakStatement) StatementNodeType() int {
	return StatementTypeSimple
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
func (it *ContinueStatement) StatementNodeType() int {
	return StatementTypeSimple
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
func (it *IncDecStatement) StatementNodeType() int {
	return StatementTypeSimple
}

// 单个变量定义的赋值部分
type VarDeclElement struct {
	Variable  *Identifier // 定义的变量标识符
	Type      *TypeDescription
	InitValue *Expression // 赋予的初始值（是个表达式）
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
func (it *VarDeclStatement) StatementNodeType() int {
	return StatementTypeSimple
}

// 语句语句节点
type ExpressionStatement struct {
	Expression *Expression
}

func (it *ExpressionStatement) NodeType() string {
	return "Simple_Statement_Expression"
}
func (it *ExpressionStatement) SimpleStatementNodeType() int {
	return SimpleStmtTypeExpression
}
func (it *ExpressionStatement) StatementNodeType() int {
	return StatementTypeSimple
}

// 同句多赋值语句定义的
type ListAssignStatement struct {
	Token   *Token // Token: '='
	Targets []*Expression
	Values  []*Expression
}

func (it *ListAssignStatement) NodeType() string {
	return "List_Assign_Statement"
}
func (it *ListAssignStatement) AssignStatementNodeType() int {
	return AssignStmtTypeList
}
func (it *ListAssignStatement) SimpleStatementNodeType() int {
	return SimpleStmtTypeAssign
}
func (it *ListAssignStatement) StatementNodeType() int {
	return StatementTypeSimple
}

// 混合赋值语句定义的
type MixAssignStatement struct {
	Token    *Token // 该混合赋值运算符
	Variable *OperandName
	Value    *Expression
}

func (it *MixAssignStatement) NodeType() string {
	return "Mix_Assign_Statement"
}
func (it *MixAssignStatement) AssignStatementNodeType() int {
	return AssignStmtTypeMix
}
func (it *MixAssignStatement) SimpleStatementNodeType() int {
	return SimpleStmtTypeAssign
}
func (it *MixAssignStatement) StatementNodeType() int {
	return StatementTypeSimple
}

// 赋值语句定义
type AssignStatement interface {
	SimpleStatement
	AssignStatementNodeType() int
}

// 简单语句定义
type SimpleStatement interface {
	Statement
	SimpleStatementNodeType() int
}

// 块语句节点
type BlockStatement struct {
	Statements []*Statement
}

func (it *BlockStatement) NodeType() string {
	return "Block_Statement"
}
func (it *BlockStatement) StatementNodeType() int {
	return StatementTypeBlock
}

// 模块名称节点
type ModuleName struct {
	NameUnits []*Identifier
}

func (it *ModuleName) GetFullModuleName() string {
	var fullModuleName string
	for i, id := range it.NameUnits {
		fullModuleName += id.Name.Str
		if i != len(it.NameUnits)-1 {
			fullModuleName += "."
		}
	}
	return fullModuleName
}

// 引入语句单元
type ImportElement struct {
	ModuleName *ModuleName
	As         *Identifier
}

func (it *ImportElement) NodeType() string {
	return "Import_Element"
}

// 单个引入语句
type SingleImportStatement struct {
	From    *ModuleName
	Element *ImportElement
}

func (it *SingleImportStatement) NodeType() string {
	return "Single_Import_Statement"
}
func (it *SingleImportStatement) ImportStatementNodeType() int {
	return ImportStatementTypeSingle
}
func (it *SingleImportStatement) StatementNodeType() int {
	return StatementTypeImport
}

// 集合引入语句
type ListImportStatement struct {
	From    *ModuleName
	Element *ImportElement
}

func (it *ListImportStatement) NodeType() string {
	return "List_Import_Statement"
}
func (it *ListImportStatement) ImportStatementNodeType() int {
	return ImportStatementTypeList
}
func (it *ListImportStatement) StatementNodeType() int {
	return StatementTypeImport
}

// 引入外部模块语句定义
type ImportStatement interface {
	Node
	ImportStatementNodeType() int
}

// 枚举单元
type EnumElement struct {
	Name  *Identifier
	Value *DecimalLit
}

func (it *EnumElement) NodeType() string {
	return "Enum_Element"
}

// 枚举语句节点
type EnumStatement struct {
	Name     *Identifier
	Elements []*EnumElement
}

func (it *EnumStatement) NodeType() string {
	return "Enum_Statement"
}
func (it *EnumStatement) StatementNodeType() int {
	return StatementTypeEnum
}

// 条件语句单元
type IfElement struct {
	Condition *Expression
	Block     *BlockStatement
}

func (it *IfElement) NodeType() string {
	return "If_Element"
}

// 条件语句节点
type IfStatement struct {
	If   *IfElement
	Elif []*IfElement
	Else *BlockStatement
}

func (it *IfStatement) NodeType() string {
	return "If_Statement"
}
func (it *IfStatement) StatementNodeType() int {
	return StatementTypeIf
}

// 分支语句单个条件接口
type SwitchStatementCase interface {
	Node
	SwitchStatementCaseNodeType() int
}

// 分支语句匹配条件单元
type SwitchStatementNormalCase struct {
	Conditions []*Expression
	Block      *BlockStatement
}

func (it *SwitchStatementNormalCase) NodeType() string {
	return "Switch_Statement_Normal_Case"
}
func (it *SwitchStatementNormalCase) SwitchStatementCaseNodeType() int {
	return SwitchStatementTypeNormal
}

// 分支语句匹配条件范围
type SwitchStatementRangeCase struct {
	From  *Expression
	To    *Expression
	Block *BlockStatement
}

func (it *SwitchStatementRangeCase) NodeType() string {
	return "Switch_Statement_Range_Case"
}
func (it *SwitchStatementRangeCase) SwitchStatementCaseNodeType() int {
	return SwitchStatementTypeRange
}

// 条件语句节点
type SwitchStatement struct {
	Keyword *Token
	Entry   *Expression
	Default *BlockStatement
	Cases   []*SwitchStatementCase
}

func (it *SwitchStatement) NodeType() string {
	return "Switch_Statement"
}
func (it *SwitchStatement) StatementNodeType() int {
	return StatementTypeSwitch
}

// while 语句
type WhileStatement struct {
	Keyword   *Token
	Condition *Expression
	Block     *BlockStatement
}

func (it *WhileStatement) NodeType() string {
	return "While_Statement"
}
func (it *WhileStatement) StatementNodeType() int {
	return StatementTypeWhile
}

// for 语句
type ForStatement struct {
	Keyword   *Token
	Initial   *Expression
	Condition *Expression
	tail      *Expression
	Block     *BlockStatement
}

func (it *ForStatement) NodeType() string {
	return "For_Statement"
}
func (it *ForStatement) StatementNodeType() int {
	return StatementTypeFor
}

// each 语句
type EachStatement struct {
	Keyword *Token
	Names   []*Identifier
	Target  *Expression
	Block   *BlockStatement
}

func (it *EachStatement) NodeType() string {
	return "Each_Statement"
}
func (it *EachStatement) StatementNodeType() int {
	return StatementTypeEach
}

// 函数形参节点
type Argument struct {
	Name *Identifier
	Type *TypeDescription
}

func (it *Argument) NodeType() string {
	return "Argument"
}

// 函数签名
type Signature struct {
	Arguments []*Argument
	Returns   []*TypeDescription
}

func (it *Signature) NodeType() string {
	return "Signature"
}

// 函数定义节点
type FunctionDeclarationStatement struct {
	Keyword   *Token
	Name      *Identifier
	Generics  *GenericArgs
	Signature *Signature
	Block     *BlockStatement
}

func (it *FunctionDeclarationStatement) NodeType() string {
	return "Function_Declaration_Statement"
}
func (it *FunctionDeclarationStatement) StatementNodeType() int {
	return StatementTypeFunctionDecl
}

// 类成员接口
type ClassMember interface {
	Node
	ClassMemberNodeType() int
}

// 类成员变量定义节点
type ClassMemberVar struct {
	Scope   *Token
	VarDecl *VarDeclStatement
}

func (it *ClassMemberVar) NodeType() string {
	return "Class_Member_Variable, scope: " + it.Scope.Str
}
func (it *ClassMemberVar) ClassMemberNodeType() int {
	return ClassMemberTypeVar
}

// 类成员方法定义节点
type ClassMemberMethod struct {
	Scope      *Token
	MethodDecl *FunctionDeclarationStatement
}

func (it *ClassMemberMethod) NodeType() string {
	return "Class_Member_Method, scope: " + it.Scope.Str
}
func (it *ClassMemberMethod) ClassMemberNodeType() int {
	return ClassMemberTypeMethod
}

// 泛型参数列表
type GenericArgs struct {
	Args []*Identifier
}

func (it *GenericArgs) NodeType() string {
	var args string
	for i, id := range it.Args {
		args += id.Name.Str
		if i != len(it.Args)-1 {
			args += ","
		}
	}
	return "Generics_Arguments: <" + args + ">"
}

// 类定义语句节点
type ClassDeclarationStatement struct {
	Name     *Identifier
	Generics *GenericArgs
	Extends  *Identifier
	Members  []*ClassMember
}

func (it *ClassDeclarationStatement) NodeType() string {
	return "Class_Declaration_Statement"
}
func (it *ClassDeclarationStatement) StatementNodeType() int {
	return StatementTypeClassDecl
}

// catch 错误捕获单元节点
type ErrorCatchUnit struct {
	Name    *Identifier
	Handler *BlockStatement
}

func (it *ErrorCatchUnit) NodeType() string {
	return "Error_Catch_Unit"
}

// try/catch 异常捕获语句节点
type TryCatchStatement struct {
	TryBlock *BlockStatement
	Handlers []*ErrorCatchUnit
	Finally  *BlockStatement
}

func (it *TryCatchStatement) NodeType() string {
	return "Try_Catch_Statement"
}
func (it *TryCatchStatement) StatementNodeType() int {
	return StatementTypeTryCatch
}

// Statement 为所有语句节点定义了接口
type Statement interface {
	Node
	StatementNodeType() int
}
