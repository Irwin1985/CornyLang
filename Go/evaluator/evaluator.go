package evaluator

import (
	"CornyLang/ast"
	"CornyLang/object"
	"CornyLang/token"
	"fmt"
)

/**
* Variables
 */
var TRUE = object.NewBooleanObj(true)
var FALSE = object.NewBooleanObj(false)
var NULL = object.NewNullObj()

/**
* Helpers functions
 */
func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

/**
* Eval Expressions
 */
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

/**
* Eval Method
 */
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.ProgramNode:
		return evalProgram(node, env)
	case *ast.BlockStmtNode:
		return evalBlock(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.NumberNode:
		return object.NewNumberObj(node.Value)
	case *ast.StringNode:
		return object.NewStringObj(node.Value)
	case *ast.BooleanNode:
		return evalBooleanNode(node)
	case *ast.NullNode:
		return NULL
	case *ast.FunctionLiteralNode:
		return evalFunctionLiteral(node, env)
	case *ast.ClassLiteralNode:
		return evalClassLiteral(node, env)
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

/**
* Eval Program
 */
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

/**
* Eval Block
 */
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

/**
* Eval Let
 */
func evalLet(letNode *ast.LetNode, env *object.Environment) object.Object {
	var valueObj = Eval(letNode.Value, env)
	if isError(valueObj) {
		return valueObj
	}
	// check for reserved word
	if _, ok := builtins[letNode.Name.Value]; ok {
		return object.NewError(fmt.Sprintf("'%s' its a reserved word.", letNode.Name.Value), object.GLOBAL_ERROR, "")
	}
	// register in symbol table
	env.Set(letNode.Name.Value, valueObj)

	return valueObj
}

/**
* Eval If Expression
 */
func evalIfExpr(ifNode *ast.IfExprNode, env *object.Environment) object.Object {
	var conditionObj = Eval(ifNode.Condition, env)
	if isError(conditionObj) {
		return conditionObj
	}
	boolConditionObj, ok := conditionObj.(*object.BooleanObj)
	if !ok {
		return object.NewError("Invalid data type for if condition.", object.GLOBAL_ERROR, "")
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

/**
* Eval Call Expression
 */
func evalCallExpr(callExpr *ast.CallExprNode, env *object.Environment) object.Object {
	// eval the callee object
	var calleeObj = Eval(callExpr.Callee, env)
	if isError(calleeObj) {
		if errorObj, ok := calleeObj.(*object.ErrorObj); ok {
			if errorObj.ErrorNro == object.VARIABLE_NOT_FOUND {
				var callObj = object.NewCallObj()
				callObj.Arguments = evalExpressions(callExpr.Arguments, env)
				callObj.Callee = errorObj.Identifier
				return callObj
			}
			return calleeObj
		}
	}

	// eval arguments
	var args = evalExpressions(callExpr.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
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
		}
	}
	return calleeObj
}

/**
* Eval Function Execution
 */
func evalFunction(functionObj *object.FunctionObj, args []object.Object) object.Object {
	// 1. create new environment and close it with FunctionObj.Env
	var newEnv = object.NewEnvironment()
	newEnv.Parent = functionObj.Env

	// 2. check for function arity
	var numArgs = len(args)
	var numParams = len(functionObj.Parameters)
	if numArgs != numParams {
		return object.NewError(fmt.Sprintf("Unexpected arguments, got: %d, want: %d", numArgs, numParams), object.GLOBAL_ERROR, "")
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

/**
* Eval Array Access
 */
func evalArrayAccess(arrayObj *object.ArrayObj, args []object.Object) object.Object {
	indexObj, ok := args[0].(*object.NumberObj)
	if !ok {
		return object.NewError("Invalid subscript reference.", object.GLOBAL_ERROR, "")
	}
	// check for out of bounds
	var index = int(indexObj.Value)
	if index < 0 || index >= len(arrayObj.Elements) {
		return object.NewError("Index out of bounds.", object.GLOBAL_ERROR, "")
	}

	return arrayObj.Elements[index]
}

/**
* Eval Hash Access
 */
func evalHashAccess(hashObj *object.HashObj, args []object.Object) object.Object {
	indexObj, ok := args[0].(*object.StringObj)
	if !ok {
		return object.NewError("Invalid subscript reference", object.GLOBAL_ERROR, "")
	}
	var key = indexObj.Value
	valueObj, ok := hashObj.Elements[key]
	if !ok {
		return NULL
	}
	return valueObj
}

/**
* Eval String Access
 */
func evalStringAccess(stringObj *object.StringObj, args []object.Object) object.Object {
	indexObj, ok := args[0].(*object.NumberObj)
	if !ok {
		return object.NewError("Invalid subscript reference", object.GLOBAL_ERROR, "")
	}
	var index = int(indexObj.Value)
	if index < 0 || index >= len(stringObj.Value) {
		return object.NewError("Index out of bounds", object.GLOBAL_ERROR, "")
	}
	return object.NewStringObj(string(stringObj.Value[index]))
	//return &object.StringObj{Value: string(stringObj.Value[index])}
}

/**
* Eval Function Literal
 */
func evalFunctionLiteral(functionLiteralNode *ast.FunctionLiteralNode, env *object.Environment) object.Object {
	var functionObj = object.NewFunctionObj()

	functionObj.Parameters = functionLiteralNode.Parameters
	functionObj.Body = functionLiteralNode.Body
	functionObj.Env = env

	return functionObj
}

/**
* Eval Class Literal
 */
func evalClassLiteral(classLiteralNode *ast.ClassLiteralNode, env *object.Environment) object.Object {
	var classObj = object.NewClassObj()
	// create new environment
	var newEnv = object.NewEnvironment()
	newEnv.Parent = env
	// fill environment
	for _, stmt := range classLiteralNode.Body.Statements {
		if letNode, ok := stmt.(*ast.LetNode); ok {
			resultObj := Eval(letNode, newEnv)
			if isError(resultObj) {
				return resultObj
			}
		} else {
			return object.NewError("Invalid statement in class definition.", object.GLOBAL_ERROR, "")
		}
	}
	// add native methods
	newEnv.Set("type", builtins["type"])
	classObj.Env = newEnv
	return classObj
}

/**
* Eval Array Literal
 */
func evalArrayLiteral(arrayNode *ast.ArrayLiteralNode, env *object.Environment) object.Object {
	var arrayObj = object.NewArrayObj()

	// evaluate the elements of the array
	if len(arrayNode.Elements) > 0 {
		arrayObj.Elements = evalExpressions(arrayNode.Elements, env)
	}

	return arrayObj
}

/**
* Eval Hash Literal
 */
func evalHashLiteral(hashNode *ast.HashLiteralNode, env *object.Environment) object.Object {
	var hashObj = object.NewHashObj()

	// evaluate expressions from keys
	var keysObj = evalExpressions(hashNode.Keys, env)
	// evaluate expressions from values
	var valuesObj = evalExpressions(hashNode.Values, env)

	hashObj.Elements = make(map[string]object.Object)

	// fill the internal hashObj dictionary
	for i := 0; i < len(keysObj); i++ {
		keyStr, ok := keysObj[i].(*object.StringObj)
		if !ok {
			return object.NewError("Invalid data type for key.", object.GLOBAL_ERROR, "")
		}
		hashObj.Elements[keyStr.Value] = valuesObj[i]
	}

	return hashObj
}

/**
* Eval Unary Expression
 */
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
			return object.NewError("Invalid data type.", object.GLOBAL_ERROR, "")
		}
		resultObj = object.NewNumberObj(numberObj.Value * -1)
	case token.TT_NOT:
		booleanObj, ok := rightObj.(*object.BooleanObj)
		if !ok {
			return object.NewError("Invalid data type.", object.GLOBAL_ERROR, "")
		}
		resultObj = &object.BooleanObj{Value: !booleanObj.Value}
	default:
		return object.NewError(fmt.Sprintf("Invalid operator: %s", unaryNode.Operator.Literal), object.GLOBAL_ERROR, "")
	}
	return resultObj
}

/**
* Eval Binary Expression
 */
func evalBinaryExpr(binOpNode *ast.BinOpNode, env *object.Environment) object.Object {
	// check for logical operation
	if binOpNode.Operator.Type == token.TT_AND || binOpNode.Operator.Type == token.TT_OR {
		return evalLogicalExpression(binOpNode, env)
	}
	// check for dot operation
	if binOpNode.Operator.Type == token.TT_DOT {
		return evalDotExpression(binOpNode, env)
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
			//leftNumber  = &object.NumberObj{Value: 0}
			leftNumber  = object.NewNumberObj(0)
			rightNumber = object.NewNumberObj(0)
			//rightNumber = &object.NumberObj{Value: 0}
		)
		if leftBoolean.Value {
			leftNumber.Value = 1
		}
		if rightBoolean.Value {
			rightNumber.Value = 1
		}

		return evalBinaryInteger(leftNumber, binOpNode.Operator, rightNumber)
	}
	return object.NewError(fmt.Sprintf("Incompatible data types: %s, %s", leftObj.Type(), rightObj.Type()), object.GLOBAL_ERROR, "")
}

