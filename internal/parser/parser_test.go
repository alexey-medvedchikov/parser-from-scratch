package parser

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexey-medvedchikov/parser-from-scratch/internal/ast"
	"github.com/alexey-medvedchikov/parser-from-scratch/internal/tokenizer"
)

var b ast.Builder

func TestParser_Parse_Literal(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `42;`,
			wantAST: b.Program(
				b.ExprStmt(b.NumericLit(42)),
			),
		}, {
			in: `"hello";`,
			wantAST: b.Program(
				b.ExprStmt(b.StringLit(`hello`)),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Sequence(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `42, 32;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.SeqExpr(
						b.NumericLit(42),
						b.NumericLit(32),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_StatementList(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `  42; "hello";`,
			wantAST: b.Program(
				b.ExprStmt(b.NumericLit(42)),
				b.ExprStmt(b.StringLit(`hello`)),
			),
		}, {
			in: `
			// This is a comment
			42;
			/*
			This is a multiline comment
			*/
			"hello";
		`,
			wantAST: b.Program(
				b.ExprStmt(b.NumericLit(42)),
				b.ExprStmt(b.StringLit(`hello`)),
			),
		}}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_BlockStatement(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `
{
	42;
	"hello";
}
		`,
			wantAST: b.Program(
				b.BlockStmt(
					b.ExprStmt(b.NumericLit(42)),
					b.ExprStmt(b.StringLit(`hello`)),
				),
			),
		}, {
			in: `
{ // This is a comment
	42;
/*
	This is a multiline comment
*/
	"hello";
}
		`,
			wantAST: b.Program(
				b.BlockStmt(
					b.ExprStmt(b.NumericLit(42)),
					b.ExprStmt(b.StringLit(`hello`)),
				),
			),
		}, {
			in: `{ }`,
			wantAST: b.Program(
				b.BlockStmt(),
			),
		}, {
			in: `
{
	42;
	{
		"hello";
	}
}
		`,
			wantAST: b.Program(
				b.BlockStmt(
					b.ExprStmt(b.NumericLit(42)),
					b.BlockStmt(
						b.ExprStmt(b.StringLit(`hello`)),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_EmptyStatement(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `;`,
			wantAST: b.Program(
				b.EmptyStmt(),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Math(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `2 + 2;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.AddBinaryOp,
						b.NumericLit(2),
						b.NumericLit(2),
					),
				),
			),
		}, {
			in: `x + x;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.AddBinaryOp,
						b.Identifier("x"),
						b.Identifier("x"),
					),
				),
			),
		}, {
			in: `2 * 2;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.MulBinaryOp,
						b.NumericLit(2),
						b.NumericLit(2),
					),
				),
			),
		}, {
			in: `3 + 2 - 2;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.SubBinaryOp,
						b.BinaryExpr(
							ast.AddBinaryOp,
							b.NumericLit(3),
							b.NumericLit(2),
						),
						b.NumericLit(2),
					),
				),
			),
		}, {
			in: `3 + 2 * 2;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.AddBinaryOp,
						b.NumericLit(3),
						b.BinaryExpr(
							ast.MulBinaryOp,
							b.NumericLit(2),
							b.NumericLit(2),
						),
					),
				),
			),
		}, {
			in: `3 * 2 * 1;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.MulBinaryOp,
						b.BinaryExpr(
							ast.MulBinaryOp,
							b.NumericLit(3),
							b.NumericLit(2),
						),
						b.NumericLit(1),
					),
				),
			),
		}, {
			in: `(3 + 2) * 2;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.MulBinaryOp,
						b.BinaryExpr(
							ast.AddBinaryOp,
							b.NumericLit(3),
							b.NumericLit(2),
						),
						b.NumericLit(2),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Assign(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `x = 2;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.AssignExpr(
						ast.SimpleAssignOp,
						b.Identifier("x"),
						b.NumericLit(2),
					),
				),
			),
		}, {
			in: `x = y = 2;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.AssignExpr(
						ast.SimpleAssignOp,
						b.Identifier("x"),
						b.AssignExpr(
							ast.SimpleAssignOp,
							b.Identifier("y"),
							b.NumericLit(2),
						),
					),
				),
			),
		}, {
			in: `x += 2;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.AssignExpr(
						ast.AddAssignOp,
						b.Identifier("x"),
						b.NumericLit(2),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Variable(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `let x = 2;`,
			wantAST: b.Program(
				b.VarStmt(
					b.VarDecl(
						b.Identifier("x"),
						b.NumericLit(2),
					),
				),
			),
		}, {
			in: `let x;`,
			wantAST: b.Program(
				b.VarStmt(
					b.VarDecl(
						b.Identifier("x"),
						nil,
					),
				),
			),
		}, {
			in: `let x, y;`,
			wantAST: b.Program(
				b.VarStmt(
					b.VarDecl(
						b.Identifier("x"),
						nil,
					),
					b.VarDecl(
						b.Identifier("y"),
						nil,
					),
				),
			),
		}, {
			in: `let x, y = 42;`,
			wantAST: b.Program(
				b.VarStmt(
					b.VarDecl(
						b.Identifier("x"),
						nil,
					),
					b.VarDecl(
						b.Identifier("y"),
						b.NumericLit(42),
					),
				),
			),
		}, {
			in: `let x = "hello", y = 42;`,
			wantAST: b.Program(
				b.VarStmt(
					b.VarDecl(
						b.Identifier("x"),
						b.StringLit(`hello`),
					),
					b.VarDecl(
						b.Identifier("y"),
						b.NumericLit(42),
					),
				),
			),
		}, {
			in: `let x = y = 42;`,
			wantAST: b.Program(
				b.VarStmt(
					b.VarDecl(
						b.Identifier("x"),
						b.AssignExpr(
							ast.SimpleAssignOp,
							b.Identifier("y"),
							b.NumericLit(42),
						),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Relational(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `x > 0;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.GtBinaryOp,
						b.Identifier("x"),
						b.NumericLit(0),
					),
				),
			),
		}, {
			in: `x < 0;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.LtBinaryOp,
						b.Identifier("x"),
						b.NumericLit(0),
					),
				),
			),
		}, {
			in: `x >= 0;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.GteBinaryOp,
						b.Identifier("x"),
						b.NumericLit(0),
					),
				),
			),
		}, {
			in: `x <= 0;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.LteBinaryOp,
						b.Identifier("x"),
						b.NumericLit(0),
					),
				),
			),
		}, {
			in: `x + 5 > 10;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.GtBinaryOp,
						b.BinaryExpr(
							ast.AddBinaryOp,
							b.Identifier("x"),
							b.NumericLit(5),
						),
						b.NumericLit(10),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Equality(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `x == 0;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.EqBinaryOp,
						b.Identifier("x"),
						b.NumericLit(0),
					),
				),
			),
		}, {
			in: `x + 5 == 10;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.EqBinaryOp,
						b.BinaryExpr(
							ast.AddBinaryOp,
							b.Identifier("x"),
							b.NumericLit(5),
						),
						b.NumericLit(10),
					),
				),
			),
		}, {
			in: `x + 5 != 10;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.NeqBinaryOp,
						b.BinaryExpr(
							ast.AddBinaryOp,
							b.Identifier("x"),
							b.NumericLit(5),
						),
						b.NumericLit(10),
					),
				),
			),
		}, {
			in: `x + 5 > 10 == true;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.EqBinaryOp,
						b.BinaryExpr(
							ast.GtBinaryOp,
							b.BinaryExpr(
								ast.AddBinaryOp,
								b.Identifier("x"),
								b.NumericLit(5),
							),
							b.NumericLit(10),
						),
						b.BoolLit(true),
					),
				),
			),
		}, {
			in: `x + 5 > 10 == null;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.BinaryExpr(
						ast.EqBinaryOp,
						b.BinaryExpr(
							ast.GtBinaryOp,
							b.BinaryExpr(
								ast.AddBinaryOp,
								b.Identifier("x"),
								b.NumericLit(5),
							),
							b.NumericLit(10),
						),
						b.NullLit(),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Logical(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `x > 0 && y < 1;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.LogicalExpr(
						ast.AndLogicalOp,
						b.BinaryExpr(
							ast.GtBinaryOp,
							b.Identifier("x"),
							b.NumericLit(0),
						),
						b.BinaryExpr(
							ast.LtBinaryOp,
							b.Identifier("y"),
							b.NumericLit(1),
						),
					),
				),
			),
		}, {
			in: `x > 0 || y < 1;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.LogicalExpr(
						ast.OrLogicalOp,
						b.BinaryExpr(
							ast.GtBinaryOp,
							b.Identifier("x"),
							b.NumericLit(0),
						),
						b.BinaryExpr(
							ast.LtBinaryOp,
							b.Identifier("y"),
							b.NumericLit(1),
						),
					),
				),
			),
		}, {
			in: `x || y && z;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.LogicalExpr(
						ast.OrLogicalOp,
						b.Identifier("x"),
						b.LogicalExpr(
							ast.AndLogicalOp,
							b.Identifier("y"),
							b.Identifier("z"),
						),
					),
				),
			),
		}, {
			in: `(x || y) && z;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.LogicalExpr(
						ast.AndLogicalOp,
						b.LogicalExpr(
							ast.OrLogicalOp,
							b.Identifier("x"),
							b.Identifier("y"),
						),
						b.Identifier("z"),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_If(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `if (x) { x = 1; }`,
			wantAST: b.Program(
				b.IfStmt(
					b.Identifier("x"),
					b.BlockStmt(
						b.ExprStmt(
							b.AssignExpr(
								ast.SimpleAssignOp,
								b.Identifier("x"),
								b.NumericLit(1),
							),
						),
					),
					nil,
				),
			),
		}, {
			in: `if (x) x = 1;`,
			wantAST: b.Program(
				b.IfStmt(
					b.Identifier("x"),
					b.ExprStmt(
						b.AssignExpr(
							ast.SimpleAssignOp,
							b.Identifier("x"),
							b.NumericLit(1),
						),
					),
					nil,
				),
			),
		}, {
			in: `if (x) if (y) {} else {}`,
			wantAST: b.Program(
				b.IfStmt(
					b.Identifier("x"),
					b.IfStmt(
						b.Identifier("y"),
						b.BlockStmt(),
						b.BlockStmt(),
					),
					nil,
				),
			),
		}, {
			in: `
if (x) {
	x = 1;
} else {
	x = 2;
}
`,
			wantAST: b.Program(
				b.IfStmt(
					b.Identifier("x"),
					b.BlockStmt(
						b.ExprStmt(
							b.AssignExpr(
								ast.SimpleAssignOp,
								b.Identifier("x"),
								b.NumericLit(1),
							),
						),
					),
					b.BlockStmt(
						b.ExprStmt(
							b.AssignExpr(
								ast.SimpleAssignOp,
								b.Identifier("x"),
								b.NumericLit(2),
							),
						),
					),
				),
			),
		}, {
			in: `if (x > 10) { x = 1; }`,
			wantAST: b.Program(
				b.IfStmt(
					b.BinaryExpr(
						ast.GtBinaryOp,
						b.Identifier("x"),
						b.NumericLit(10),
					),
					b.BlockStmt(
						b.ExprStmt(
							b.AssignExpr(
								ast.SimpleAssignOp,
								b.Identifier("x"),
								b.NumericLit(1),
							),
						),
					),
					nil,
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Unary(t *testing.T) {
	type test struct {
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `-x;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.UnaryExpr(
						ast.NegUnaryOp,
						b.Identifier("x"),
					),
				),
			),
		}, {
			in: `!x;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.UnaryExpr(
						ast.NotUnaryOp,
						b.Identifier("x"),
					),
				),
			),
		}, {
			in: `!(x && y);`,
			wantAST: b.Program(
				b.ExprStmt(
					b.UnaryExpr(
						ast.NotUnaryOp,
						b.LogicalExpr(
							ast.AndLogicalOp,
							b.Identifier("x"),
							b.Identifier("y"),
						),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Loops(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `while (x > 10) { x -= 1;}`,
			wantAST: b.Program(
				b.WhileStmt(
					b.BinaryExpr(
						ast.GtBinaryOp,
						b.Identifier("x"),
						b.NumericLit(10),
					),
					b.BlockStmt(
						b.ExprStmt(
							b.AssignExpr(
								ast.SubAssignOp,
								b.Identifier("x"),
								b.NumericLit(1),
							),
						),
					),
				),
			),
		}, {
			in: `do { x -= 1; } while (x > 10);`,
			wantAST: b.Program(
				b.DoWhileStmt(
					b.BinaryExpr(
						ast.GtBinaryOp,
						b.Identifier("x"),
						b.NumericLit(10),
					),
					b.BlockStmt(
						b.ExprStmt(
							b.AssignExpr(
								ast.SubAssignOp,
								b.Identifier("x"),
								b.NumericLit(1),
							),
						),
					),
				),
			),
		}, {
			in: `for (let i = 0; i < 10; i += 10) {	x += 1; }`,
			wantAST: b.Program(
				b.ForStmt(
					b.VarStmt(
						b.VarDecl(
							b.Identifier("i"),
							b.NumericLit(0),
						),
					),
					b.BinaryExpr(
						ast.LtBinaryOp,
						b.Identifier("i"),
						b.NumericLit(10),
					),
					b.AssignExpr(
						ast.AddAssignOp,
						b.Identifier("i"),
						b.NumericLit(10),
					),
					b.BlockStmt(
						b.ExprStmt(
							b.AssignExpr(
								ast.AddAssignOp,
								b.Identifier("x"),
								b.NumericLit(1),
							),
						),
					),
				),
			),
		}, {
			in: `for (i = 0; i < 10; i += 10) {	x += 1; }`,
			wantAST: b.Program(
				b.ForStmt(
					b.AssignExpr(
						ast.SimpleAssignOp,
						b.Identifier("i"),
						b.NumericLit(0),
					),
					b.BinaryExpr(
						ast.LtBinaryOp,
						b.Identifier("i"),
						b.NumericLit(10),
					),
					b.AssignExpr(
						ast.AddAssignOp,
						b.Identifier("i"),
						b.NumericLit(10),
					),
					b.BlockStmt(
						b.ExprStmt(
							b.AssignExpr(
								ast.AddAssignOp,
								b.Identifier("x"),
								b.NumericLit(1),
							),
						),
					),
				),
			),
		}, {
			in: `for (j = 10, i = 0; i < 10; i += 10) { x += 1; }`,
			wantAST: b.Program(
				b.ForStmt(
					b.SeqExpr(
						b.AssignExpr(
							ast.SimpleAssignOp,
							b.Identifier("j"),
							b.NumericLit(10),
						),
						b.AssignExpr(
							ast.SimpleAssignOp,
							b.Identifier("i"),
							b.NumericLit(0),
						),
					),
					b.BinaryExpr(
						ast.LtBinaryOp,
						b.Identifier("i"),
						b.NumericLit(10),
					),
					b.AssignExpr(
						ast.AddAssignOp,
						b.Identifier("i"),
						b.NumericLit(10),
					),
					b.BlockStmt(
						b.ExprStmt(
							b.AssignExpr(
								ast.AddAssignOp,
								b.Identifier("x"),
								b.NumericLit(1),
							),
						),
					),
				),
			),
		}, {
			in: `for (; ;) {}`,
			wantAST: b.Program(
				b.ForStmt(
					nil,
					nil,
					nil,
					b.BlockStmt(),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Func(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `
def square(x) {
	return x * x;
}
`,
			wantAST: b.Program(
				b.FuncDecl(
					b.Identifier("square"),
					[]ast.Node{b.Identifier("x")},
					b.BlockStmt(
						b.ReturnStmt(
							b.BinaryExpr(
								ast.MulBinaryOp,
								b.Identifier("x"),
								b.Identifier("x"),
							),
						),
					),
				),
			),
		}, {
			in: `
def empty() {
	return;
}
`,
			wantAST: b.Program(
				b.FuncDecl(
					b.Identifier("empty"),
					nil,
					b.BlockStmt(
						b.ReturnStmt(nil),
					),
				),
			),
		}, {
			in: `
def multiply(x, y) {
	return x * y;
}
`,
			wantAST: b.Program(
				b.FuncDecl(
					b.Identifier("multiply"),
					[]ast.Node{
						b.Identifier("x"),
						b.Identifier("y"),
					},
					b.BlockStmt(
						b.ReturnStmt(
							b.BinaryExpr(
								ast.MulBinaryOp,
								b.Identifier("x"),
								b.Identifier("y"),
							),
						),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Member(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `x.y;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.MemberExpr(
						false,
						b.Identifier("x"),
						b.Identifier("y"),
					),
				),
			),
		}, {
			in: `x.y = 1;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.AssignExpr(
						ast.SimpleAssignOp,
						b.MemberExpr(
							false,
							b.Identifier("x"),
							b.Identifier("y"),
						),
						b.NumericLit(1),
					),
				),
			),
		}, {
			in: `x[0] = 1;`,
			wantAST: b.Program(
				b.ExprStmt(
					b.AssignExpr(
						ast.SimpleAssignOp,
						b.MemberExpr(
							true,
							b.Identifier("x"),
							b.NumericLit(0),
						),
						b.NumericLit(1),
					),
				),
			),
		}, {
			in: `a.b.c['d'];`,
			wantAST: b.Program(
				b.ExprStmt(
					b.MemberExpr(
						true,
						b.MemberExpr(
							false,
							b.MemberExpr(
								false,
								b.Identifier("a"),
								b.Identifier("b"),
							),
							b.Identifier("c"),
						),
						b.StringLit("d"),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_FuncCalls(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `empty();`,
			wantAST: b.Program(
				b.ExprStmt(
					b.CallExpr(
						b.Identifier("empty"),
						nil,
					),
				),
			),
		}, {
			in: `square(10);`,
			wantAST: b.Program(
				b.ExprStmt(
					b.CallExpr(
						b.Identifier("square"),
						[]ast.Node{b.NumericLit(10)},
					),
				),
			),
		}, {
			in: `multiple(10, x, 20);`,
			wantAST: b.Program(
				b.ExprStmt(
					b.CallExpr(
						b.Identifier("multiple"),
						[]ast.Node{
							b.NumericLit(10),
							b.Identifier("x"),
							b.NumericLit(20),
						},
					),
				),
			),
		}, {
			in: `clojure(10)(x, y);`,
			wantAST: b.Program(
				b.ExprStmt(
					b.CallExpr(
						b.CallExpr(
							b.Identifier("clojure"),
							[]ast.Node{b.NumericLit(10)},
						),
						[]ast.Node{
							b.Identifier("x"),
							b.Identifier("y"),
						},
					),
				),
			),
		}, {
			in: `console.log(x, "hello");`,
			wantAST: b.Program(
				b.ExprStmt(
					b.CallExpr(
						b.MemberExpr(
							false,
							b.Identifier("console"),
							b.Identifier("log"),
						),
						[]ast.Node{
							b.Identifier("x"),
							b.StringLit("hello"),
						},
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Complex(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `

// This is a string
let s = "Hello, world!";
/*
 This is an integer
*/
let i = 0;

def square(x) {
	return x * x;
}

while (i < s.length) {
	console.log(i, s[i]);
	square(2 + i);
	getCallback()();
	i += 1;
}

`,
			wantAST: b.Program(
				b.VarStmt(
					b.VarDecl(
						b.Identifier("s"),
						b.StringLit("Hello, world!"),
					),
				),
				b.VarStmt(
					b.VarDecl(
						b.Identifier("i"),
						b.NumericLit(0),
					),
				),
				b.FuncDecl(
					b.Identifier("square"),
					[]ast.Node{b.Identifier("x")},
					b.BlockStmt(
						b.ReturnStmt(
							b.BinaryExpr(
								ast.MulBinaryOp,
								b.Identifier("x"),
								b.Identifier("x"),
							),
						),
					),
				),
				b.WhileStmt(
					b.BinaryExpr(
						ast.LtBinaryOp,
						b.Identifier("i"),
						b.MemberExpr(
							false,
							b.Identifier("s"),
							b.Identifier("length"),
						),
					),
					b.BlockStmt(
						b.ExprStmt(
							b.CallExpr(
								b.MemberExpr(
									false,
									b.Identifier("console"),
									b.Identifier("log"),
								),
								[]ast.Node{
									b.Identifier("i"),
									b.MemberExpr(
										true,
										b.Identifier("s"),
										b.Identifier("i"),
									),
								},
							),
						),
						b.ExprStmt(
							b.CallExpr(
								b.Identifier("square"),
								[]ast.Node{
									b.BinaryExpr(
										ast.AddBinaryOp,
										b.NumericLit(2),
										b.Identifier("i"),
									),
								},
							),
						),

						b.ExprStmt(
							b.CallExpr(
								b.CallExpr(
									b.Identifier("getCallback"),
									nil,
								),
								nil,
							),
						),
						b.ExprStmt(
							b.AssignExpr(
								ast.AddAssignOp,
								b.Identifier("i"),
								b.NumericLit(1),
							),
						),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_Class(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `
class Point {
	def constructor(x, y) {
		this.x = x;
		this.y = y;
	}

	def calc() {
		return this.x + this.y;
	}
}
`,
			wantAST: b.Program(
				b.ClassDecl(
					b.Identifier("Point"),
					nil,
					b.BlockStmt(
						b.FuncDecl(
							b.Identifier("constructor"),
							[]ast.Node{
								b.Identifier("x"),
								b.Identifier("y"),
							},
							b.BlockStmt(
								b.ExprStmt(
									b.AssignExpr(
										ast.SimpleAssignOp,
										b.MemberExpr(
											false,
											b.ThisExpr(),
											b.Identifier("x"),
										),
										b.Identifier("x"),
									),
								),
								b.ExprStmt(
									b.AssignExpr(
										ast.SimpleAssignOp,
										b.MemberExpr(
											false,
											b.ThisExpr(),
											b.Identifier("y"),
										),
										b.Identifier("y"),
									),
								),
							),
						),
						b.FuncDecl(
							b.Identifier("calc"),
							nil,
							b.BlockStmt(
								b.ReturnStmt(
									b.BinaryExpr(
										ast.AddBinaryOp,
										b.MemberExpr(
											false,
											b.ThisExpr(),
											b.Identifier("x"),
										),
										b.MemberExpr(
											false,
											b.ThisExpr(),
											b.Identifier("y"),
										),
									),
								),
							),
						),
					),
				),
			),
		}, {
			in: `
class Point3D extends Point {
	def constructor(x, y, z) {
		super(x, y);
		this.z = z;
	}

	def calc() {
		return super() + this.z;
	}
}
`,
			wantAST: b.Program(
				b.ClassDecl(
					b.Identifier("Point3D"),
					b.Identifier("Point"),
					b.BlockStmt(
						b.FuncDecl(
							b.Identifier("constructor"),
							[]ast.Node{
								b.Identifier("x"),
								b.Identifier("y"),
								b.Identifier("z"),
							},
							b.BlockStmt(
								b.ExprStmt(
									b.CallExpr(
										b.SuperCall(),
										[]ast.Node{
											b.Identifier("x"),
											b.Identifier("y"),
										},
									),
								),
								b.ExprStmt(
									b.AssignExpr(
										ast.SimpleAssignOp,
										b.MemberExpr(
											false,
											b.ThisExpr(),
											b.Identifier("z"),
										),
										b.Identifier("z"),
									),
								),
							),
						),
						b.FuncDecl(
							b.Identifier("calc"),
							nil,
							b.BlockStmt(
								b.ReturnStmt(
									b.BinaryExpr(
										ast.AddBinaryOp,
										b.CallExpr(
											b.SuperCall(),
											nil,
										),
										b.MemberExpr(
											false,
											b.ThisExpr(),
											b.Identifier("z"),
										),
									),
								),
							),
						),
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func TestParser_Parse_New(t *testing.T) {
	type test struct {
		name    string
		in      string
		wantAST ast.Node
	}
	tests := []test{
		{
			in: `new Point3D(10, 20, 30);`,
			wantAST: b.Program(
				b.ExprStmt(
					b.NewExpr(
						b.Identifier("Point3D"),
						[]ast.Node{
							b.NumericLit(10),
							b.NumericLit(20),
							b.NumericLit(30),
						},
					),
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			testOk(t, tc.in, tc.wantAST)
		})
	}
}

func testOk(t *testing.T, in string, wantAST ast.Node) {
	tok := tokenizer.NewTokenizer(tokenizer.DefaultRules, in)
	p := NewParser(tok, b)
	node, err := p.Parse()
	assert.NoError(t, err)
	if !assert.Exactly(t, wantAST, node) {
		assert.Exactly(t, dumpJSON(t, wantAST), dumpJSON(t, node))
	}
}

func dumpJSON(t *testing.T, node ast.Node) string {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(node)
	assert.NoError(t, err)
	return buf.String()
}
