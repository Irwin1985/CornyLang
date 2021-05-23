package parser

import (
	"CornyLang/ast"
	"CornyLang/lexer"
	"CornyLang/token"
	"fmt"
	"os"
	"strconv"
)

type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer}
	parser.nextToken() // feed peekToken
	parser.nextToken() // feed curToken
	return parser
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) advance(tokenType token.TokenType) {
	if p.curToken.Type == tokenType {
		p.nextToken()
	} else {
		msg := fmt.Sprintf("unexpected token. got %s instead of %s", tokenType, p.curToken.Type)
		p.errors = append(p.errors, msg)
	}
}

// program ::= ( statement )*
func (p *Parser) Program() *ast.ProgramNode {
	var programNode = &ast.ProgramNode{}
	programNode.Statements = []ast.Statement{}

	for p.curToken.Type != token.TT_EOF {
		var stmt = p.parseStatement()
		if stmt != nil {
			programNode.Statements = append(programNode.Statements, stmt)
		}
	}

	return programNode
}

// parseBlockStmt ::= LBRACE ( statement )* RBRACE
func (p *Parser) parseBlockStmt() *ast.BlockStmtNode {
	var blockNode = &ast.BlockStmtNode{Token: p.curToken}

	blockNode.Statements = []ast.Statement{}

	p.advance(token.TT_LBRACE)
	for p.curToken.Type != token.TT_RBRACE {
		blockNode.Statements = append(blockNode.Statements, p.parseStatement())
	}
	p.advance(token.TT_RBRACE)

	return blockNode
}

// statement ::= parseLetStatement
// 			  	| parseReturnStatement
// 			  	| parseFunctionStatement
//				| parseExpressionStatement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.TT_LET:
		return p.parseLetStmt()
	case token.TT_RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExpressionStmt()
	}

}

// parseLetStmt ::= let IDENTIFIER = parseExpression
func (p *Parser) parseLetStmt() ast.Statement {
	var letNode = &ast.LetNode{Token: p.curToken}

	p.advance(token.TT_LET)
	letNode.Name = p.parseIdentifier()

	p.advance(token.TT_ASSIGN)
	letNode.Value = p.parseExpression()

	if p.curToken.Type == token.TT_SEMICOLON {
		p.advance(token.TT_SEMICOLON)
	}

	return letNode
}

// parseReturnStmt
func (p *Parser) parseReturnStmt() ast.Statement {
	var returnNode = &ast.ReturnNode{Token: p.curToken}
	p.advance(token.TT_RETURN)

	returnNode.Value = p.parseExpression()

	if p.curToken.Type == token.TT_SEMICOLON {
		p.advance(token.TT_SEMICOLON)
	}

	return returnNode
}

// parseExpressionStmt ::= logicOr
func (p *Parser) parseExpressionStmt() *ast.ExpressionStatement {
	var expressionStmt = &ast.ExpressionStatement{}
	expressionStmt.Expression = p.parseExpression()
	if p.curToken.Type == token.TT_SEMICOLON {
		p.advance(token.TT_SEMICOLON)
	}
	return expressionStmt
}

func (p *Parser) parseExpression() ast.Expression {
	return p.logicOr()
}

func (p *Parser) logicOr() ast.Expression {
	var node = p.logicAnd()
	for p.curToken.Type == token.TT_OR {
		var token = p.curToken
		p.advance(token.Type)
		node = &ast.BinOpNode{Left: node, Operator: token, Right: p.logicAnd()}
	}
	return node
}

func (p *Parser) logicAnd() ast.Expression {
	var node = p.equality()
	for p.curToken.Type == token.TT_AND {
		var token = p.curToken
		p.advance(token.Type)
		node = &ast.BinOpNode{Left: node, Operator: token, Right: p.equality()}
	}
	return node
}

func (p *Parser) equality() ast.Expression {
	var node = p.comparison()
	for p.curToken.Type == token.TT_EQUAL || p.curToken.Type == token.TT_NOT_EQ {
		var token = p.curToken
		p.advance(token.Type)
		node = &ast.BinOpNode{Left: node, Operator: token, Right: p.comparison()}
	}
	return node
}

func (p *Parser) comparison() ast.Expression {
	var node = p.term()
	for p.curToken.Type == token.TT_LESS ||
		p.curToken.Type == token.TT_LESS_EQ ||
		p.curToken.Type == token.TT_GREATER ||
		p.curToken.Type == token.TT_GREATER_EQ {
		var token = p.curToken
		p.advance(token.Type)
		node = &ast.BinOpNode{Left: node, Operator: token, Right: p.term()}
	}
	return node
}