/**
* Eval Logical Expression
 */
func evalLogicalExpression(binOpNode *ast.BinOpNode, env *object.Environment) object.Object {
	// evaluate the left operator
	var leftObj = Eval(binOpNode.Left, env)
	if isError(leftObj) {
		return leftObj
	}
	leftBinObj, ok := leftObj.(*object.BooleanObj)
	if !ok {
		return object.NewError("Invalid left hand type operand.", object.GLOBAL_ERROR, "")
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
			return object.NewError("Invalid right hand type operand.", object.GLOBAL_ERROR, "")
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
			return object.NewError("Invalid right hand type operand.", object.GLOBAL_ERROR, "")
		}
		if rightBinOp.Value {
			return TRUE
		} else {
			return FALSE
		}
	}
	return object.NewError(fmt.Sprintf("Invalid operator for logical expression: %s", binOpNode.Operator.Type), object.GLOBAL_ERROR, "")
}

/**
* Eval Dot accessor Expression
 */
func evalDotExpression(binOpNode *ast.BinOpNode, env *object.Environment) object.Object {
	var msg string
	// 1. Resolve left hand operand
	var leftObj = Eval(binOpNode.Left, env)
	if isError(leftObj) {
		return leftObj
	}
	// 3. Resolve the right hand operand
	var rightObj object.Object
	if classObj, ok := leftObj.(*object.ClassObj); ok {
		rightObj = Eval(binOpNode.Right, classObj.Env)
	} else {
		rightObj = Eval(binOpNode.Right, env)
	}
	if isError(rightObj) {
		return rightObj
	}
	var methodName string
	var newArgs []object.Object
	// check for rightObj type.
	if callObj, ok := rightObj.(*object.CallObj); ok {
		methodName = callObj.Callee
		if callObj.Arguments != nil {
			newArgs = callObj.Arguments
		}
	} else {
		return rightObj
	}
	// 5. ask for method to respond
	if leftObj.HashMethod(methodName) {
		// get the method
		if classObj, ok := leftObj.(*object.ClassObj); ok {
			methodObj, _ := classObj.Env.Get(methodName)
			if funcObj, ok := methodObj.(*object.FunctionObj); ok {
				return evalFunction(funcObj, newArgs)
			} else if builtinObj, ok := methodObj.(*object.BuiltinObj); ok {
				return builtinObj.Func([]object.Object{leftObj}...)
			}
			return methodObj
		} else {
			if builtinObj, ok := builtins[methodName]; ok {
				return builtinObj.Func([]object.Object{leftObj}...)
			}
		}
	} else {
		msg = fmt.Sprintf("Method '%s' not found.", methodName)
	}
	return object.NewError(msg, object.GLOBAL_ERROR, "")
}

