package ast

const (
	// 定义类型标注的种类来区分
	TypeDescriptionTypeName = iota
	TypeDescriptionFunction
	TypeDescriptionTypeArrayLit
	TypeDescriptionTypeGenerics

	// 定义所有语句的种类来区分
	StatementTypeSimple
	StatementTypePackage
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
	StatementTypeInterfaceDecl

	// 定义引入外部模块语句的种类来区分
	ImportStatementTypeSingleGlobal
	ImportStatementTypeSingleFrom
	ImportStatementTypeList

	// 定义简单语句的种类来区分
	SimpleStmtTypeExpression
	SimpleStmtIncDecStmt
	SimpleStmtTypeVariableDecl
	SimpleStmtTypeAssignList

	// 定义条件语句匹配项的种类来区分
	SwitchStatementTypeNormal
	SwitchStatementTypeRange

	// 定义类成员的种类来区分
	ClassMemberTypeVar
	ClassMemberTypeMethod

	// 类成员的公开与否 枚举：
	ClassMemberScopePrivate
	ClassMemberScopePublic

	// 定义表达式的类型来区分
	ExpressionTypePrimary
	ExpressionTypeNewInstance
	ExpressionTypeUnary
	ExpressionTypeBinary
	ExpressionTypeRange

	// 定义基本表达式的类型来区分
	PrimaryExprTypeBasic
	PrimaryExprTypeIndex
	PrimaryExprTypeSlice
	PrimaryExprTypeCall
	PrimaryExprTypeMember

	// 定义操作数的类型来区分
	OperandTypeLiteral
	OperandTypeName

	// 定义字面量的类型来区分
	LiteralNodeTypeNil
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
	LiteralNodeTypeThis
	LiteralNodeTypeSuper
)
