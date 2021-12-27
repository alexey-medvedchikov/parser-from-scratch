package ast

import (
	"bytes"
	"encoding/json"
)

type Node *concreteNode

type Fields interface{}

type concreteNode struct {
	Type NodeType
	Fields
}

func (c *concreteNode) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{
		"type": c.Type.String(),
	}
	b, err := jsonMarshal(c.Fields)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}

	return jsonMarshal(result)
}

func jsonMarshal(val interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(val)
	return buf.Bytes(), err
}

type Program struct {
	Body []Node `json:"body"`
}

type StringLit struct {
	Value string `json:"value"`
}

type NumericLit struct {
	Value int `json:"value"`
}

type BoolLit struct {
	Value bool `json:"value"`
}

type NullLit struct{}

type ExprStmt struct {
	Expr Node `json:"expr"`
}

type BlockStmt struct {
	Body []Node `json:"body"`
}

type EmptyStmt struct{}

type BinaryExpr struct {
	Op    BinaryOp `json:"op"`
	Left  Node     `json:"left"`
	Right Node     `json:"right"`
}

type BinaryOp int

const (
	InvalidBinaryOp BinaryOp = iota

	AddBinaryOp
	SubBinaryOp
	MulBinaryOp
	DivBinaryOp

	GtBinaryOp
	LtBinaryOp
	GteBinaryOp
	LteBinaryOp

	EqBinaryOp
	NeqBinaryOp
)

var binaryOpStrings = [...]string{
	"InavlidBinaryOp",

	"+", // AddBinaryOp
	"-", // SubBinaryOp
	"*", // MulBinaryOp
	"/", // DivBinaryOp

	">",  // GtBinaryOp
	"<",  // LtBinaryOp
	">=", // GteBinaryOp
	"<=", // LteBinaryOp

	"==", // EqBinaryOp
	"!=", // NeqBinaryOp
}

func (b BinaryOp) String() string {
	if b >= 0 && int(b) < len(binaryOpStrings) {
		return binaryOpStrings[b]
	}

	return binaryOpStrings[InvalidBinaryOp]
}

func (b BinaryOp) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

var binaryOpMap = func() map[string]BinaryOp {
	result := map[string]BinaryOp{}
	for i, v := range binaryOpStrings {
		result[v] = BinaryOp(i)
	}

	return result
}()

func BinaryOpFromString(v string) BinaryOp {
	op, ok := binaryOpMap[v]
	if !ok {
		return InvalidBinaryOp
	}
	return op
}

type AssignExpr struct {
	Op    AssignOp `json:"op"`
	Left  Node     `json:"left"`
	Right Node     `json:"right"`
}

type AssignOp int

const (
	InvalidAssignOp AssignOp = iota

	SimpleAssignOp
	AddAssignOp
	SubAssignOp
	MulAssignOp
	DivAssignOp
)

var assignOpStrings = [...]string{
	"InvalidAssignOp",

	"=",  // SimpleAssignOp
	"+=", // AddAssignOp
	"-=", // SubAssignOp
	"*=", // MulAssignOp
	"/=", // DivAssignOp
}

func (a AssignOp) String() string {
	if a >= 0 && int(a) < len(assignOpStrings) {
		return assignOpStrings[a]
	}

	return assignOpStrings[InvalidAssignOp]
}

func (a AssignOp) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

var assignOpMap = func() map[string]AssignOp {
	result := map[string]AssignOp{}
	for i, v := range assignOpStrings {
		result[v] = AssignOp(i)
	}

	return result
}()

func AssignOpFromString(v string) AssignOp {
	op, ok := assignOpMap[v]
	if !ok {
		return InvalidAssignOp
	}
	return op
}

type SeqExpr struct {
	Body []Node `json:"body"`
}

type NewExpr struct {
	Callee Node   `json:"callee"`
	Args   []Node `json:"args"`
}

type LogicalExpr struct {
	Op    LogicalOp `json:"op"`
	Left  Node      `json:"left"`
	Right Node      `json:"right"`
}

type ThisExpr struct{}

type LogicalOp int

const (
	InvalidLogicalOp LogicalOp = iota

	AndLogicalOp
	OrLogicalOp
)

var logicalOpStrings = [...]string{
	"InvalidLogicalOp",

	"&&", // LandBinaryOp
	"||", // LorBinaryOp
}

func (l LogicalOp) String() string {
	if l >= 0 && int(l) < len(logicalOpStrings) {
		return logicalOpStrings[l]
	}

	return logicalOpStrings[InvalidLogicalOp]
}

func (l LogicalOp) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

var logicalOpMap = func() map[string]LogicalOp {
	result := map[string]LogicalOp{}
	for i, v := range logicalOpStrings {
		result[v] = LogicalOp(i)
	}

	return result
}()

func LogicalOpFromString(v string) LogicalOp {
	op, ok := logicalOpMap[v]
	if !ok {
		return InvalidLogicalOp
	}
	return op
}

type UnaryExpr struct {
	Op  UnaryOp `json:"op"`
	Arg Node    `json:"arg"`
}

type UnaryOp int

const (
	InvalidUnaryOp UnaryOp = iota

	NotUnaryOp
	NegUnaryOp
)

var unaryOpStrings = [...]string{
	"InvalidUnaryOp",

	"!", // NotUnaryOp
	"-", // NegUnaryOp
}

func (u UnaryOp) String() string {
	if u >= 0 && int(u) < len(unaryOpStrings) {
		return unaryOpStrings[u]
	}

	return unaryOpStrings[InvalidUnaryOp]
}

func (u UnaryOp) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

var unaryOpMap = func() map[string]UnaryOp {
	result := map[string]UnaryOp{}
	for i, v := range unaryOpStrings {
		result[v] = UnaryOp(i)
	}

	return result
}()

func UnaryOpFromString(v string) UnaryOp {
	op, ok := unaryOpMap[v]
	if !ok {
		return InvalidUnaryOp
	}
	return op
}

type Identifier struct {
	Name string `json:"name"`
}

type VarStmt struct {
	Decls []Node `json:"decls"`
}

type VarDecl struct {
	ID   Node `json:"id"`
	Init Node `json:"init"`
}

type IfStmt struct {
	Cond Node `json:"cond"`
	Cons Node `json:"cons"`
	Alt  Node `json:"alt"`
}

type WhileStmt struct {
	Cond Node `json:"cond"`
	Body Node `json:"body"`
}

type DoWhileStmt struct {
	Cond Node `json:"cond"`
	Body Node `json:"body"`
}

type ForStmt struct {
	Init Node `json:"init"`
	Cond Node `json:"cond"`
	Step Node `json:"step"`
	Body Node `json:"body"`
}

type FuncDecl struct {
	Name   Node   `json:"name"`
	Params []Node `json:"params"`
	Body   Node   `json:"body"`
}

type ReturnStmt struct {
	Arg Node `json:"arg"`
}

type MemberExpr struct {
	Computed bool `json:"computed"`
	Obj      Node `json:"obj"`
	Prop     Node `json:"prop"`
}

type CallExpr struct {
	Callee Node   `json:"callee"`
	Args   []Node `json:"args"`
}

type ClassDecl struct {
	ID    Node `json:"id"`
	Super Node `json:"super"`
	Body  Node `json:"body"`
}

type SuperCall struct{}
