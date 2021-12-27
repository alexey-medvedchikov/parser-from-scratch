package ast

type NodeType int

const (
	NumericLitType NodeType = iota
	StringLitType
	BoolLitType
	NullLitType
	SuperCallType
	ProgramType
	ExprStmtType
	BlockStmtType
	EmptyStmtType
	MemberExprType
	BinaryExprType
	LogicalExprType
	UnaryExprType
	AssignExprType
	SeqExprType
	ThisExprType
	NewExprType
	CallExprType
	IdentifierType
	VarStmtType
	VarDeclType
	IfStmtType
	WhileStmtType
	DoWhileStmtType
	ForStmtType
	FuncDeclType
	ClassDeclType
	ReturnStmtType
)

var nodeTypeNames = [...]string{
	"NumericLitType",
	"StringLitType",
	"BoolLitType",
	"NullLitType",
	"SuperCallType",
	"ProgramType",
	"ExprStmtType",
	"BlockStmtType",
	"EmptyStmtType",
	"MemberExprType",
	"BinaryExprType",
	"LogicalExprType",
	"UnaryExprType",
	"AssignExprType",
	"SeqExprType",
	"ThisExprType",
	"NewExprType",
	"CallExprType",
	"IdentifierType",
	"VarStmtType",
	"VarDeclType",
	"IfStmtType",
	"WhileStmtType",
	"DoWhileStmtType",
	"ForStmtType",
	"FuncDeclType",
	"ClassDeclType",
	"ReturnStmtType",
}

func (n NodeType) String() string {
	if n >= 0 && int(n) < len(nodeTypeNames) {
		return nodeTypeNames[n]
	}

	return ""
}

type Builder struct{}

func (b Builder) Program(body ...Node) Node {
	return &concreteNode{
		Type: ProgramType,
		Fields: &Program{
			Body: body,
		},
	}
}

func (b Builder) StringLit(s string) Node {
	return &concreteNode{
		Type: StringLitType,
		Fields: &StringLit{
			Value: s,
		},
	}
}

func (b Builder) NumericLit(n int) Node {
	return &concreteNode{
		Type: NumericLitType,
		Fields: &NumericLit{
			Value: n,
		},
	}
}

func (b Builder) BoolLit(v bool) Node {
	return &concreteNode{
		Type: BoolLitType,
		Fields: &BoolLit{
			Value: v,
		},
	}
}

func (b Builder) NullLit() Node {
	return &concreteNode{
		Type:   NullLitType,
		Fields: &NullLit{},
	}
}

func (b Builder) ExprStmt(expr Node) Node {
	return &concreteNode{
		Type: ExprStmtType,
		Fields: &ExprStmt{
			Expr: expr,
		},
	}
}

func (b Builder) BlockStmt(body ...Node) Node {
	return &concreteNode{
		Type: BlockStmtType,
		Fields: &BlockStmt{
			Body: body,
		},
	}
}

func (b Builder) EmptyStmt() Node {
	return &concreteNode{
		Type:   EmptyStmtType,
		Fields: &EmptyStmt{},
	}
}

func (b Builder) BinaryExpr(op BinaryOp, left Node, right Node) Node {
	return &concreteNode{
		Type: BinaryExprType,
		Fields: &BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
		},
	}
}

func (b Builder) UnaryExpr(op UnaryOp, arg Node) Node {
	return &concreteNode{
		Type: UnaryExprType,
		Fields: &UnaryExpr{
			Op:  op,
			Arg: arg,
		},
	}
}

func (b Builder) LogicalExpr(op LogicalOp, left Node, right Node) Node {
	return &concreteNode{
		Type: LogicalExprType,
		Fields: &LogicalExpr{
			Op:    op,
			Left:  left,
			Right: right,
		},
	}
}

func (b Builder) AssignExpr(op AssignOp, left Node, right Node) Node {
	return &concreteNode{
		Type: AssignExprType,
		Fields: &AssignExpr{
			Op:    op,
			Left:  left,
			Right: right,
		},
	}
}

func (b Builder) SeqExpr(body ...Node) Node {
	return &concreteNode{
		Type: SeqExprType,
		Fields: &SeqExpr{
			Body: body,
		},
	}
}

func (b Builder) Identifier(name string) Node {
	return &concreteNode{
		Type: IdentifierType,
		Fields: &Identifier{
			Name: name,
		},
	}
}

func (b Builder) VarStmt(decl ...Node) Node {
	return &concreteNode{
		Type: VarStmtType,
		Fields: &VarStmt{
			Decls: decl,
		},
	}
}

func (b Builder) VarDecl(id Node, init Node) Node {
	return &concreteNode{
		Type: VarDeclType,
		Fields: &VarDecl{
			ID:   id,
			Init: init,
		},
	}
}

func (b Builder) IfStmt(cond Node, cons Node, alt Node) Node {
	return &concreteNode{
		Type: IfStmtType,
		Fields: &IfStmt{
			Cond: cond,
			Cons: cons,
			Alt:  alt,
		},
	}
}

func (b Builder) WhileStmt(cond Node, body Node) Node {
	return &concreteNode{
		Type: WhileStmtType,
		Fields: &WhileStmt{
			Cond: cond,
			Body: body,
		},
	}
}

func (b Builder) DoWhileStmt(cond Node, body Node) Node {
	return &concreteNode{
		Type: DoWhileStmtType,
		Fields: &DoWhileStmt{
			Cond: cond,
			Body: body,
		},
	}
}

func (b Builder) ForStmt(init Node, cond Node, step Node, body Node) Node {
	return &concreteNode{
		Type: ForStmtType,
		Fields: &ForStmt{
			Init: init,
			Cond: cond,
			Step: step,
			Body: body,
		},
	}
}

func (b Builder) FuncDecl(name Node, params []Node, body Node) Node {
	return &concreteNode{
		Type: FuncDeclType,
		Fields: &FuncDecl{
			Name:   name,
			Params: params,
			Body:   body,
		},
	}
}

func (b Builder) ReturnStmt(arg Node) Node {
	return &concreteNode{
		Type: ReturnStmtType,
		Fields: &ReturnStmt{
			Arg: arg,
		},
	}
}

func (b Builder) MemberExpr(computed bool, obj Node, prop Node) Node {
	return &concreteNode{
		Type: MemberExprType,
		Fields: &MemberExpr{
			Computed: computed,
			Obj:      obj,
			Prop:     prop,
		},
	}
}

func (b Builder) CallExpr(callee Node, args []Node) Node {
	return &concreteNode{
		Type: CallExprType,
		Fields: &CallExpr{
			Callee: callee,
			Args:   args,
		},
	}
}

func (b Builder) ClassDecl(id Node, super Node, body Node) Node {
	return &concreteNode{
		Type: ClassDeclType,
		Fields: &ClassDecl{
			ID:    id,
			Super: super,
			Body:  body,
		},
	}
}

func (b Builder) SuperCall() Node {
	return &concreteNode{
		Type:   SuperCallType,
		Fields: &SuperCall{},
	}
}

func (b Builder) NewExpr(callee Node, args []Node) Node {
	return &concreteNode{
		Type: NewExprType,
		Fields: &NewExpr{
			Callee: callee,
			Args:   args,
		},
	}
}

func (b Builder) ThisExpr() Node {
	return &concreteNode{
		Type:   ThisExprType,
		Fields: &ThisExpr{},
	}
}
