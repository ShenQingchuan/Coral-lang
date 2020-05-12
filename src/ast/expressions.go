package ast

import (
	. "coral-lang/src/lexer"
)

// 定义表达式的类型来区分
const (
	ExpressionTypePrimary = iota
	ExpressionTypeNewInstance
	ExpressionTypeUnary
	ExpressionTypeBinary
)

// 定义基本表达式的类型来区分
const (
	PrimaryExprTypeBasic = iota
	PrimaryExprTypeIndex
	PrimaryExprTypeSlice
	PrimaryExprTypeCall
)

// 定义操作数的类型来区分
const (
	OperandTypeLiteral = iota
	OperandTypeName
)

// 定义字面量的类型来区分
const (
	LiteralNodeTypeNil = iota
	LiteralNodeTypeTrue
	LiteralNodeTypeFalse
	LiteralNodeTypeDecimal
	LiteralNodeTypeHexadecimal
	LiteralNodeTypeOctal
	LiteralNodeTypeBinary
	LiteralNodeTypeFloat
	LiteralNodeTypeExponent
	LiteralNodeTypeChar
	LiteralNodeTypeString
	LiteralNodeTypeArray
	LiteralNodeTypeMap
	LiteralNodeTypeLambda
)

// ----- 各种字面量 ------

// nil
type NilLit struct {
	Value *Token
}

