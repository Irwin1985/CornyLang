package evaluator

import (
	"CornyLang/ast"
	"CornyLang/object"
	"CornyLang/token"
	"fmt"
)

var TRUE = &object.BooleanObj{Value: true}
var FALSE = &object.BooleanObj{Value: false}
var NULL = &object.NullObj{}

/**
* Helpers functions
 */
func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.ProgramNode:
		return evalProgram(node, env)
	case *ast.BlockStmtNode:
		return evalBlock(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.NumberNode:
		return &object.NumberObj{Value: node.Value}
	case *ast.StringNode:
		return &object.StringObj{Value: node.Value}
	case *ast.BooleanNode:
		return evalBooleanNode(node)
	case *ast.NullNode:
		return NULL
	case *ast.FunctionLiteralNode:
		return evalFunctionLiteralNode(node, env)
	case *ast.ArrayLiteralNode:
		return evalArrayLiteral(node, env)
	case *ast.HashLiteralNode:
		return evalHashLiteral(node, env)
	case *ast.ReturnNode:
		return Eval(node.Value, env)
	case *ast.LetNode:
		return evalLet(node, env)
	case *ast.IfExprNode:
		return evalIfExpr(node, env)
	case *ast.IdentifierNode:
		return evalIdentifier(node, env)
	case *ast.CallExprNode:
		return evalCallExpr(node, env)
	case *ast.UnaryNode:
		return evalUnaryExpr(node, env)
	case *ast.BinOpNode:
		return evalBinaryExpr(node, env)
	default:
		return nil
	}
}

func evalProgram(program *ast.ProgramNode, env *object.Environment) object.Object {
	var result object.Object

	for _, statememt := range program.Statements {
		result = Eval(statememt, env)
		switch result := result.(type) {
		case *object.ReturnObj:
			return result.Value // return the value object.
		case *object.ErrorObj:
			return result // return the error object.
		}
	}

	return result
}

func evalBlock(block *ast.BlockStmtNode, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			if result.Type() == object.ERROR_OBJ || result.Type() == object.RETURN_OBJ {
				return result
			}
		}
	}

	return result
}

func evalLet(letNode *ast.LetNode, env *object.Environment) object.Object {
	var valueObj = Eval(letNode.Value, env)
	if isError(valueObj) {
		return valueObj
	}
	// register in symbol table
	env.Set(letNode.Name.Value, valueObj)

	return valueObj
}

// evalIfExpr
func evalIfExpr(ifNode *ast.IfExprNode, env *object.Environment) object.Object {
	var conditionObj = Eval(ifNode.Condition, env)
	if isError(conditionObj) {
		return conditionObj
	}
	boolConditionObj, ok := conditionObj.(*object.BooleanObj)
	if !ok {
		return object.NewError("Invalid data type for if condition.\n")
	}
	if boolConditionObj.Value {
		return Eval(ifNode.Consequence, env)
	} else {
		if ifNode.Alternative != nil {
			return Eval(ifNode.Alternative, env)
		} else {
			return NULL
		}
	}
}

