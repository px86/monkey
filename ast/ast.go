package ast

import (
	"bytes"
	"fmt"
	"github.com/px86/monkey/token"
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
	for i, s := range p.Statements {
		out.WriteString(s.String())
		if i != len(p.Statements)-1 {
			out.WriteString(" ")
		}
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

type BlockStatement struct {
	Token      token.Token // of type token.LEFT_BRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("(block ")
	for i, s := range bs.Statements {
		out.WriteString(s.String())
		if i != len(bs.Statements)-1 {
			out.WriteString(" ")
		}
	}
	out.WriteString(")")
	return out.String()
}

type IfExpression struct {
	Token     token.Token
	Condition Expression
	ThenBlock *BlockStatement
	ElseBlock *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) String() string {
	if ie.ElseBlock != nil {
		return fmt.Sprintf("(%s %s %s %s)", token.AsString(ie.Token.Type),
			ie.Condition.String(), ie.ThenBlock.String(), ie.ElseBlock.String())
	}
	return fmt.Sprintf("(%s %s %s)", token.AsString(ie.Token.Type),
		ie.Condition.String(), ie.ThenBlock.String())
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}
func (b *Boolean) expressionNode() {}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", i.Value)
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
	return fmt.Sprintf("(%s %s)", token.AsString(pe.Operator.Type), pe.Expression.String())
}
func (pe *PrefixExpr) expressionNode() {}

type InfixExpr struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (be *InfixExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", token.AsString(be.Operator.Type), be.Left.String(), be.Right.String())
}
func (be *InfixExpr) expressionNode() {}

type FunctionExpr struct {
	Token token.Token // fn
	Args  []*Identifier
	Body  *BlockStatement
}

func (fe *FunctionExpr) String() string {
	var buff bytes.Buffer
	buff.WriteString("(" + token.AsString(fe.Token.Type) + " (")
	for i, arg := range fe.Args {
		buff.WriteString(arg.Value)
		if i != len(fe.Args)-1 {
			buff.WriteString(" ")
		}
	}
	buff.WriteString(") ")
	buff.WriteString(fe.Body.String())
	buff.WriteString(")")
	return buff.String()
}
func (fe *FunctionExpr) expressionNode() {}

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
