package parser

import (
	"strconv"

	"github.com/alexey-medvedchikov/parser-from-scratch/internal/ast"
	"github.com/alexey-medvedchikov/parser-from-scratch/internal/tokenizer"
)

type Tokenizer interface {
	NextToken() (*tokenizer.Token, error)
}

type Parser struct {
	tokenizer Tokenizer
	lookahead *tokenizer.Token
	builder   ast.Builder
}

func NewParser(t Tokenizer, b ast.Builder) *Parser {
	return &Parser{
		tokenizer: t,
		builder:   b,
	}
}

func (p *Parser) Parse() (ast.Node, error) {
	var err error
	p.lookahead, err = p.tokenizer.NextToken()
	if err != nil {
		return nil, err
	}

	return p.program()
}

// Program
//   : StatementList
//   ;
func (p *Parser) program() (ast.Node, error) {
	body, err := p.stmtList(tokenizer.EOF)
	if err != nil {
		return nil, err
	}

	return p.builder.Program(body...), nil
}

// StmtList
//   : Stmt
//   | StmtList Stmt
//   ;
func (p *Parser) stmtList(stopLookahead tokenizer.TokenType) ([]ast.Node, error) {
	statement, err := p.stmt()
	if err != nil {
		return nil, err
	}
	statementList := []ast.Node{statement}

	for p.lookahead != nil && p.lookahead.Type != stopLookahead {
		statement, err := p.stmt()
		if err != nil {
			return nil, err
		}
		statementList = append(statementList, statement)
	}

	return statementList, nil
}

// Stmt
//   : ExprStmt
//   | BlockStmt
//   | EmptyStmt
//   | VarStmt
//   | IfStmt
//   | IterStmt
//   | FuncDecl
//   | ReturnStmt
//   | ClassDecl
//   ;
func (p *Parser) stmt() (ast.Node, error) {
	switch p.lookahead.Type {
	case tokenizer.Semicolon:
		return p.emptyStmt()
	case tokenizer.OpenCurlyBrace:
		return p.blockStmt()
	case tokenizer.LetKeyword:
		return p.varStmt()
	case tokenizer.IfKeyword:
		return p.ifStmt()
	case tokenizer.WhileKeyword, tokenizer.DoKeyword, tokenizer.ForKeyword:
		return p.iterStmt()
	case tokenizer.DefKeyword:
		return p.funcDecl()
	case tokenizer.ClassKeyword:
		return p.classDecl()
	case tokenizer.ReturnKeyword:
		return p.returnStmt()
	default:
		return p.exprStmt()
	}
}

// ExprStmt
//   : SeqExpr ';'
//   ;
func (p *Parser) exprStmt() (ast.Node, error) {
	node, err := p.seqExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.Semicolon); err != nil {
		return nil, err
	}

	return p.builder.ExprStmt(node), nil
}

// BlockStmt
//   : '{' OptStmtList '}'
//   ;
func (p *Parser) blockStmt() (ast.Node, error) {
	if _, err := p.consume(tokenizer.OpenCurlyBrace); err != nil {
		return nil, err
	}

	var body []ast.Node
	if p.lookahead.Type != tokenizer.CloseCurlyBrace {
		var err error
		body, err = p.stmtList(tokenizer.CloseCurlyBrace)
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(tokenizer.CloseCurlyBrace); err != nil {
		return nil, err
	}

	return p.builder.BlockStmt(body...), nil
}

// EmptyStmt
//   : ';'
//   ;
func (p *Parser) emptyStmt() (ast.Node, error) {
	if _, err := p.consume(tokenizer.Semicolon); err != nil {
		return nil, err
	}

	return p.builder.EmptyStmt(), nil
}

// VarStmt
//   : VarStmtInit ';'
//   ;
func (p *Parser) varStmt() (ast.Node, error) {
	node, err := p.varStmtInit()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.Semicolon); err != nil {
		return nil, err
	}

	return node, nil
}

