package analyzer

import (
	. "coral-lang/src/ast"
)

func (analyzer *Analyzer) CheckStatement(stmt Statement) {
	switch stmt.StatementNodeType() {
	case StatementTypeSimple:
		simpleStmt := stmt.(SimpleStatement)
		analyzer.CheckSimpleStatement(simpleStmt)
	case StatementTypeEnum:
		enumStmt := stmt.(*EnumStatement)
		analyzer.CheckEnumStatement(enumStmt)
	case StatementTypeBlock:
		analyzer.EnterNewBlockScope()
		blockStmt := stmt.(*BlockStatement)
		analyzer.CheckBlockStatement(blockStmt)
		analyzer.LeaveCurrentBlockScope()
	case StatementTypeIf:
	case StatementTypeSwitch:
	case StatementTypeWhile:
	case StatementTypeFor:
	case StatementTypeEach:
	case StatementTypeFunctionDecl:
	case StatementTypeClassDecl:
	case StatementTypeInterfaceDecl:
	case StatementTypeTryCatch:
	}
}

func (analyzer *Analyzer) CheckSimpleStatement(simpleStmt SimpleStatement) {
	switch simpleStmt.SimpleStatementNodeType() {
	case SimpleStmtTypeExpression:
	case SimpleStmtTypeVariableDecl:
	case SimpleStmtTypeAssignList:
	case SimpleStmtTypeIncDecStmt:
	}
}

func (analyzer *Analyzer) CheckEnumStatement(enumStmt *EnumStatement) {
	enumSymbol := new(EnumSymbol)
	enumSymbol.Token = enumStmt.Name.Token
	enumSymbol.CollectionName = enumStmt.Name.GetName()
	for _, enumElement := range enumStmt.Elements {
		enumSymbol.ElementsMap[enumElement.Name.GetName()] = enumElement
	}
	analyzer.CurrentScope.SymbolMap[enumSymbol.CollectionName] = enumSymbol
}

func (analyzer *Analyzer) CheckBlockStatement(blockStmt *BlockStatement) {
	for _, stmt := range blockStmt.Statements {
		analyzer.CheckStatement(stmt)
	}
}
