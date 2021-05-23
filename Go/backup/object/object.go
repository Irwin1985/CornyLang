package object

import (
	"CornyLang/ast"
	"bytes"
	"fmt"
	"strings"
)

type ObjType string

func NewError(message string) *ErrorObj {
	return &ErrorObj{Message: message}
}

const (
	ERROR_OBJ    = "ERROR"
	NUMBER_OBJ   = "NUMBER"
	BOOLEAN_OBJ  = "BOOLEAN"
	STRING_OBJ   = "STRING"
	NULL_OBJ     = "NULL"
	FUNCTION_OBJ = "FUNCTION"
	ARRAY_OBJ    = "ARRAY"
	HASH_OBJ     = "HASH"
	RETURN_OBJ   = "RETURN"
	BUILTIN_OBJ  = "BUILTIN"
)

/**
* Object interface
 */
type Object interface {
	Type() ObjType
	Inspect() string
}

/**
* NumberObj
 */
type NumberObj struct {
	Value float64
}

func (numObj *NumberObj) Type() ObjType {
	return NUMBER_OBJ
}

func (numObj *NumberObj) Inspect() string {
	return fmt.Sprintf("%f", numObj.Value)
}

/**
* ErrorObj
 */
type ErrorObj struct {
	Message string
}

func (errObj *ErrorObj) Type() ObjType {
	return ERROR_OBJ
}

func (errObj *ErrorObj) Inspect() string {
	return errObj.Message
}

/**
* StringObj
 */
type StringObj struct {
	Value string
}

func (stringObj *StringObj) Type() ObjType {
	return STRING_OBJ
}

func (stringObj *StringObj) Inspect() string {
	return string("\"" + stringObj.Value + "\"")
}

/**
* BooleanObj
 */
type BooleanObj struct {
	Value bool
}

func (boolObj *BooleanObj) Type() ObjType {
	return BOOLEAN_OBJ
}

func (boolObj *BooleanObj) Inspect() string {
	if boolObj.Value {
		return "true"
	} else {
		return "false"
	}
}

/**
* NullObj
 */
type NullObj struct {
}

func (nullObj *NullObj) Type() ObjType {
	return NULL_OBJ
}

func (nullObj *NullObj) Inspect() string {
	return "null"
}

/**
* FunctionObj
 */
type FunctionObj struct {
	Parameters []*ast.IdentifierNode
	Body       *ast.BlockStmtNode
	Env        *Environment
}

func (functionObj *FunctionObj) Type() ObjType {
	return FUNCTION_OBJ
}

func (functionObj *FunctionObj) Inspect() string {
	var out bytes.Buffer
	out.WriteString("fn(")
	if functionObj.Parameters != nil {
		params := []string{}
		for _, param := range functionObj.Parameters {
			params = append(params, param.String())
		}
		out.WriteString(strings.Join(params, ", "))
	}
	out.WriteString(")")
	// function Body
	out.WriteString(functionObj.Body.String())
	return out.String()
}

/**
* ArrayObj
 */
type ArrayObj struct {
	Elements []Object
}

func (arrayObj *ArrayObj) Type() ObjType {
	return ARRAY_OBJ
}

func (arrayObj *ArrayObj) Inspect() string {
	var out bytes.Buffer

	out.WriteString("[")

	var elements = []string{}

	for _, elem := range arrayObj.Elements {
		elements = append(elements, elem.Inspect())
	}

	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

/**
* HashObj
 */
type HashObj struct {
	Elements map[string]Object
}

func (hashObj *HashObj) Type() ObjType {
	return HASH_OBJ
}

func (hashObj *HashObj) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")

	var elements = []string{}

	for key, value := range hashObj.Elements {
		elements = append(elements, string(string("\""+key+"\"")+":"+value.Inspect()))
	}

	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

/**
* ReturnObj
 */
type ReturnObj struct {
	Value Object
}

func (returnObj *ReturnObj) Type() ObjType {
	return RETURN_OBJ
}

func (returnObj *ReturnObj) Inspect() string {
	return returnObj.Value.Inspect()
}

/**
* BuiltinObj
 */
type BuiltinFunction func(args ...Object) Object
type BuiltinObj struct {
	Func      BuiltinFunction
	ExtraArgs bool
}

func (builtinObj *BuiltinObj) Type() ObjType {
	return BUILTIN_OBJ
}

func (builtinObj *BuiltinObj) Inspect() string {
	return ""
}