// VarStmtInit
//   : 'let' VarDeclList
//   ;
func (p *Parser) varStmtInit() (ast.Node, error) {
	if _, err := p.consume(tokenizer.LetKeyword); err != nil {
		return nil, err
	}

	declarations, err := p.varDeclList()
	if err != nil {
		return nil, err
	}

	return p.builder.VarStmt(declarations...), nil
}

// IfStmt
//   : 'if' '(' SeqExpr ')' Stmt
//   | 'if' '(' SeqExpr ')' Stmt 'else' Stmt
//   ;
func (p *Parser) ifStmt() (ast.Node, error) {
	if _, err := p.consume(tokenizer.IfKeyword); err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.OpenParens); err != nil {
		return nil, err
	}

	cond, err := p.seqExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.CloseParens); err != nil {
		return nil, err
	}

	cons, err := p.stmt()
	if err != nil {
		return nil, err
	}

	var alt ast.Node
	if p.lookahead.Type == tokenizer.ElseKeyword {
		if _, err := p.consume(tokenizer.ElseKeyword); err != nil {
			return nil, err
		}
		alt, err = p.stmt()
		if err != nil {
			return nil, err
		}
	}

	return p.builder.IfStmt(cond, cons, alt), nil
}

// IterStmt
//   : WhileStmt
//   | DoWhileStmt
//   | ForStmt
//   ;
func (p *Parser) iterStmt() (ast.Node, error) {
	switch p.lookahead.Type {
	case tokenizer.WhileKeyword:
		return p.whileStmt()
	case tokenizer.DoKeyword:
		return p.doWhileStmt()
	case tokenizer.ForKeyword:
		return p.forStmt()
	default:
		return nil, &ErrUnexpectedToken{
			Type:         p.lookahead.Type,
			ExpectedType: "Iteration",
			Value:        p.lookahead.Value,
		}
	}
}

// FuncDecl
//   : 'def' Identifier '(' OptFormalParamList ')' BlockStmt
//   ;
func (p *Parser) funcDecl() (ast.Node, error) {
	if _, err := p.consume(tokenizer.DefKeyword); err != nil {
		return nil, err
	}

	name, err := p.identifier()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.OpenParens); err != nil {
		return nil, err
	}

	var params []ast.Node
	if p.lookahead.Type != tokenizer.CloseParens {
		var err error
		if params, err = p.formalParamList(); err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(tokenizer.CloseParens); err != nil {
		return nil, err
	}

	body, err := p.blockStmt()
	if err != nil {
		return nil, err
	}

	return p.builder.FuncDecl(name, params, body), nil
}

// FormalParamList
//   : Identifier
//   | FormalParamList ',' Identifier
//   ;
func (p *Parser) formalParamList() ([]ast.Node, error) {
	var params []ast.Node

	for {
		param, err := p.identifier()
		if err != nil {
			return nil, err
		}
		params = append(params, param)
		if p.lookahead.Type != tokenizer.Comma {
			break
		}
		if _, err := p.consume(tokenizer.Comma); err != nil {
			return nil, err
		}
	}

	return params, nil
}

