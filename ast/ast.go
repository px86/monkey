package ast

import (
	"bytes"
	"fmt"
	"github.com/px86/monkey/token"
	"os"
	"strconv"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	tokstr, ok := ls.Token.Value.(string) // "let"
	if !ok {
		fmt.Fprintf(os.Stderr, "expected %q, got=%v", "let", ls.Token.Value)
	}

	out.WriteString(tokstr + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string {
	s, _ := i.Token.Value.(string)
	return s
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	tokstr, _ := rs.Token.Value.(string) // "return"
	out.WriteString(tokstr + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // first token of expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BinaryOperation struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (bo *BinaryOperation) String() string {
	return "(" + bo.Left.String() + " " + bo.Right.String() + ")"
}
func (bo *BinaryOperation) expressionNode() {}

type Integer struct {
	Token token.Token
	Value int
}

func (i *Integer) String() string {
	return strconv.Itoa(i.Value)
}
func (i *Integer) expressionNode() {}

type Str struct {
	Token token.Token
	Value string
}

func (s *Str) String() string {
	return s.Value
}
func (s *Str) expressionNode() {}

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (be *BinaryExpression) String() string {
	return be.Left.String() + " " + token.TypeStr(be.Operator.Type) + " " + be.Right.String()
}
func (be *BinaryExpression) expressionNode() {}
