package parser

import (
	"fmt"

	"github.com/alexey-medvedchikov/parser-from-scratch/internal/ast"
	"github.com/alexey-medvedchikov/parser-from-scratch/internal/tokenizer"
)

type ErrUnknownLiteral struct {
	Type  tokenizer.TokenType
	Value string
}

func (e *ErrUnknownLiteral) Error() string {
	return fmt.Sprintf("unknown literal type %s: \"%s\"", e.Type, e.Value)
}

type ErrUnexpectedEndOfInput struct {
	Type tokenizer.TokenType
}

func (e *ErrUnexpectedEndOfInput) Error() string {
	return fmt.Sprintf("unexpected end of input, expected: \"%s\"", e.Type)
}

type ErrUnexpectedToken struct {
	Type         tokenizer.TokenType
	ExpectedType tokenizer.TokenType
	Value        string
}

func (e *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token, \"%v(%s)\", expected: \"%s\"", e.Type, e.Value, e.ExpectedType)
}

type ErrUnknownLogicalOp struct {
	Op string
}

func (e *ErrUnknownLogicalOp) Error() string {
	return fmt.Sprintf("unknown logical operator: \"%s\"", e.Op)
}

type ErrUnknownBinaryOp struct {
	Op string
}

func (e *ErrUnknownBinaryOp) Error() string {
	return fmt.Sprintf("unknown binary operator: \"%s\"", e.Op)
}

type ErrUnknownUnaryOp struct {
	Op string
}

func (e *ErrUnknownUnaryOp) Error() string {
	return fmt.Sprintf("unknown unary operator: \"%s\"", e.Op)
}

type ErrUnknownAssignOp struct {
	Op string
}

func (e *ErrUnknownAssignOp) Error() string {
	return fmt.Sprintf("unknown assign operator: \"%s\"", e.Op)
}

type ErrInvalidLvalue struct {
	Node ast.Node
}

func (e *ErrInvalidLvalue) Error() string {
	return fmt.Sprintf("invalid lvalue in assignment: %+v", e.Node)
}