func (it *NilLit) NodeType() string {
	return "Nil_Lit"
}
func (it *NilLit) LiteralNodeType() int {
	return LiteralNodeTypeNil
}
func (it *NilLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// true
type TrueLit struct {
	Value *Token
}

func (it *TrueLit) NodeType() string {
	return "True_Lit"
}
func (it *TrueLit) LiteralNodeType() int {
	return LiteralNodeTypeTrue
}
func (it *TrueLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// false
type FalseLit struct {
	Value *Token
}

func (it *FalseLit) NodeType() string {
	return "False_Lit"
}
func (it *FalseLit) LiteralNodeType() int {
	return LiteralNodeTypeFalse
}
func (it *FalseLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 十进制整数
type DecimalLit struct {
	Value *Token
}

func (it *DecimalLit) NodeType() string {
	return "Decimal_Lit, value: " + it.Value.Str
}
func (it *DecimalLit) LiteralNodeType() int {
	return LiteralNodeTypeDecimal
}
func (it *DecimalLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 十六进制整数
type HexadecimalLit struct {
	Value *Token
}

func (it *HexadecimalLit) NodeType() string {
	return "Hexadecimal_Lit, value: " + it.Value.Str
}
func (it *HexadecimalLit) LiteralNodeType() int {
	return LiteralNodeTypeHexadecimal
}
func (it *HexadecimalLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 八进制整数
type OctalLit struct {
	Value *Token
}

func (it *OctalLit) NodeType() string {
	return "Octal_Lit, value: " + it.Value.Str
}
func (it *OctalLit) LiteralNodeType() int {
	return LiteralNodeTypeOctal
}
func (it *OctalLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 二进制整数
type BinaryLit struct {
	Value *Token
}

func (it *BinaryLit) NodeType() string {
	return "Binary_Lit, value: " + it.Value.Str
}
func (it *BinaryLit) LiteralNodeType() int {
	return LiteralNodeTypeBinary
}
func (it *BinaryLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 浮点数
type FloatLit struct {
	Value *Token
}

func (it *FloatLit) NodeType() string {
	return "Float_Lit, value: " + it.Value.Str
}
func (it *FloatLit) LiteralNodeType() int {
	return LiteralNodeTypeFloat
}
func (it *FloatLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 科学记数法
type ExponentLit struct {
	Value *Token
}

func (it *ExponentLit) NodeType() string {
	return "Exponent_Lit, value: " + it.Value.Str
}
func (it *ExponentLit) LiteralNodeType() int {
	return LiteralNodeTypeExponent
}
func (it *ExponentLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 字符
type CharLit struct {
	Value *Token
}

func (it *CharLit) NodeType() string {
	return "Char_Lit, value: " + it.Value.Str
}
func (it *CharLit) LiteralNodeType() int {
	return LiteralNodeTypeChar
}
func (it *CharLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 字符串
type StringLit struct {
	Value *Token
}

func (it *StringLit) NodeType() string {
	return "String_Lit, value: " + it.Value.Str
}
func (it *StringLit) LiteralNodeType() int {
	return LiteralNodeTypeString
}
func (it *StringLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 数组
type ArrayLit struct {
	ValueList []*Expression
}

func (it *ArrayLit) NodeType() string {
	return "Array_Lit"
}
func (it *ArrayLit) LiteralNodeType() int {
	return LiteralNodeTypeArray
}
func (it *ArrayLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 字典元素
type MapElement struct {
	Key   *Expression
	Value *Expression
}

func (it *MapElement) NodeType() string {
	return "Map_Element"
}

// 字典
type MapLit struct {
	KeyValueList []*MapElement
}

func (it *MapLit) NodeType() string {
	return "Map_Lit"
}
func (it *MapLit) LiteralNodeType() int {
	return LiteralNodeTypeMap
}
func (it *MapLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 箭头函数
type LambdaLit struct {
	Literal
	Arguments []*Expression
	Block     *BlockStatement
}

func (it *LambdaLit) NodeType() string {
	return "Lambda_Lit"
}
func (it *LambdaLit) LiteralNodeType() int {
	return LiteralNodeTypeLambda
}
func (it *LambdaLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// ----- 各种字面量 END ------

// 字面量节点
type Literal interface {
	Operand
	LiteralNodeType() int
}

// 操作数名节点
type OperandName struct {
	NameList []*Identifier
}

func (it *OperandName) NodeType() string {
	return "Operand_Name: " + it.GetFullName()
}
func (it *OperandName) GetFullName() string {
	var typeName string
	for i, id := range it.NameList {
		typeName += id.Name.Str
		if i != len(it.NameList)-1 {
			typeName += "."
		}
	}
	return typeName
}
func (it *OperandName) OperandNodeType() int {
	return OperandTypeName
}

// 结合语法定义可知 primaryExpr 有四种可能性
// 应当抽取出 operand，之后的三种情况可以继承
type Operand interface {
	OperandNodeType() int
}

// 只是操作数本身的 primaryExpr
type BasicPrimaryExpression struct {
	Operand Operand
}

func (it *BasicPrimaryExpression) NodeType() string {
	return "Basic_Primary_Expression"
}
func (it *BasicPrimaryExpression) PrimaryExpressionNode() int {
	return PrimaryExprTypeBasic
}
func (it *BasicPrimaryExpression) ExpressionNodeType() int {
	return ExpressionTypePrimary
}
func (it *BasicPrimaryExpression) SimpleStatementNodeType() int {
	return SimpleStmtTypeExpression
}
func (it *BasicPrimaryExpression) StatementNodeType() int {
	return StatementTypeSimple
}

// 索引访问表达式节点
type IndexExpression struct {
	Operand *BasicPrimaryExpression // 继承其操作数，其他三种下同
	Index   *Expression             // 索引表达式
}

func (it *IndexExpression) ExpressionNodeType() int {
	return ExpressionTypePrimary
}
func (it *IndexExpression) PrimaryExpressionNode() int {
	return PrimaryExprTypeIndex
}
func (it *IndexExpression) NodeType() string {
	return "Index_Expression"
}
func (it *IndexExpression) SimpleStatementNodeType() int {
	return SimpleStmtTypeExpression
}
func (it *IndexExpression) StatementNodeType() int {
	return StatementTypeSimple
}

// 切片访问表达式节点
type SliceExpression struct {
	Operand *BasicPrimaryExpression
	Start   *Expression // 切片位置起点
	End     *Expression // 切片位置终点
}

func (it *SliceExpression) NodeType() string {
	return "Slice_Expression"
}
func (it *SliceExpression) ExpressionNodeType() int {
	return ExpressionTypePrimary
}
func (it *SliceExpression) PrimaryExpressionNode() int {
	return PrimaryExprTypeSlice
}
func (it *SliceExpression) SimpleStatementNodeType() int {
	return SimpleStmtTypeExpression
}
func (it *SliceExpression) StatementNodeType() int {
	return StatementTypeSimple
}

// 函数调用表达式节点
type CallExpression struct {
	Operand *BasicPrimaryExpression
	Params  []*Expression // 函数实参列表
}

func (it *CallExpression) ExpressionNodeType() int {
	return ExpressionTypePrimary
}
func (it *CallExpression) PrimaryExpressionNode() int {
	return PrimaryExprTypeCall
}
func (it *CallExpression) NodeType() string {
	return "Call_Expression"
}
func (it *CallExpression) SimpleStatementNodeType() int {
	return SimpleStmtTypeExpression
}
func (it *CallExpression) StatementNodeType() int {
	return StatementTypeSimple
}

// 基本表达式节点
type PrimaryExpression interface {
	Expression
	PrimaryExpressionNode() int
}

// 新建对象实例表达式节点
type NewInstanceExpression struct {
	Class      *TypeDescription
	InitParams []*Expression
}

func (it *NewInstanceExpression) ExpressionNodeType() int {
	return ExpressionTypeNewInstance
}
func (it *NewInstanceExpression) NodeType() string {
	return "New_Instance_Expression"
}
func (it *NewInstanceExpression) SimpleStatementNodeType() int {
	return SimpleStmtTypeExpression
}
func (it *NewInstanceExpression) StatementNodeType() int {
	return StatementTypeSimple
}

// 一元表达式节点
type UnaryExpression struct {
	Operator *Token
	Operand  *Expression
}

func (it *UnaryExpression) ExpressionNodeType() int {
	return ExpressionTypeUnary
}
func (it *UnaryExpression) NodeType() string {
	return "Unary_Expression"
}
func (it *UnaryExpression) SimpleStatementNodeType() int {
	return SimpleStmtTypeExpression
}
func (it *UnaryExpression) StatementNodeType() int {
	return StatementTypeSimple
}

// 二元表达式节点
type BinaryExpression struct {
	Operator *Token
	left     *Expression
	right    *Expression
}

func (it *BinaryExpression) ExpressionNodeType() int {
	return ExpressionTypeBinary
}
func (it *BinaryExpression) NodeType() string {
	return "Binary_Expression"
}
func (it *BinaryExpression) SimpleStatementNodeType() int {
	return SimpleStmtTypeExpression
}
func (it *BinaryExpression) StatementNodeType() int {
	return StatementTypeSimple
}

// Expression 为所有表达式节点定义了接口
type Expression interface {
	SimpleStatement
	ExpressionNodeType() int
}