func (p *Parser) term() ast.Expression {
	var node = p.factor()
	for p.curToken.Type == token.TT_PLUS || p.curToken.Type == token.TT_MINUS {
		var token = p.curToken
		p.advance(token.Type)
		node = &ast.BinOpNode{Left: node, Operator: token, Right: p.factor()}
	}
	return node
}

func (p *Parser) factor() ast.Expression {
	var node = p.unary()
	for p.curToken.Type == token.TT_MUL || p.curToken.Type == token.TT_DIV {
		var token = p.curToken
		p.advance(token.Type)
		node = &ast.BinOpNode{Left: node, Operator: token, Right: p.unary()}
	}
	return node
}

// unary ::= ( ( '-' | '!' ) unary )* | call
func (p *Parser) unary() ast.Expression {
	for p.curToken.Type == token.TT_MINUS || p.curToken.Type == token.TT_NOT {
		var token = p.curToken
		p.advance(token.Type)
		return &ast.UnaryNode{Operator: token, Right: p.unary()}
	}
	return p.call()
}

// call ::= primary '(' arguments? ')'
func (p *Parser) call() ast.Expression {
	var node = p.primary()
	for {
		if p.curToken.Type == token.TT_LPAREN ||
			p.curToken.Type == token.TT_LBRACKET ||
			p.curToken.Type == token.TT_DOT {
			node = p.parseCallExpr(node)
		} else if p.curToken.Type == token.TT_QUESTION {
			node = p.parseTernaryExpr(node)
		} else {
			break
		}
	}
	return node
}

// primary ::= NUMBRE | STRING | TRUE | FALSE | NULL | IDENT | FUNCTION | IF-ELSE
func (p *Parser) primary() ast.Expression {
	var tok = p.curToken
	switch tok.Type {
	case token.TT_NUMBER:
		p.advance(token.TT_NUMBER)
		value, _ := strconv.ParseFloat(tok.Literal, 64)
		return &ast.NumberNode{Token: tok, Value: value}
	case token.TT_TRUE:
		p.advance(token.TT_TRUE)
		return &ast.BooleanNode{Token: tok, Value: true}
	case token.TT_FALSE:
		p.advance(token.TT_FALSE)
		return &ast.BooleanNode{Token: tok, Value: false}
	case token.TT_NULL:
		p.advance(token.TT_NULL)
		return &ast.NullNode{Token: tok}
	case token.TT_IDENT:
		return p.parseIdentifier()
	case token.TT_LPAREN:
		p.advance(token.TT_LPAREN)
		var resultExpr = p.parseExpression()
		p.advance(token.TT_RPAREN)
		return resultExpr
	case token.TT_STRING:
		p.advance(token.TT_STRING)
		return &ast.StringNode{Token: tok, Value: tok.Literal}
	case token.TT_FUNCTION:
		return p.parseFunctionLiteral()
	case token.TT_LBRACKET:
		return p.parseArrayLiteral()
	case token.TT_LBRACE:
		return p.parseHashLiteral()
	case token.TT_IF:
		return p.parseIfExpr()
	default:
		fmt.Printf("Unknown expression token: %s\n", tok.Literal)
		os.Exit(4)
	}
	return nil
}

// parseCallExpr ::= FUNCTIONCALL | ARRAYCALL | HASHCALL
func (p *Parser) parseCallExpr(callee ast.Expression) ast.Expression {
	var callExprNode = &ast.CallExprNode{Token: p.curToken}
	callExprNode.Callee = callee
	callExprNode.Arguments = []ast.Expression{}

	if p.curToken.Type == token.TT_LBRACKET { // array or hash call
		p.advance(token.TT_LBRACKET)
		if p.curToken.Type == token.TT_RBRACKET {
			fmt.Print("Invalid subscript reference for array or hash types\n")
			os.Exit(4)
		}

		callExprNode.Arguments = append(callExprNode.Arguments, p.parseExpression())
		p.advance(token.TT_RBRACKET)
	} else if p.curToken.Type == token.TT_LPAREN { // function call
		p.advance(token.TT_LPAREN)
		if p.curToken.Type != token.TT_RPAREN {
			callExprNode.Arguments = append(callExprNode.Arguments, p.parseExpression())
			for p.curToken.Type == token.TT_COMMA {
				p.advance(token.TT_COMMA)
				callExprNode.Arguments = append(callExprNode.Arguments, p.parseExpression())
			}
		}
		p.advance(token.TT_RPAREN)
	} else if p.curToken.Type == token.TT_DOT {
		p.advance(token.TT_DOT)
		callExprNode.Arguments = append(callExprNode.Arguments, p.parseExpression())
	}
	return callExprNode
}

