package ast

import (
	"CornyLang/token"
	"bytes"
	"strings"
)

type NodeType string

const (
	NT_PROGRAM  = "PROGRAM"
	NT_BLOCK    = "BLOCK"
	NT_NUMBER   = "NUMBER"
	NT_LET      = "LET"
	NT_RETURN   = "RETURN"
	NT_FUNCTION = "FUNCTION"
	NT_ARRAY    = "ARRAY"
	NT_HASH     = "HASH"
	NT_CALL     = "CALL"
	NT_IDENT    = "IDENT"
	NT_BINARY   = "BINARY"
	NT_UNARY    = "UNARY"
	NT_STRING   = "STRING"
	NT_BOOLEAN  = "BOOLEAN"
	NT_NULL     = "NULL"
	NT_IF       = "IF"
	NT_EXPRSTMT = "EXPRSTMT"
	NT_CLASS    = "CLASS"
)

/**
* Node Interface
 */
type Node interface {
	TokenLiteral() string
	String() string
	Type() NodeType
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

func (p *ProgramNode) Type() NodeType {
	return NT_PROGRAM
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

func (bs *BlockStmtNode) Type() NodeType {
	return NT_BLOCK
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

func (bin *BinOpNode) Type() NodeType {
	return NT_BINARY
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

func (un *UnaryNode) Type() NodeType {
	return NT_UNARY
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

func (n *NumberNode) Type() NodeType {
	return NT_NUMBER
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

func (bn *BooleanNode) Type() NodeType {
	return NT_BOOLEAN
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

func (nn *NullNode) Type() NodeType {
	return NT_NULL
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

func (sn *StringNode) Type() NodeType {
	return NT_STRING
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

func (in *IdentifierNode) Type() NodeType {
	return NT_IDENT
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

func (ln *LetNode) Type() NodeType {
	return NT_LET
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

func (es *ExpressionStatement) Type() NodeType {
	return NT_EXPRSTMT
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

func (ce *CallExprNode) Type() NodeType {
	return NT_CALL
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

func (fl *IfExprNode) Type() NodeType {
	return NT_IF
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
func (rn *ReturnNode) Type() NodeType {
	return NT_RETURN
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
func (fl *FunctionLiteralNode) Type() NodeType {
	return NT_FUNCTION
}

/**
* ClassLiteralNode
 */
type ClassLiteralNode struct {
	Token token.Token
	Body  *BlockStmtNode
}

func (cl *ClassLiteralNode) expressionNode() {}
func (cl *ClassLiteralNode) TokenLiteral() string {
	return cl.Token.Literal
}

func (cl *ClassLiteralNode) String() string {
	var out bytes.Buffer
	out.WriteString("class")
	out.WriteString(cl.Body.String())
	return out.String()
}
func (cl *ClassLiteralNode) Type() NodeType {
	return NT_CLASS
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
func (al *ArrayLiteralNode) Type() NodeType {
	return NT_ARRAY
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
func (hl *HashLiteralNode) Type() NodeType {
	return NT_HASH
}
