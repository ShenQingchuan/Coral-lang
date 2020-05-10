package ast

// 类型标注节点
type TypeDescription interface {
	Node
	TypeDescriptionNode() int
}

// 类型的名称节点
type TypeName struct {
	NameList []*Identifier
}

func (it *TypeName) GetFullName() string {
	var typeName string
	for i, id := range it.NameList {
		typeName += id.Name.Str
		if i != len(it.NameList)-1 {
			typeName += "."
		}
	}
	return typeName
}
func (it *TypeName) NodeType() string {
	return "Type_Name: " + it.GetFullName()
}
func (it *TypeName) TypeDescriptionNode() int {
	return TypeDescriptionTypeName
}

// 类型字面量节点
type TypeLit interface {
	Node
	TypeLitNode()
}

// 数组类型标识 eg: [T]
type ArrayTypeLit struct {
	ElementType *TypeDescription
}

func (it *ArrayTypeLit) NodeType() string {
	return "Array_Type_Lit"
}
func (it *ArrayTypeLit) TypeLitNode() {}
func (it *ArrayTypeLit) TypeDescriptionNode() int {
	return TypeDescriptionTypeLit
}

// 带泛型参数的标识 eg: A<B,C>
type GenericsTypeLit struct {
	BasicType    *TypeName
	GenericsArgs []*TypeName
}

func (it *GenericsTypeLit) NodeType() string {
	return "Generics_Type_Lit"
}
func (it *GenericsTypeLit) TypeLitNode() {}
func (it *GenericsTypeLit) TypeDescriptionNode() int {
	return TypeDescriptionTypeLit
}
