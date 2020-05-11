package ast

import (
	. "coral-lang/src/lexer"
)

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
	LiteralTypeNil = iota
	LiteralTypeDecimal
	LiteralTypeHexadecimal
	LiteralTypeOctal
	LiteralTypeBinary
	LiteralTypeFloat
	LiteralTypeExponent
	LiteralTypeChar
	LiteralTypeString
	LiteralTypeArray
	LiteralTypeMap
	LiteralTypeLambda
)

// ----- 各种字面量 ------

// nil
type NilLit struct {
	Literal
	Value *Token
}

func (it *NilLit) NodeType() string {
	return "Nil_Lit"
}
func (it *NilLit) LiteralNodeType() int {
	return LiteralTypeNil
}
func (it *NilLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 十进制整数
type DecimalLit struct {
	Literal
	Value *Token
}

func (it *DecimalLit) NodeType() string {
	return "Decimal_Lit, value: " + it.Value.Str
}
func (it *DecimalLit) LiteralNodeType() int {
	return LiteralTypeDecimal
}
func (it *DecimalLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 十六进制整数
type HexadecimalLit struct {
	Literal
	Value *Token
}

func (it *HexadecimalLit) NodeType() string {
	return "Hexadecimal_Lit, value: " + it.Value.Str
}
func (it *HexadecimalLit) LiteralType() int {
	return LiteralTypeHexadecimal
}
func (it *HexadecimalLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 八进制整数
type OctalLit struct {
	Literal
	Value *Token
}

func (it *OctalLit) NodeType() string {
	return "Octal_Lit, value: " + it.Value.Str
}
func (it *OctalLit) LiteralType() int {
	return LiteralTypeOctal
}
func (it *OctalLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 二进制整数
type BinaryLit struct {
	Literal
	Value *Token
}

func (it *BinaryLit) NodeType() string {
	return "Binary_Lit, value: " + it.Value.Str
}
func (it *BinaryLit) LiteralType() int {
	return LiteralTypeBinary
}
func (it *BinaryLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 浮点数
type FloatLit struct {
	Literal
	Value *Token
}

func (it *FloatLit) NodeType() string {
	return "Float_Lit, value: " + it.Value.Str
}
func (it *FloatLit) LiteralType() int {
	return LiteralTypeFloat
}
func (it *FloatLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 科学记数法
type ExponentLit struct {
	Literal
	Value *Token
}

func (it *ExponentLit) NodeType() string {
	return "Exponent_Lit, value: " + it.Value.Str
}
func (it *ExponentLit) LiteralType() int {
	return LiteralTypeExponent
}
func (it *ExponentLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 字符
type CharLit struct {
	Literal
	Value *Token
}

func (it *CharLit) NodeType() string {
	return "Char_Lit, value: " + it.Value.Str
}
func (it *CharLit) LiteralType() int {
	return LiteralTypeChar
}
func (it *CharLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 字符串
type StringLit struct {
	Literal
	Value *Token
}

func (it *StringLit) NodeType() string {
	return "String_Lit, value: " + it.Value.Str
}
func (it *StringLit) LiteralType() int {
	return LiteralTypeString
}
func (it *StringLit) OperandNodeType() int {
	return OperandTypeLiteral
}

// 数组
type ArrayLit struct {
	Literal
	ValueList []*Expression
}

func (it *ArrayLit) NodeType() string {
	return "Array_Lit"
}
func (it *ArrayLit) LiteralType() int {
	return LiteralTypeArray
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
	Literal
	KeyValueList []*MapElement
}

func (it *MapLit) NodeType() string {
	return "Map_Lit"
}
func (it *MapLit) LiteralType() int {
	return LiteralTypeMap
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
func (it *LambdaLit) LiteralType() int {
	return LiteralTypeLambda
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
	Node
	OperandNodeType() int
}

// 只是操作数本身的 primaryExpr
type BasicPrimaryExpression struct {
	Operand *Operand
}

func (it *BasicPrimaryExpression) NodeType() string {
	return "Basic_Primary_Expression"
}
func (it *BasicPrimaryExpression) PrimaryExpressionNode() {}
func (it *BasicPrimaryExpression) ExpressionNodeType() int {
	return ExpressionTypePrimary
}

// 索引访问表达式节点
type IndexExpression struct {
	*BasicPrimaryExpression             // 继承其操作数，其他三种下同
	index                   *Expression // 索引表达式
}

func (it *IndexExpression) ExpressionNodeType() int {
	return ExpressionTypePrimary
}

// 切片访问表达式节点
type SliceExpression struct {
	*BasicPrimaryExpression
	start *Expression // 切片位置起点
	end   *Expression // 切片位置终点
}

func (it *SliceExpression) ExpressionNodeType() int {
	return ExpressionTypePrimary
}

// 函数调用表达式节点
type CallExpression struct {
	*BasicPrimaryExpression
	params []*Expression // 函数实参列表
}

func (it *CallExpression) ExpressionNodeType() int {
	return ExpressionTypePrimary
}

// 基本表达式节点
type PrimaryExpression interface {
	Node
	PrimaryExpressionNode()
}

// 新建对象实例表达式节点
type NewInstanceExpression struct {
	Class      *TypeDescription
	InitParams []*Expression
}

func (it *NewInstanceExpression) ExpressionNodeType() int {
	return ExpressionTypeNewInstance
}

// 一元表达式节点
type UnaryExpression struct {
	Operator *Token
	Operand  *Expression
}

func (it *UnaryExpression) ExpressionNodeType() int {
	return ExpressionTypeUnary
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

// Expression 为所有表达式节点定义了接口
type Expression interface {
	Node
	ExpressionNodeType() int
}
