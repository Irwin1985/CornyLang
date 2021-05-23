package object

import (
	"CornyLang/ast"
	"bytes"
	"fmt"
	"strings"
)

type ObjType string

type ErrorNro int

const (
	VARIABLE_NOT_FOUND = iota
	GLOBAL_ERROR
)

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
	CLASS_OBJ    = "CLASS"
	CALL_OBJ     = "CALL"
)

/**
* Object interface
 */
type Object interface {
	Type() ObjType
	Inspect() string
	HashMethod(method string) bool
}

/**
* ErrorObj
 */
type ErrorObj struct {
	Methods    []string
	Message    string
	ErrorNro   ErrorNro
	Identifier string
}

func (errObj *ErrorObj) Type() ObjType {
	return ERROR_OBJ
}

func (errObj *ErrorObj) Inspect() string {
	return errObj.Message
}

func (errObj *ErrorObj) HashMethod(method string) bool {
	for _, value := range errObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

func NewError(message string, nro ErrorNro, identifier string) *ErrorObj {
	return &ErrorObj{
		Message:    fmt.Sprintf("Error %d: %s", nro, message),
		ErrorNro:   nro,
		Identifier: identifier,
	}
}

/**
* NumberObj
 */
type NumberObj struct {
	Methods []string
	Value   float64
}

func (numObj *NumberObj) Type() ObjType {
	return NUMBER_OBJ
}

func (numObj *NumberObj) Inspect() string {
	return fmt.Sprintf("%f", numObj.Value)
}

func (numObj *NumberObj) HashMethod(method string) bool {
	for _, value := range numObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

func NewNumberObj(value float64) *NumberObj {
	return &NumberObj{
		Methods: []string{
			"type",
		},
		Value: value,
	}
}

/**
* StringObj
 */
type StringObj struct {
	Methods []string
	Value   string
}

func (stringObj *StringObj) Type() ObjType {
	return STRING_OBJ
}

func (stringObj *StringObj) Inspect() string {
	return string("\"" + stringObj.Value + "\"")
}

func (stringObj *StringObj) HashMethod(method string) bool {
	for _, value := range stringObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

func NewStringObj(value string) *StringObj {
	return &StringObj{
		Methods: []string{
			"size",
			"type",
		},
		Value: value,
	}
}

/**
* BooleanObj
 */
type BooleanObj struct {
	Methods []string
	Value   bool
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

func (boolObj *BooleanObj) HashMethod(method string) bool {
	for _, value := range boolObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

func NewBooleanObj(value bool) *BooleanObj {
	return &BooleanObj{
		Methods: []string{
			"type",
		},
		Value: value,
	}
}

/**
* NullObj
 */
type NullObj struct {
	Methods []string
}

func (nullObj *NullObj) Type() ObjType {
	return NULL_OBJ
}

func (nullObj *NullObj) Inspect() string {
	return "null"
}

func (nullObj *NullObj) HashMethod(method string) bool {
	for _, value := range nullObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

func NewNullObj() *NullObj {
	return &NullObj{
		Methods: []string{
			"type",
		},
	}
}

/**
* FunctionObj
 */
type FunctionObj struct {
	Methods    []string
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

func (functionObj *FunctionObj) HashMethod(method string) bool {
	for _, value := range functionObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

func NewFunctionObj() *FunctionObj {
	return &FunctionObj{
		Methods: []string{
			"type",
		},
	}
}

/**
* ArrayObj
 */
type ArrayObj struct {
	Methods  []string
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

func (arrayObj *ArrayObj) HashMethod(method string) bool {
	for _, value := range arrayObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

func NewArrayObj() *ArrayObj {
	return &ArrayObj{
		Methods: []string{
			"type",
			"size",
		},
	}
}

/**
* HashObj
 */
type HashObj struct {
	Methods  []string
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

func (hashObj *HashObj) HashMethod(method string) bool {
	for _, value := range hashObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

func NewHashObj() *HashObj {
	return &HashObj{
		Methods: []string{
			"type",
			"size",
		},
	}
}

/**
* ReturnObj
 */
type ReturnObj struct {
	Methods []string
	Value   Object
}

func (returnObj *ReturnObj) Type() ObjType {
	return RETURN_OBJ
}

func (returnObj *ReturnObj) Inspect() string {
	return returnObj.Value.Inspect()
}

func (returnObj *ReturnObj) HashMethod(method string) bool {
	for _, value := range returnObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

/**
* BuiltinObj
 */
type BuiltinFunction func(args ...Object) Object
type BuiltinObj struct {
	Methods   []string
	Func      BuiltinFunction
	ExtraArgs bool
}

func (builtinObj *BuiltinObj) Type() ObjType {
	return BUILTIN_OBJ
}

func (builtinObj *BuiltinObj) Inspect() string {
	return ""
}

func (builtinObj *BuiltinObj) HashMethod(method string) bool {
	for _, value := range builtinObj.Methods {
		if value == method {
			return true
		}
	}
	return false
}

/**
* ClassObj
 */
type ClassObj struct {
	Methods []string
	Env     *Environment
}

func (classObj *ClassObj) Type() ObjType {
	return CLASS_OBJ
}

func (classObj *ClassObj) Inspect() string {
	return ""
}

func (classObj *ClassObj) HashMethod(method string) bool {
	for _, value := range classObj.Methods {
		if value == method {
			return true
		}
	}
	// try find the method in the internal dictionary (Pems)
	if _, ok := classObj.Env.Get(method); ok {
		return true
	}
	// no method found
	return false
}

func NewClassObj() *ClassObj {
	return &ClassObj{
		Methods: []string{
			"type",
		},
	}
}

/**
* CallObj
 */
type CallObj struct {
	Methods   []string
	Callee    string
	Arguments []Object
}

func (callObj *CallObj) Type() ObjType {
	return CALL_OBJ
}

func (callObj *CallObj) Inspect() string {
	return "call:ok"
}

func (callObj *CallObj) HashMethod(method string) bool {
	for _, value := range callObj.Methods {
		if value == method {
			return true
		}
	}
	// no method found
	return false
}

func NewCallObj() *CallObj {
	return &CallObj{
		Methods: []string{
			"type",
		},
	}
}
