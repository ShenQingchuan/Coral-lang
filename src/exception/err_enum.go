package exception

const (
	NormalError = iota
	FileSystemOpenFileError
	CompilerRegexExpCreatingFailed
	LexFloatFormatError
	LexExponentFormatError
	LexUnicodeEscapeFormatError
	LexIdentifierFirstRuneCanNotBeDigit
	LexBlockCommentTooNested
)
