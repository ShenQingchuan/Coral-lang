package analyzer

type NameTableKind int

const ( // 名字表
	NameKindConstant = iota
	NameKindVariable
	NameKindBuiltInType
	NameKindFunction
	NameKindClass
	NameKindInterface
	NameKindEnum
)

type NameTableItem struct {
	Name  string
	Kind  NameTableKind
	Level int
	Link  *NameTableItem
}
type BlockTableItem struct {
	First *NameTableItem
	Last  *NameTableItem
}

var StandardLibraries = []string{
	"stdlib",
	"math",
	"json",
	"regexp",
	"httplib",
	"orm",
}

func isStandardLibrary(packageName string) bool {
	for _, libName := range StandardLibraries {
		if libName == packageName {
			return true
		}
	}
	return false
}
