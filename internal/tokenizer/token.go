package tokenizer

type Token struct {
	Type  TokenType
	Value string
}

type TokenType string

const (
	// EOF is a special type of token that indicates the end of the file
	EOF TokenType = "EOF"
	// Skip are tokens such as whitespace or comments
	Skip             TokenType = "Skip"
	Semicolon        TokenType = ";"
	OpenCurlyBrace   TokenType = "{"
	CloseCurlyBrace  TokenType = "}"
	OpenParens       TokenType = "("
	CloseParens      TokenType = ")"
	Comma            TokenType = ","
	Dot              TokenType = "."
	OpenSquare       TokenType = "["
	CloseSquare      TokenType = "]"
	LetKeyword       TokenType = "let"
	DefKeyword       TokenType = "def"
	ReturnKeyword    TokenType = "return"
	IfKeyword        TokenType = "if"
	WhileKeyword     TokenType = "while"
	DoKeyword        TokenType = "do"
	ClassKeyword     TokenType = "class"
	ThisKeyword      TokenType = "this"
	ExtendsKeyword   TokenType = "extends"
	SuperKeyword     TokenType = "super"
	NewKeyword       TokenType = "new"
	ForKeyword       TokenType = "for"
	ElseKeyword      TokenType = "else"
	TrueKeyword      TokenType = "true"
	FalseKeyword     TokenType = "false"
	NullKeyword      TokenType = "null"
	Number           TokenType = "Number"           // 10
	String           TokenType = "String"           // "hello"
	Identifier       TokenType = "Identifier"       // name of variable
	EqualityOp       TokenType = "EqualityOp"       // == !=
	SimpleAssign     TokenType = "="                // =
	ComplexAssign    TokenType = "ComplexAssign"    // *= /= += -=
	RelationalOp     TokenType = "RelationalOp"     // > < >= <=
	AndLogicalOp     TokenType = "AndLogicalOp"     // &&
	OrLogicalOp      TokenType = "OrLogicalOp"      // ||
	NotLogicalOp     TokenType = "NotLogicalOp"     // !
	AdditiveOp       TokenType = "AdditiveOp"       // + or -
	MultiplicativeOp TokenType = "MultiplicativeOp" // * or /
)

func (t TokenType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t + `"`), nil
}
