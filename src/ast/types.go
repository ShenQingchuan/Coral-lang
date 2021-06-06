package ast

// 类型标注节点
type TypeDescription interface {
	Node
	TypeDescriptionNode() int
}

// 类型的名称节点
type TypeName struct {
	Identifier *Identifier
}

func (it *TypeName) NodeType() string {
	return "Type_Name"
}
func (it *TypeName) TypeDescriptionNode() int {
	return TypeDescriptionTypeName
}

type FuncType struct {
	ArgTypes    []TypeDescription
	ReturnTypes []TypeDescription
}

func (it *FuncType) NodeType() string {
	return "Func_Type"
}
func (it *FuncType) TypeDescriptionNode() int {
	return TypeDescriptionFunction
}

// 数组类型标识 eg: T[]
type ArrayTypeLit struct {
	ElementType TypeDescription
	ArrayLength int
}

func (it *ArrayTypeLit) NodeType() string {
	return "Array_Type_Lit"
}
func (it *ArrayTypeLit) TypeDescriptionNode() int {
	return TypeDescriptionTypeArrayLit
}

// 带泛型参数的标识 eg: A<B,C>
type GenericsTypeLit struct {
	BasicType    *TypeName
	GenericsArgs []TypeDescription
}

func (it *GenericsTypeLit) NodeType() string {
	return "Generics_Type_Lit"
}
func (it *GenericsTypeLit) TypeDescriptionNode() int {
	return TypeDescriptionTypeGenerics
}
