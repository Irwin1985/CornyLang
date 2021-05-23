package ast

import (
	"CornyLang/token"
	"bytes"
	"strings"
)

/**
* Node Interface
 */
type Node interface {
	TokenLiteral() string
	String() string
}

/**
* Statement Interface
 */
type Statement interface {
	Node
	statementNode()
}

/**
* Expression Interface
 */
type Expression interface {
	Node
	expressionNode()
}

/**
* ProgramNode
 */
type ProgramNode struct {
	Statements []Statement
}

func (p *ProgramNode) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *ProgramNode) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

/**
* BlockStmtNode
 */
type BlockStmtNode struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStmtNode) statementNode() {}
func (bs *BlockStmtNode) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStmtNode) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}
	out.WriteString("}")
	return out.String()
}

/**
* BinOpNode
 */
type BinOpNode struct {
	Token    token.Token
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (bin *BinOpNode) expressionNode() {}
func (bin *BinOpNode) TokenLiteral() string {
	return bin.Token.Literal
}
func (bin *BinOpNode) String() string {
	var out bytes.Buffer
	out.WriteString("( ")
	out.WriteString(bin.Left.String())
	out.WriteString(" " + bin.Operator.Literal + " ")
	out.WriteString(bin.Right.String())
	out.WriteString(" )")
	return out.String()
}

/**
* UnaryNode
 */
type UnaryNode struct {
	Token    token.Token
	Operator token.Token
	Right    Expression
}

func (un *UnaryNode) expressionNode() {}
func (un *UnaryNode) TokenLiteral() string {
	return un.Token.Literal
}
func (un *UnaryNode) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(un.Operator.Literal)
	out.WriteString(un.Right.String())
	out.WriteString(")")
	return out.String()
}

/**
* NumberNode
 */
type NumberNode struct {
	Token token.Token
	Value float64
}

func (n *NumberNode) expressionNode() {}
func (n *NumberNode) TokenLiteral() string {
	return n.Token.Literal
}
func (n *NumberNode) String() string {
	return n.Token.Literal
}

/**
* BooleanNode
 */
type BooleanNode struct {
	Token token.Token
	Value bool
}

func (bn *BooleanNode) expressionNode() {}
func (bn *BooleanNode) TokenLiteral() string {
	return bn.Token.Literal
}
func (bn *BooleanNode) String() string {
	return bn.Token.Literal
}

/**
* NullNode
 */
type NullNode struct {
	Token token.Token
}

func (nn *NullNode) expressionNode() {}
func (nn *NullNode) TokenLiteral() string {
	return nn.Token.Literal
}
func (nn *NullNode) String() string {
	return nn.Token.Literal
}

/**
* StringNode
 */
type StringNode struct {
	Token token.Token
	Value string
}

func (sn *StringNode) expressionNode() {}
func (sn *StringNode) TokenLiteral() string {
	return sn.Token.Literal
}
func (sn *StringNode) String() string {
	return string("\"" + sn.Token.Literal + "\"")
}

/**
* IdentifierNode
 */
type IdentifierNode struct {
	Token token.Token
	Value string
}

func (in *IdentifierNode) expressionNode() {}
func (in *IdentifierNode) TokenLiteral() string {
	return in.Token.Literal
}
func (in *IdentifierNode) String() string {
	return in.Value
}

/**
* LetNode
 */
type LetNode struct {
	Token token.Token
	Name  *IdentifierNode
	Value Expression
}

func (ln *LetNode) statementNode() {}
func (ln *LetNode) TokenLiteral() string {
	return ln.Token.Literal
}

func (ln *LetNode) String() string {
	var out bytes.Buffer
	out.WriteString("let ")
	out.WriteString(ln.Name.String())
	out.WriteString(" = ")
	out.WriteString(ln.Value.String())
	out.WriteString(";")

	return out.String()
}

/**
* ExpressionStatement
 */
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

/**
* CallExprNode
 */
type CallExprNode struct {
	Token     token.Token
	Callee    Expression
	Arguments []Expression
}

func (ce *CallExprNode) expressionNode() {}
func (ce *CallExprNode) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExprNode) String() string {
	var out bytes.Buffer
	// callee
	out.WriteString(ce.Callee.String())

	// arguments list
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

/**
* IfExprNode
 */
type IfExprNode struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStmtNode
	Alternative *BlockStmtNode
}

func (fl *IfExprNode) expressionNode() {}
func (fl *IfExprNode) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *IfExprNode) String() string {
	var out bytes.Buffer

	out.WriteString("if (")
	out.WriteString(fl.Condition.String())
	out.WriteString(")")

	out.WriteString(fl.Consequence.String())

	if fl.Alternative != nil {
		out.WriteString("else")
		out.WriteString(fl.Alternative.String())
	}

	return out.String()
}

/**
* ReturnNode
 */
type ReturnNode struct {
	Token token.Token
	Value Expression
}

func (rn *ReturnNode) statementNode() {}
func (rn *ReturnNode) TokenLiteral() string {
	return rn.Token.Literal
}
func (rn *ReturnNode) String() string {
	var out bytes.Buffer

	out.WriteString("return ")
	out.WriteString(rn.Value.String())

	return out.String()
}

/**
* FunctionLiteralNode
 */
type FunctionLiteralNode struct {
	Token      token.Token
	Parameters []*IdentifierNode
	Body       *BlockStmtNode
}

func (fl *FunctionLiteralNode) expressionNode() {}
func (fl *FunctionLiteralNode) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteralNode) String() string {
	var out bytes.Buffer
	out.WriteString("fn(")
	if fl.Parameters != nil {
		params := []string{}
		for _, param := range fl.Parameters {
			params = append(params, param.String())
		}
		out.WriteString(strings.Join(params, ", "))
	}
	out.WriteString(")")
	// function Body
	out.WriteString(fl.Body.String())
	return out.String()
}

/**
* ArrayLiteralNode
 */
type ArrayLiteralNode struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteralNode) expressionNode() {}
func (al *ArrayLiteralNode) TokenLiteral() string {
	return al.Token.Literal
}
func (al *ArrayLiteralNode) String() string {
	var out bytes.Buffer

	out.WriteString("[")

	var elements = []string{}

	for _, elem := range al.Elements {
		elements = append(elements, elem.String())
	}

	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

/**
* HashLiteralNode
 */
type HashLiteralNode struct {
	Token  token.Token
	Keys   []Expression
	Values []Expression
}

func (hl *HashLiteralNode) expressionNode() {}
func (hl *HashLiteralNode) TokenLiteral() string {
	return hl.Token.Literal
}
func (hl *HashLiteralNode) String() string {
	var out bytes.Buffer

	out.WriteString("{")

	var elements = []string{}

	for i, key := range hl.Keys {
		elements = append(elements, key.String()+":"+hl.Values[i].String())
	}

	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}
