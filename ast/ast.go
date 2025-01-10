package ast

import (
	"bytes"
	"fmt"
	"github.com/px86/monkey/token"
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
	out.WriteString("(prog ")
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(")")
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
	return fmt.Sprintf("%q", s.Value)
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
	return fmt.Sprintf("(%s %s %s)", token.TypeStr(be.Operator.Type), be.Left.String(), be.Right.String())
}
func (be *InfixExpr) expressionNode() {}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) String() string {
	kwlet, _ := ls.Token.Value.(string) // "let"
	return fmt.Sprintf("(%s %s %s)", kwlet, ls.Name.String(), ls.Value.String())
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
	Name *Identifier
	Args []Expression
}

func (fc *FunctionCall) expressionNode() {}

func (fc *FunctionCall) String() string {
	var buff bytes.Buffer
	buff.WriteString("(")
	buff.WriteString(fc.Name.String())
	for _, arg := range fc.Args {
		buff.WriteString(" ")
		buff.WriteString(arg.String())
	}
	buff.WriteString(")")
	return buff.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) String() string {
	rtrn, _ := rs.Token.Value.(string) // "return"
	return fmt.Sprintf("(%s %s)", rtrn, rs.ReturnValue.String())
}