// ReturnStmt
//   : 'return' OptSeqExpr
//   ;
func (p *Parser) returnStmt() (ast.Node, error) {
	if _, err := p.consume(tokenizer.ReturnKeyword); err != nil {
		return nil, err
	}

	var arg ast.Node
	if p.lookahead.Type != tokenizer.Semicolon {
		var err error
		if arg, err = p.seqExpr(); err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(tokenizer.Semicolon); err != nil {
		return nil, err
	}

	return p.builder.ReturnStmt(arg), nil
}

// ClassDecl
//   : 'class' Identifier OptClassExtends BlockStmt
//   ;
func (p *Parser) classDecl() (ast.Node, error) {
	if _, err := p.consume(tokenizer.ClassKeyword); err != nil {
		return nil, err
	}

	id, err := p.identifier()
	if err != nil {
		return nil, err
	}

	var superClass ast.Node
	if p.lookahead.Type == tokenizer.ExtendsKeyword {
		if superClass, err = p.classExtends(); err != nil {
			return nil, err
		}
	}

	body, err := p.blockStmt()
	if err != nil {
		return nil, err
	}

	return p.builder.ClassDecl(id, superClass, body), nil
}

// ClassExtends
//   : 'extends' Identifier
func (p *Parser) classExtends() (ast.Node, error) {
	if _, err := p.consume(tokenizer.ExtendsKeyword); err != nil {
		return nil, err
	}

	return p.identifier()
}

// WhileStmt
//   : 'while' '(' SeqExpr ')' Stmt
//   ;
func (p *Parser) whileStmt() (ast.Node, error) {
	if _, err := p.consume(tokenizer.WhileKeyword); err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.OpenParens); err != nil {
		return nil, err
	}

	cond, err := p.seqExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.CloseParens); err != nil {
		return nil, err
	}

	body, err := p.stmt()
	if err != nil {
		return nil, err
	}

	return p.builder.WhileStmt(cond, body), nil
}

// DoWhileStmt
//   : 'do' Stmt 'while' '(' SeqExpr ')' ';'
//   ;
func (p *Parser) doWhileStmt() (ast.Node, error) {
	if _, err := p.consume(tokenizer.DoKeyword); err != nil {
		return nil, err
	}

	body, err := p.stmt()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.WhileKeyword); err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.OpenParens); err != nil {
		return nil, err
	}

	cond, err := p.seqExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.CloseParens); err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.Semicolon); err != nil {
		return nil, err
	}

	return p.builder.DoWhileStmt(cond, body), nil
}