// parseTernaryExpr
func (p *Parser) parseTernaryExpr(conditionNode ast.Expression) ast.Expression {
	var ifNode = &ast.IfExprNode{Token: p.curToken}
	p.advance(token.TT_QUESTION)

	ifNode.Condition = conditionNode
	ifNode.Consequence = &ast.BlockStmtNode{Statements: []ast.Statement{p.parseStatement()}}

	p.advance(token.TT_COLON)
	ifNode.Alternative = &ast.BlockStmtNode{Statements: []ast.Statement{p.parseStatement()}}

	return ifNode
}

// parseFunctionLiteral
func (p *Parser) parseFunctionLiteral() ast.Expression {
	var functionNode = &ast.FunctionLiteralNode{Token: p.curToken}
	p.advance(token.TT_FUNCTION)
	p.advance(token.TT_LPAREN)

	if p.curToken.Type != token.TT_RPAREN {
		functionNode.Parameters = []*ast.IdentifierNode{}
		functionNode.Parameters = append(functionNode.Parameters, p.parseIdentifier())
		for p.curToken.Type == token.TT_COMMA {
			p.advance(token.TT_COMMA)
			functionNode.Parameters = append(functionNode.Parameters, p.parseIdentifier())
		}
	}
	p.advance(token.TT_RPAREN)
	functionNode.Body = p.parseBlockStmt()

	return functionNode
}

// parseIdentifier
func (p *Parser) parseIdentifier() *ast.IdentifierNode {
	var identifierNode = &ast.IdentifierNode{Token: p.curToken}
	identifierNode.Value = p.curToken.Literal

	p.advance(token.TT_IDENT)

	return identifierNode
}

// parseArrayLiteral
func (p *Parser) parseArrayLiteral() *ast.ArrayLiteralNode {
	var arrayNode = &ast.ArrayLiteralNode{Token: p.curToken}

	p.advance(token.TT_LBRACKET)
	if p.curToken.Type != token.TT_RBRACKET {
		arrayNode.Elements = []ast.Expression{}
		arrayNode.Elements = append(arrayNode.Elements, p.parseExpression())
		for p.curToken.Type == token.TT_COMMA {
			p.advance(token.TT_COMMA)
			arrayNode.Elements = append(arrayNode.Elements, p.parseExpression())
		}
	}
	p.advance(token.TT_RBRACKET)

	return arrayNode
}

// parseHashLiteral
func (p *Parser) parseHashLiteral() *ast.HashLiteralNode {
	var hashNode = &ast.HashLiteralNode{Token: p.curToken}
	p.advance(token.TT_LBRACE)
	if p.curToken.Type != token.TT_RBRACE {
		hashNode.Keys = []ast.Expression{}
		hashNode.Values = []ast.Expression{}

		// parse key-value pairs expressions
		hashNode.Keys = append(hashNode.Keys, p.parseExpression())
		p.advance(token.TT_COLON)
		hashNode.Values = append(hashNode.Values, p.parseExpression())

		for p.curToken.Type == token.TT_COMMA {
			p.advance(token.TT_COMMA)
			// parse key-value pairs expressions
			// TODO: apply refactoring in duplicated code.
			hashNode.Keys = append(hashNode.Keys, p.parseExpression())
			p.advance(token.TT_COLON)
			hashNode.Values = append(hashNode.Values, p.parseExpression())
		}

	}
	p.advance(token.TT_RBRACE)
	return hashNode
}

// parseIfExpr
func (p *Parser) parseIfExpr() *ast.IfExprNode {
	var ifNode = &ast.IfExprNode{Token: p.curToken}

	p.advance(token.TT_IF)

	p.advance(token.TT_LPAREN)
	ifNode.Condition = p.parseExpression()
	p.advance(token.TT_RPAREN)

	ifNode.Consequence = p.parseBlockStmt()

	if p.curToken.Type == token.TT_ELSE {
		p.advance(token.TT_ELSE)
		ifNode.Alternative = p.parseBlockStmt()
	}

	return ifNode
}