/**
* Eval Binary Integer
 */
func evalBinaryInteger(leftObj *object.NumberObj, ope token.Token, rightObj *object.NumberObj) object.Object {
	var resultObj object.Object
	switch ope.Type {
	case token.TT_PLUS:
		//resultObj = &object.NumberObj{Value: float64(leftObj.Value + rightObj.Value)}
		resultObj = object.NewNumberObj(leftObj.Value + rightObj.Value)
	case token.TT_MINUS:
		//resultObj = &object.NumberObj{Value: float64(leftObj.Value - rightObj.Value)}
		resultObj = object.NewNumberObj(leftObj.Value - rightObj.Value)
	case token.TT_MUL:
		//resultObj = &object.NumberObj{Value: float64(leftObj.Value * rightObj.Value)}
		resultObj = object.NewNumberObj(leftObj.Value * rightObj.Value)
	case token.TT_DIV:
		var rightValue = rightObj.Value
		if rightValue < 0 {
			return object.NewError("Division by zero.", object.GLOBAL_ERROR, "")
		}
		//resultObj = &object.NumberObj{Value: float64(leftObj.Value / rightObj.Value)}
		resultObj = object.NewNumberObj(leftObj.Value / rightObj.Value)
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
		resultObj = object.NewError(fmt.Sprintf("Invalid operator: %s", ope.Literal), object.GLOBAL_ERROR, "")
	}

	return resultObj
}

/**
* Eval Binary String
 */
func evalBinaryString(leftObj *object.StringObj, ope token.Token, rightObj *object.StringObj) object.Object {
	var resultObj object.Object
	switch ope.Type {
	case token.TT_PLUS:
		resultObj = object.NewStringObj(string(leftObj.Value + rightObj.Value))
		//resultObj = &object.StringObj{Value: string(leftObj.Value + rightObj.Value)}
	default:
		resultObj = object.NewError(fmt.Sprintf("Invalid operator: %s", ope.Literal), object.GLOBAL_ERROR, "")
	}

	return resultObj
}

/**
* Eval Boolean Node
 */
func evalBooleanNode(booleanNode *ast.BooleanNode) object.Object {
	if booleanNode.Value {
		return TRUE
	} else {
		return FALSE
	}
}

/**
* Eval Eval Identifier
 */
func evalIdentifier(identifierNode *ast.IdentifierNode, env *object.Environment) object.Object {
	if valueObj, ok := env.Get(identifierNode.Value); ok {
		return valueObj
	}
	// definetely is not registered.
	return object.NewError(fmt.Sprintf("identifier not found: '%s'", identifierNode.Value), object.VARIABLE_NOT_FOUND, identifierNode.Value)
}