// ForStmt
//   : 'for' '(' OptForStmtInit ';' OptSeqExpr ';' OptSeqExpr ')' Stmt
//   ;
func (p *Parser) forStmt() (ast.Node, error) {
	if _, err := p.consume(tokenizer.ForKeyword); err != nil {
		return nil, err
	}

	if _, err := p.consume(tokenizer.OpenParens); err != nil {
		return nil, err
	}

	var init ast.Node
	if p.lookahead.Type != tokenizer.Semicolon {
		var err error
		if init, err = p.forStmtInit(); err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(tokenizer.Semicolon); err != nil {
		return nil, err
	}

	var cond ast.Node
	if p.lookahead.Type != tokenizer.Semicolon {
		var err error
		if cond, err = p.seqExpr(); err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(tokenizer.Semicolon); err != nil {
		return nil, err
	}

	var step ast.Node
	if p.lookahead.Type != tokenizer.CloseParens {
		var err error
		if step, err = p.seqExpr(); err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(tokenizer.CloseParens); err != nil {
		return nil, err
	}

	body, err := p.stmt()
	if err != nil {
		return nil, err
	}

	return p.builder.ForStmt(init, cond, step, body), nil
}

// ForStmtInit
//   : VarStmtInit
//   | SeqExpr
//   ;
func (p *Parser) forStmtInit() (ast.Node, error) {
	if p.lookahead.Type == tokenizer.LetKeyword {
		return p.varStmtInit()
	}
	return p.seqExpr()
}

// VarDeclList
//   : VarDecl
//   | VarDeclList ',' VarDecl
//   ;
func (p *Parser) varDeclList() ([]ast.Node, error) {
	var declarations []ast.Node

	for {
		declaration, err := p.varDecl()
		if err != nil {
			return nil, err
		}
		declarations = append(declarations, declaration)
		if p.lookahead.Type != tokenizer.Comma {
			break
		}
		_, _ = p.consume(tokenizer.Comma)
	}

	return declarations, nil
}

// VarDecl
//   : Identifier OptVarInit
//   ;
func (p *Parser) varDecl() (ast.Node, error) {
	id, err := p.identifier()
	if err != nil {
		return nil, err
	}

	var init ast.Node
	if p.lookahead.Type != tokenizer.Comma && p.lookahead.Type != tokenizer.Semicolon {
		init, err = p.varInit()
		if err != nil {
			return nil, err
		}
	}

	return p.builder.VarDecl(id, init), nil
}

// VarInit
//   : SIMPLE_ASSIGN AssignExpr
//   ;
func (p *Parser) varInit() (ast.Node, error) {
	if _, err := p.consume(tokenizer.SimpleAssign); err != nil {
		return nil, err
	}

	return p.assignExpr()
}

// SeqExpr
//   : Expr
//   | SeqExpr ',' Expr
//   ;
func (p *Parser) seqExpr() (ast.Node, error) {
	var body []ast.Node

	for {
		expr, err := p.expr()
		if err != nil {
			return nil, err
		}

		body = append(body, expr)

		if p.lookahead.Type != tokenizer.Comma {
			break
		}
		_, _ = p.consume(tokenizer.Comma)
	}

	if len(body) == 1 {
		return body[0], nil
	}
	return p.builder.SeqExpr(body...), nil
}

// Expr
//   : AssignExpr
//   ;
func (p *Parser) expr() (ast.Node, error) {
	return p.assignExpr()
}

// AssignExpr
//   : EqualExpr
//   | LeftHandSideExpr AssignOp AssignExpr
//   ;
func (p *Parser) assignExpr() (ast.Node, error) {
	left, err := p.logicalOrExpr()
	if err != nil {
		return nil, err
	}

	if p.lookahead.Type != tokenizer.SimpleAssign &&
		p.lookahead.Type != tokenizer.ComplexAssign {
		return left, nil
	}

	opTok, err := p.assignOp()
	if err != nil {
		return nil, err
	}

	op := ast.AssignOpFromString(opTok.Value)
	if op == ast.InvalidAssignOp {
		return nil, &ErrUnknownAssignOp{
			Op: opTok.Value,
		}
	}

	if err := checkValidAssignTarget(left); err != nil {
		return nil, err
	}

	right, err := p.assignExpr()
	if err != nil {
		return nil, err
	}

	return p.builder.AssignExpr(op, left, right), nil
}

// AssignOp
//   : SIMPLE_ASSIGN
//   | COMPLEX_ASSIGN
//   ;
func (p *Parser) assignOp() (*tokenizer.Token, error) {
	if p.lookahead.Type == tokenizer.SimpleAssign {
		return p.consume(tokenizer.SimpleAssign)
	}
	return p.consume(tokenizer.ComplexAssign)
}

// Identifier
//   : IDENTIFIER
//   ;
func (p *Parser) identifier() (ast.Node, error) {
	tok, err := p.consume(tokenizer.Identifier)
	if err != nil {
		return nil, err
	}

	return p.builder.Identifier(tok.Value), nil
}

// ThisExpr
//   : 'this'
//   ;
func (p *Parser) thisExpr() (ast.Node, error) {
	if _, err := p.consume(tokenizer.ThisKeyword); err != nil {
		return nil, err
	}

	return p.builder.ThisExpr(), nil
}

// SuperCall
//   : 'super'
//   ;
func (p *Parser) superCall() (ast.Node, error) {
	if _, err := p.consume(tokenizer.SuperKeyword); err != nil {
		return nil, err
	}

	return p.builder.SuperCall(), nil
}

func checkValidAssignTarget(n ast.Node) error {
	switch n.Type {
	case ast.IdentifierType, ast.MemberExprType:
		return nil
	}

	return &ErrInvalidLvalue{Node: n}
}

// LogicalOrExpr
//   : LogicalAndExpr LOGICAL_OR LogicalOrExpr
//   | LogicalAndExpression
//   ;
func (p *Parser) logicalOrExpr() (ast.Node, error) {
	return p.logicalExpr(p.logicalAndExpr, tokenizer.OrLogicalOp)
}

// LogicalAndExpr
//   : EqualExpr LOGICAL_AND LogicalAndExpr
//   | EqualExpr
//   ;
func (p *Parser) logicalAndExpr() (ast.Node, error) {
	return p.logicalExpr(p.equalExpr, tokenizer.AndLogicalOp)
}

// EqualExpr
//   : RelExpr
//   | RelExpr EQUALITY_OP EqualExpr
func (p *Parser) equalExpr() (ast.Node, error) {
	return p.binaryExpr(p.relExpr, tokenizer.EqualityOp)
}

// RelExpr
//   : AddExpr
//   | RelExpr RELATIONAL_OP AddExpr
//   ;
func (p *Parser) relExpr() (ast.Node, error) {
	return p.binaryExpr(p.addExpr, tokenizer.RelationalOp)
}

// AddExpr
//   : MultExpr
//   | AddExpr ADDITIVE_OP MultExpr
//   ;
func (p *Parser) addExpr() (ast.Node, error) {
	return p.binaryExpr(p.multExpr, tokenizer.AdditiveOp)
}

// MultExpr
//   : UnaryExpr
//   | MultExpr ADDITIVE_OP UnaryExpr
//   ;
func (p *Parser) multExpr() (ast.Node, error) {
	return p.binaryExpr(p.unaryExpr, tokenizer.MultiplicativeOp)
}

func (p *Parser) binaryExpr(buildFunc func() (ast.Node, error), tokenType tokenizer.TokenType,
) (ast.Node, error) {
	left, err := buildFunc()
	if err != nil {
		return nil, err
	}

	for p.lookahead.Type == tokenType {
		opToken, err := p.consume(tokenType)
		if err != nil {
			return nil, err
		}

		op := ast.BinaryOpFromString(opToken.Value)
		if op == ast.InvalidBinaryOp {
			return nil, &ErrUnknownBinaryOp{Op: opToken.Value}
		}

		right, err := buildFunc()
		if err != nil {
			return nil, err
		}

		left = p.builder.BinaryExpr(op, left, right)
	}

	return left, nil
}

func (p *Parser) logicalExpr(buildFunc func() (ast.Node, error), tokenType tokenizer.TokenType,
) (ast.Node, error) {
	left, err := buildFunc()
	if err != nil {
		return nil, err
	}

	for p.lookahead.Type == tokenType {
		opToken, err := p.consume(tokenType)
		if err != nil {
			return nil, err
		}

		op := ast.LogicalOpFromString(opToken.Value)
		if op == ast.InvalidLogicalOp {
			return nil, &ErrUnknownLogicalOp{Op: opToken.Value}
		}

		right, err := buildFunc()
		if err != nil {
			return nil, err
		}

		left = p.builder.LogicalExpr(op, left, right)
	}

	return left, nil
}

// UnaryExpr
//   : LeftHandSideExpr
//   | ADDITIVE_OP UnaryExpr
//   | LOGICAL_NOT UnaryExpr
//   ;
func (p *Parser) unaryExpr() (ast.Node, error) {
	var opTok *tokenizer.Token
	var err error
	switch p.lookahead.Type {
	case tokenizer.AdditiveOp:
		if opTok, err = p.consume(tokenizer.AdditiveOp); err != nil {
			return nil, err
		}
	case tokenizer.NotLogicalOp:
		if opTok, err = p.consume(tokenizer.NotLogicalOp); err != nil {
			return nil, err
		}
	default:
		return p.leftHandSideExpr()
	}

	op := ast.UnaryOpFromString(opTok.Value)
	if op == ast.InvalidUnaryOp {
		return nil, &ErrUnknownUnaryOp{Op: opTok.Value}
	}

	arg, err := p.unaryExpr()
	if err != nil {
		return nil, err
	}

	return p.builder.UnaryExpr(op, arg), nil
}

// LeftHandSideExpr
//   : CallMemberExpr
//   ;
func (p *Parser) leftHandSideExpr() (ast.Node, error) {
	return p.callMemberExpr()
}

// CallMemberExpr
//   : MemberExpr
//   | CallExpr
//   | SuperCall CallExpr
//   ;
func (p *Parser) callMemberExpr() (ast.Node, error) {
	if p.lookahead.Type == tokenizer.SuperKeyword {
		super, err := p.superCall()
		if err != nil {
			return nil, err
		}
		return p.callExpr(super)
	}

	member, err := p.memberExpr()
	if err != nil {
		return nil, err
	}

	if p.lookahead.Type == tokenizer.OpenParens {
		return p.callExpr(member)
	}

	return member, nil
}

// CallExpr
//   : Callee CallArgs
//   ;
//
// Calee
//   : MemberExpr
//   | CallExpr
//   ;
func (p *Parser) callExpr(callee ast.Node) (ast.Node, error) {
	args, err := p.callArgs()
	if err != nil {
		return nil, err
	}

	callExpr := p.builder.CallExpr(callee, args)

	if p.lookahead.Type == tokenizer.OpenParens {
		callExpr, err = p.callExpr(callExpr)
		if err != nil {
			return nil, err
		}
	}

	return callExpr, nil
}

// NewExpr
//   : 'new' MemberExpression CallArgs
//   ;
func (p *Parser) newExpr() (ast.Node, error) {
	if _, err := p.consume(tokenizer.NewKeyword); err != nil {
		return nil, err
	}

	member, err := p.memberExpr()
	if err != nil {
		return nil, err
	}

	args, err := p.callArgs()
	if err != nil {
		return nil, err
	}

	return p.builder.NewExpr(member, args), nil
}

// CallArgs
//   : '(' OptArgList ')'
//   ;
func (p *Parser) callArgs() ([]ast.Node, error) {
	if _, err := p.consume(tokenizer.OpenParens); err != nil {
		return nil, err
	}

	var argList []ast.Node
	if p.lookahead.Type != tokenizer.CloseParens {
		var err error
		if argList, err = p.argList(); err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(tokenizer.CloseParens); err != nil {
		return nil, err
	}

	return argList, nil
}

// ArgList
//   : AssignExpr
//   | ArgList ',' AssignExpr
func (p *Parser) argList() ([]ast.Node, error) {
	var result []ast.Node

	for {
		arg, err := p.assignExpr()
		if err != nil {
			return nil, err
		}

		result = append(result, arg)

		if p.lookahead.Type != tokenizer.Comma {
			break
		}

		if _, err := p.consume(tokenizer.Comma); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// MemberExpr
//   : PrimaryExpr
//   | MemberExpr '.' Identifier
//   | MemberExpr '[' SeqExpr ']'
//   ;
func (p *Parser) memberExpr() (ast.Node, error) {
	obj, err := p.primaryExpr()
	if err != nil {
		return nil, err
	}

	for {
		if p.lookahead.Type == tokenizer.Dot {
			if _, err := p.consume(tokenizer.Dot); err != nil {
				return nil, err
			}
			prop, err := p.identifier()
			if err != nil {
				return nil, err
			}
			obj = p.builder.MemberExpr(false, obj, prop)
		} else if p.lookahead.Type == tokenizer.OpenSquare {
			if _, err := p.consume(tokenizer.OpenSquare); err != nil {
				return nil, err
			}
			prop, err := p.seqExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.consume(tokenizer.CloseSquare); err != nil {
				return nil, err
			}
			obj = p.builder.MemberExpr(true, obj, prop)
		} else {
			break
		}
	}

	return obj, nil
}

// PrimaryExpr
//   : Literal
//   | ParensExpr
//   | Identifier
//   | ThisExpr
//   | NewExpr
//   ;
func (p *Parser) primaryExpr() (ast.Node, error) {
	if isLiteral(p.lookahead.Type) {
		return p.literal()
	}
	switch p.lookahead.Type {
	case tokenizer.OpenParens:
		return p.parensExpr()
	case tokenizer.Identifier:
		return p.identifier()
	case tokenizer.ThisKeyword:
		return p.thisExpr()
	case tokenizer.NewKeyword:
		return p.newExpr()
	default:
		return p.leftHandSideExpr()
	}
}

func isLiteral(t tokenizer.TokenType) bool {
	switch t {
	case tokenizer.String,
		tokenizer.Number,
		tokenizer.TrueKeyword,
		tokenizer.FalseKeyword,
		tokenizer.NullKeyword:
		return true
	default:
		return false
	}
}

// ParensExpr
//   : '(' SeqExpr ')'
//   ;
func (p *Parser) parensExpr() (ast.Node, error) {
	if _, err := p.consume(tokenizer.OpenParens); err != nil {
		return nil, err
	}

	expr, err := p.seqExpr()
	if err != nil {
		return nil, err
	}

	if _, err = p.consume(tokenizer.CloseParens); err != nil {
		return nil, err
	}

	return expr, nil
}

// Literal
//   : NumericLit
//   | StringLit
//   | BoolLit
//   | NullLit
//   ;
func (p *Parser) literal() (ast.Node, error) {
	switch p.lookahead.Type {
	case tokenizer.Number:
		return p.numericLit()
	case tokenizer.String:
		return p.stringLit()
	case tokenizer.TrueKeyword:
		return p.boolLit(true)
	case tokenizer.FalseKeyword:
		return p.boolLit(false)
	case tokenizer.NullKeyword:
		return p.nullLit()
	default:
		return nil, &ErrUnknownLiteral{
			Type:  p.lookahead.Type,
			Value: p.lookahead.Value,
		}
	}
}

// NumericLit
//   : NUMBER
//   ;
func (p *Parser) numericLit() (ast.Node, error) {
	token, err := p.consume(tokenizer.Number)
	if err != nil {
		return nil, err
	}

	n, err := strconv.ParseInt(token.Value, 10, 64)
	if err != nil {
		return nil, err
	}

	return p.builder.NumericLit(int(n)), nil
}

// StringLit
//   : STRING
//   ;
func (p *Parser) stringLit() (ast.Node, error) {
	token, err := p.consume(tokenizer.String)
	if err != nil {
		return nil, err
	}

	return p.builder.StringLit(token.Value[1 : len(token.Value)-1]), nil
}

// BoolLit
//   : 'true'
//   | 'false'
//   ;
func (p *Parser) boolLit(v bool) (ast.Node, error) {
	tokType := tokenizer.FalseKeyword
	if v {
		tokType = tokenizer.TrueKeyword
	}

	if _, err := p.consume(tokType); err != nil {
		return nil, err
	}

	return p.builder.BoolLit(v), nil
}

// NullLit
//   : 'null'
//   ;
func (p *Parser) nullLit() (ast.Node, error) {
	if _, err := p.consume(tokenizer.NullKeyword); err != nil {
		return nil, err
	}

	return p.builder.NullLit(), nil
}

func (p *Parser) consume(tokType tokenizer.TokenType) (*tokenizer.Token, error) {
	token := p.lookahead

	if token == nil || token.Type == tokenizer.EOF {
		return nil, &ErrUnexpectedEndOfInput{Type: tokType}
	}

	if token.Type != tokType {
		return nil, &ErrUnexpectedToken{
			Type:         token.Type,
			Value:        token.Value,
			ExpectedType: tokType,
		}
	}

	var err error
	p.lookahead, err = p.tokenizer.NextToken()
	if err != nil {
		return nil, err
	}

	return token, nil
}
