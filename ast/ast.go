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

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (i *IntegerLiteral) String() string {
	return strconv.Itoa(i.Value)
}
func (i *IntegerLiteral) expressionNode() {}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) String() string {
	return s.Value
}
func (s *StringLiteral) expressionNode() {}

type PrefixExpr struct {
	Operator   token.Token
	Expression Expression
}

func (pe *PrefixExpr) String() string {
	return token.TypeStr(pe.Operator.Type) + pe.Expression.String()
}
func (pe *PrefixExpr) expressionNode() {}

type InfixExpr struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (be *InfixExpr) String() string {
	return be.Left.String() + " " + token.TypeStr(be.Operator.Type) + " " + be.Right.String()
}
func (be *InfixExpr) expressionNode() {}

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

type FunctionCall struct {
	Identifier *Identifier
	Args       []Expression
}

func (fc *FunctionCall) expressionNode() {}

func (fc *FunctionCall) String() string {
	var buff bytes.Buffer
	buff.WriteString(fc.Identifier.String() + "(")
	for _, arg := range fc.Args {
		buff.WriteString(arg.String())
	}
	return buff.String()
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
