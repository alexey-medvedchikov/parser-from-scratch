package tokenizer

import "regexp"

type Rule struct {
	Type   TokenType
	Regexp *regexp.Regexp
}

var DefaultRules = []Rule{
	{Type: Skip, Regexp: regexp.MustCompile(`^\s+`)},
	{Type: Skip, Regexp: regexp.MustCompile(`^//.*`)},
	{Type: Skip, Regexp: regexp.MustCompile(`^/\*[\s\S]*?\*/`)},
	{Type: Semicolon, Regexp: regexp.MustCompile(`^;`)},
	{Type: OpenCurlyBrace, Regexp: regexp.MustCompile(`^{`)},
	{Type: CloseCurlyBrace, Regexp: regexp.MustCompile(`^}`)},
	{Type: OpenParens, Regexp: regexp.MustCompile(`^\(`)},
	{Type: CloseParens, Regexp: regexp.MustCompile(`^\)`)},
	{Type: Comma, Regexp: regexp.MustCompile(`^,`)},
	{Type: Dot, Regexp: regexp.MustCompile(`^\.`)},
	{Type: OpenSquare, Regexp: regexp.MustCompile(`^\[`)},
	{Type: CloseSquare, Regexp: regexp.MustCompile(`^]`)},
	{Type: LetKeyword, Regexp: regexp.MustCompile(`^\blet\b`)},
	{Type: DefKeyword, Regexp: regexp.MustCompile(`^\bdef\b`)},
	{Type: ReturnKeyword, Regexp: regexp.MustCompile(`^\breturn\b`)},
	{Type: IfKeyword, Regexp: regexp.MustCompile(`^\bif\b`)},
	{Type: WhileKeyword, Regexp: regexp.MustCompile(`^\bwhile\b`)},
	{Type: DoKeyword, Regexp: regexp.MustCompile(`^\bdo\b`)},
	{Type: ClassKeyword, Regexp: regexp.MustCompile(`^\bclass\b`)},
	{Type: ThisKeyword, Regexp: regexp.MustCompile(`^\bthis\b`)},
	{Type: ExtendsKeyword, Regexp: regexp.MustCompile(`^\bextends\b`)},
	{Type: SuperKeyword, Regexp: regexp.MustCompile(`^\bsuper\b`)},
	{Type: NewKeyword, Regexp: regexp.MustCompile(`^\bnew\b`)},
	{Type: ForKeyword, Regexp: regexp.MustCompile(`^\bfor\b`)},
	{Type: ElseKeyword, Regexp: regexp.MustCompile(`^\belse\b`)},
	{Type: TrueKeyword, Regexp: regexp.MustCompile(`^\btrue\b`)},
	{Type: FalseKeyword, Regexp: regexp.MustCompile(`^\bfalse\b`)},
	{Type: NullKeyword, Regexp: regexp.MustCompile(`^\bnull\b`)},
	{Type: Number, Regexp: regexp.MustCompile(`^\d+`)},
	{Type: String, Regexp: regexp.MustCompile(`^"[^"]*"`)},
	{Type: String, Regexp: regexp.MustCompile(`^'[^"]*'`)},
	{Type: Identifier, Regexp: regexp.MustCompile(`^\w+`)},
	{Type: EqualityOp, Regexp: regexp.MustCompile(`^[=!]=`)},
	{Type: SimpleAssign, Regexp: regexp.MustCompile(`^=`)},
	{Type: ComplexAssign, Regexp: regexp.MustCompile(`^[+\-*/]=`)},
	{Type: NotLogicalOp, Regexp: regexp.MustCompile(`^!`)},
	{Type: AndLogicalOp, Regexp: regexp.MustCompile(`^&&`)},
	{Type: OrLogicalOp, Regexp: regexp.MustCompile(`^\|\|`)},
	{Type: RelationalOp, Regexp: regexp.MustCompile(`^[<>]=?`)},
	{Type: AdditiveOp, Regexp: regexp.MustCompile(`^[+\-]`)},
	{Type: MultiplicativeOp, Regexp: regexp.MustCompile(`^[*/]`)},
}

type Tokenizer struct {
	expr   string
	cursor int
	rules  []Rule
}

func NewTokenizer(rules []Rule, expr string) *Tokenizer {
	return &Tokenizer{
		expr:   expr,
		cursor: 0,
		rules:  rules,
	}
}

func (t *Tokenizer) NextToken() (*Token, error) {
	if t.cursor >= len(t.expr) {
		return &Token{
			Type: EOF,
		}, nil
	}

	for _, spec := range t.rules {
		rest := t.expr[t.cursor:]

		if matched, ok := t.match(spec.Regexp, rest); ok {
			if spec.Type == Skip {
				return t.NextToken()
			}

			return &Token{
				Type:  spec.Type,
				Value: matched,
			}, nil
		}
	}

	return nil, &ErrUnexpectedToken{
		Position:   t.cursor,
		CodeString: t.expr[t.cursor:],
	}
}

func (t *Tokenizer) match(re *regexp.Regexp, s string) (string, bool) {
	if m := re.FindStringIndex(s); m != nil {
		t.cursor += m[1] - m[0]
		return s[m[0]:m[1]], true
	}

	return "", false
}