// evalCallExpr
func evalCallExpr(callExpr *ast.CallExprNode, env *object.Environment) object.Object {
	// eval the callee
	var calleeObj = Eval(callExpr.Callee, env)

	if isError(calleeObj) {
		return calleeObj
	}

	var args []object.Object
	if callExpr.Token.Type != token.TT_DOT {
		// eval arguments
		args = evalExpressions(callExpr.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
	} else {
		if len(callExpr.Arguments) != 1 {
			return object.NewError("wrong number of arguments for builtin function.")
		}
		if callExprNode, ok := callExpr.Arguments[0].(*ast.CallExprNode); ok {
			// check for identifierNode type
			if identifierNode, ok := callExprNode.Callee.(*ast.IdentifierNode); ok {
				// identifierNode must be resolved into builtinObj
				var builtinObj = evalIdentifier(identifierNode, env)
				if builtinObj, ok := builtinObj.(*object.BuiltinObj); ok {
					// this builtin function takes extra arguments?
					if !builtinObj.ExtraArgs && len(callExprNode.Arguments) > 0 {
						return object.NewError("Unexpected arguments.")
					}
					var newArgs = []object.Object{calleeObj} // we always pass the callee reference.
					if builtinObj.ExtraArgs {
						// pass extra arguments gathered by `callExpr.Arguments`
						var evaluatedArgs = evalExpressions(callExprNode.Arguments, env)
						newArgs = append(newArgs, evaluatedArgs...)
					}
					// execute the builtin function and get the result
					var resultObj = builtinObj.Func(newArgs...)
					// get the identifier from callExpr.Callee
					if identifierNode, ok := callExpr.Callee.(*ast.IdentifierNode); ok {
						env.Set(identifierNode.Value, resultObj)
					}
					return resultObj

				}
			}
		}
		return object.NewError(fmt.Sprintf("Invalid argument for builtin function using the dot(.) semantic: %s\n", callExpr.Arguments[0].String()))
	}

	switch calleeObj := calleeObj.(type) {
	case *object.FunctionObj:
		return evalFunction(calleeObj, args)
	case *object.ArrayObj:
		if callExpr.Token.Type == token.TT_LBRACKET {
			return evalArrayAccess(calleeObj, args)
		}
	case *object.HashObj:
		if callExpr.Token.Type == token.TT_LBRACKET {
			return evalHashAccess(calleeObj, args)
		}
	case *object.StringObj:
		if callExpr.Token.Type == token.TT_LBRACKET {
			return evalStringAccess(calleeObj, args)
		} else if callExpr.Token.Type == token.TT_DOT {
			return NULL
		}
	case *object.BuiltinObj: // NOTE: i'm not sure if this is necessary
		return calleeObj.Func(args...)
	default:
		//return object.NewError(fmt.Sprintf("not a function: %s\n", calleeObj.Type()))
		return calleeObj // TODO: this must be checked deeper.
	}
	return nil
}

func evalFunction(functionObj *object.FunctionObj, args []object.Object) object.Object {
	// 1. create new environment and close it with FunctionObj.Env
	var newEnv = object.NewEnvironment()
	newEnv.Parent = functionObj.Env

	// 2. check for function arity
	var numArgs = len(args)
	var numParams = len(functionObj.Parameters)
	if numArgs != numParams {
		return object.NewError(fmt.Sprintf("Unexpected arguments, got: %d, want: %d\n", numArgs, numParams))
	}
	// 3. fill the new environment with arguments
	for paramIdx, param := range functionObj.Parameters {
		newEnv.Set(param.Value, args[paramIdx])
	}
	// 4. execute the function with new environment
	var resultObj = Eval(functionObj.Body, newEnv)
	if isError(resultObj) {
		return resultObj
	}
	// check for return object
	if resultObj, ok := resultObj.(*object.ReturnObj); ok {
		return resultObj.Value
	}

	return resultObj
}

// evalArrayAccess
func evalArrayAccess(arrayObj *object.ArrayObj, args []object.Object) object.Object {
	indexObj, ok := args[0].(*object.NumberObj)
	if !ok {
		return object.NewError("Invalid subscript reference.")
	}
	// check for out of bounds
	var index = int(indexObj.Value)
	if index < 0 || index >= len(arrayObj.Elements) {
		return object.NewError("Index out of bounds.")
	}

	return arrayObj.Elements[index]
}

// evalHashAccess
func evalHashAccess(hashObj *object.HashObj, args []object.Object) object.Object {
	indexObj, ok := args[0].(*object.StringObj)
	if !ok {
		return object.NewError("Invalid subscript reference")
	}
	var key = indexObj.Value
	valueObj, ok := hashObj.Elements[key]
	if !ok {
		return NULL
	}
	return valueObj
}

// evalStringAccess
func evalStringAccess(stringObj *object.StringObj, args []object.Object) object.Object {
	indexObj, ok := args[0].(*object.NumberObj)
	if !ok {
		return object.NewError("Invalid subscript reference")
	}
	var index = int(indexObj.Value)
	if index < 0 || index >= len(stringObj.Value) {
		return object.NewError("Index out of bounds")
	}

	return &object.StringObj{Value: string(stringObj.Value[index])}
}

// evalArrayLiteral
func evalArrayLiteral(arrayNode *ast.ArrayLiteralNode, env *object.Environment) object.Object {
	var arrayObj = &object.ArrayObj{}

	// evaluate the elements of the array
	if len(arrayNode.Elements) > 0 {
		arrayObj.Elements = evalExpressions(arrayNode.Elements, env)
	}

	return arrayObj
}

// evalHashLiteral
func evalHashLiteral(hashNode *ast.HashLiteralNode, env *object.Environment) object.Object {
	var hashObj = &object.HashObj{}

	// evaluate expressions from keys
	var keysObj = evalExpressions(hashNode.Keys, env)
	// evaluate expressions from values
	var valuesObj = evalExpressions(hashNode.Values, env)

	hashObj.Elements = make(map[string]object.Object)

	// fill the internal hashObj dictionary
	for i := 0; i < len(keysObj); i++ {
		keyStr, ok := keysObj[i].(*object.StringObj)
		if !ok {
			return object.NewError("Invalid data type for key.")
		}
		hashObj.Elements[keyStr.Value] = valuesObj[i]
	}

	return hashObj
}

// evalUnaryExpr
func evalUnaryExpr(unaryNode *ast.UnaryNode, env *object.Environment) object.Object {
	var resultObj object.Object
	var rightObj = Eval(unaryNode.Right, env)

	if isError(rightObj) {
		return rightObj
	}
	// check the token operator
	switch unaryNode.Operator.Type {
	case token.TT_MINUS:
		numberObj, ok := rightObj.(*object.NumberObj)
		if !ok {
			return object.NewError("Invalid data type.")
		}
		resultObj = &object.NumberObj{Value: numberObj.Value * -1}
	case token.TT_NOT:
		booleanObj, ok := rightObj.(*object.BooleanObj)
		if !ok {
			return object.NewError("Invalid data type.")
		}
		resultObj = &object.BooleanObj{Value: !booleanObj.Value}
	default:
		return object.NewError(fmt.Sprintf("Invalid operator: %s", unaryNode.Operator.Literal))
	}
	return resultObj
}

// evalBinaryExpr
func evalBinaryExpr(binOpNode *ast.BinOpNode, env *object.Environment) object.Object {
	if binOpNode.Operator.Type == token.TT_AND || binOpNode.Operator.Type == token.TT_OR {
		return evalLogicalExpression(binOpNode, env)
	}
	// we need to know both operands types
	var leftObj = Eval(binOpNode.Left, env)
	if isError(leftObj) {
		return leftObj
	}
	var rightObj = Eval(binOpNode.Right, env)
	if isError(rightObj) {
		return rightObj
	}
	// now based on their types we perform the correct operation
	if leftObj.Type() == object.NUMBER_OBJ && rightObj.Type() == object.NUMBER_OBJ {
		leftNumber, _ := leftObj.(*object.NumberObj)
		rightNumber, _ := rightObj.(*object.NumberObj)

		return evalBinaryInteger(leftNumber, binOpNode.Operator, rightNumber)
	}
	if leftObj.Type() == object.STRING_OBJ && rightObj.Type() == object.STRING_OBJ {
		leftString, _ := leftObj.(*object.StringObj)
		rightString, _ := rightObj.(*object.StringObj)
		return evalBinaryString(leftString, binOpNode.Operator, rightString)
	}
	if leftObj.Type() == object.BOOLEAN_OBJ && rightObj.Type() == object.BOOLEAN_OBJ {
		leftBoolean, _ := leftObj.(*object.BooleanObj)
		rightBoolean, _ := rightObj.(*object.BooleanObj)
		var (
			leftNumber  = &object.NumberObj{Value: 0}
			rightNumber = &object.NumberObj{Value: 0}
		)
		if leftBoolean.Value {
			leftNumber.Value = 1
		}
		if rightBoolean.Value {
			rightNumber.Value = 1
		}

		return evalBinaryInteger(leftNumber, binOpNode.Operator, rightNumber)
	}
	return object.NewError(fmt.Sprintf("Incompatible data types: %s, %s\n", leftObj.Type(), rightObj.Type()))
}

// evalLogicalExpression
func evalLogicalExpression(binOpNode *ast.BinOpNode, env *object.Environment) object.Object {
	// evaluate the left operator
	var leftObj = Eval(binOpNode.Left, env)
	if isError(leftObj) {
		return leftObj
	}
	leftBinObj, ok := leftObj.(*object.BooleanObj)
	if !ok {
		return object.NewError("Invalid left hand type operand.")
	}
	if binOpNode.Operator.Type == token.TT_AND {
		// if leftObj is false then we do nothing
		if !leftBinObj.Value {
			return FALSE
		}
		// otherwise we need to evaluate the right hand operand
		var rightObj = Eval(binOpNode.Right, env)
		if isError(rightObj) {
			return rightObj
		}
		rightBinOp, ok := rightObj.(*object.BooleanObj)
		if !ok {
			return object.NewError("Invalid right hand type operand.")
		}
		if rightBinOp.Value {
			return TRUE
		} else {
			return FALSE
		}
	} else if binOpNode.Operator.Type == token.TT_OR {
		// if leftObj is true then we do nothing
		if leftBinObj.Value {
			return TRUE
		}
		// otherwise we need to evaluate the right hand operand
		var rightObj = Eval(binOpNode.Right, env)
		if isError(rightObj) {
			return rightObj
		}
		rightBinOp, ok := rightObj.(*object.BooleanObj)
		if !ok {
			return object.NewError("Invalid right hand type operand.")
		}
		if rightBinOp.Value {
			return TRUE
		} else {
			return FALSE
		}
	}
	return object.NewError(fmt.Sprintf("Invalid operator for logical expression: %s\n", binOpNode.Operator.Type))
}

// evalBinaryInteger
func evalBinaryInteger(leftObj *object.NumberObj, ope token.Token, rightObj *object.NumberObj) object.Object {
	var resultObj object.Object
	switch ope.Type {
	case token.TT_PLUS:
		resultObj = &object.NumberObj{Value: float64(leftObj.Value + rightObj.Value)}
	case token.TT_MINUS:
		resultObj = &object.NumberObj{Value: float64(leftObj.Value - rightObj.Value)}
	case token.TT_MUL:
		resultObj = &object.NumberObj{Value: float64(leftObj.Value * rightObj.Value)}
	case token.TT_DIV:
		var rightValue = rightObj.Value
		if rightValue < 0 {
			return object.NewError("Division by zero.")
		}
		resultObj = &object.NumberObj{Value: float64(leftObj.Value / rightObj.Value)}
	case token.TT_LESS:
		if leftObj.Value < rightObj.Value {
			return TRUE
		} else {
			return FALSE
		}
	case token.TT_LESS_EQ:
		if leftObj.Value <= rightObj.Value {
			return TRUE
		} else {
			return FALSE
		}
	case token.TT_GREATER:
		if leftObj.Value > rightObj.Value {
			return TRUE
		} else {
			return FALSE
		}
	case token.TT_GREATER_EQ:
		if leftObj.Value >= rightObj.Value {
			return TRUE
		} else {
			return FALSE
		}
	case token.TT_EQUAL:
		if leftObj.Value == rightObj.Value {
			return TRUE
		} else {
			return FALSE
		}
	case token.TT_NOT_EQ:
		if leftObj.Value != rightObj.Value {
			return TRUE
		} else {
			return FALSE
		}
	default:
		resultObj = object.NewError(fmt.Sprintf("Invalid operator: %s", ope.Literal))
	}

	return resultObj
}

// evalBinaryString
func evalBinaryString(leftObj *object.StringObj, ope token.Token, rightObj *object.StringObj) object.Object {
	var resultObj object.Object
	switch ope.Type {
	case token.TT_PLUS:
		resultObj = &object.StringObj{Value: string(leftObj.Value + rightObj.Value)}
	default:
		resultObj = object.NewError(fmt.Sprintf("Invalid operator: %s\n", ope.Literal))
	}

	return resultObj
}

// evalBooleanNode
func evalBooleanNode(booleanNode *ast.BooleanNode) object.Object {
	if booleanNode.Value {
		return TRUE
	} else {
		return FALSE
	}
}

// evalIdentifier
func evalIdentifier(identifierNode *ast.IdentifierNode, env *object.Environment) object.Object {
	if valueObj, ok := env.Get(identifierNode.Value); ok {
		return valueObj
	}
	// check for builtin function
	if builtinObj, ok := builtins[identifierNode.Value]; ok {
		return builtinObj
	}
	// definetely is not registered.
	return object.NewError(fmt.Sprintf("identifier not found: '%s'", identifierNode.Value))
}

// evalFunctionLiteralNode
func evalFunctionLiteralNode(functionLiteralNode *ast.FunctionLiteralNode, env *object.Environment) object.Object {
	return &object.FunctionObj{
		Parameters: functionLiteralNode.Parameters,
		Body:       functionLiteralNode.Body,
		Env:        env,
	}
}
