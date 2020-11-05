package analyzer

const (
	KeywordSymbol = iota
	BuiltinsTypeSymbol
	VariableIdentifierSymbol
	VariableClassTypeSymbol
	FunctionNameSymbol
)

type BlockSymbolTableItem struct {
	Name       string
	SymbolType int
}
type BlockSymbolTable struct {
	OuterLevel *BlockSymbolTable
	SymbolList []*BlockSymbolTableItem
}

type Analyzer struct {
	SymbolTableRoot *BlockSymbolTable
}

func loadKeywordsSymbol(table *BlockSymbolTable) *BlockSymbolTable {
	keywordsList := []string{"import", "package", "from", "as", "enum", "break",
		"continue", "return", "var", "val", "if", "elif", "else", "switch", "default",
		"case", "while", "for", "each", "in", "fn", "class", "interface", "this", "super",
		"static", "public", "private", "new", "nil", "true", "false", "try", "catch", "finally", "throws",
	}
	for _, keyword := range keywordsList {
		table.SymbolList = append(table.SymbolList, &BlockSymbolTableItem{
			Name:       keyword,
			SymbolType: KeywordSymbol,
		})
	}

	return table
}
func loadBuiltinsTypeSymbol(table *BlockSymbolTable) *BlockSymbolTable {
	builtinsTypes := []string{
		"int8", "int16", "int", "int64",
		"float", "float64",
		"bool",
		"rune",
	}
	for _, keyword := range builtinsTypes {
		table.SymbolList = append(table.SymbolList, &BlockSymbolTableItem{
			Name:       keyword,
			SymbolType: BuiltinsTypeSymbol,
		})
	}

	return table
}
func (analyzer *Analyzer) initSymbolTable() {
	loadKeywordsSymbol(loadBuiltinsTypeSymbol(analyzer.SymbolTableRoot))
}
func (analyzer *Analyzer) initialize() {
	analyzer.initSymbolTable()
}
