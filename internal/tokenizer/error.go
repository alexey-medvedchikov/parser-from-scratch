package tokenizer

import "fmt"

type ErrUnexpectedToken struct {
	Position   int
	CodeString string
}

func (u *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token at position %d: \"%s\"", u.Position, u.CodeString)
}
