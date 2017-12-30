package evaluator

import (
	"dito/src/ast"
	"dito/src/object"
	"dito/src/token"
)

/*
	Statements:
*/

func evalAssignment(node *ast.AssignmentStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	// just ordinary assignement
	if node.Token == token.REASSIGN || node.Token == token.NEWASSIGN {
		env.Set(node.Name.Value, val)
		return nil
	}
	// * Currently don't enforce scope initialation assignment `:=`
	// * as we are only having single statement functions right now.
	// if node.Token == token.REASSIGN {
	// 	if ident := evalIdentifier(node.Name, env); isError(ident) {
	// 		return ident
	// 	}
	// }
	// ######### TODO.
	// this should be changed in the future but for now just implement
	// inplace operators like binary ones. This has just been copy and
	// pasted from evalInfixEpression.
	var newVal object.Object
	operator := string(node.Token.String()[0])
	right := val
	left := evalIdentifier(node.Name, env)
	if isError(left) {
		return left
	}
	switch {
	case isNumericType(left) && isNumericType(right):
		if left.Type() == object.IntergerObj && right.Type() == object.IntergerObj {
			newVal = evalIntegerInfixExpression(operator, left, right)
		} else {
			newVal = evalFloatInfixExpression(operator, left, right)
		}
	case left.Type() != right.Type():
		newVal = newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.StringObj:
		newVal = evalStringExpression(operator, left, right)
	case left.Type() == object.ArrayObj:
		newVal = evalArrayExpression(operator, left, right)
	default:
		newVal = newError("Unknown operator: %s %s= %s", left.Type(), operator, right.Type())
	}
	if isError(newVal) {
		return newVal
	}
	env.Set(node.Name.Value, newVal)
	return nil
}

func evalIndexAssignment(node *ast.IndexAssignmentStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	ident := evalIdentifier(node.IdxExp.Left, env)
	if isError(ident) {
		return ident
	}
	index := Eval(node.IdxExp.Index, env)
	if isError(index) {
		return index
	}
	if ident.Type() != object.ArrayObj || index.Type() != object.IntergerObj {
		return newError("Index operator not supported: %s[%s].", ident.Type(), index.Type())
	}
	arrayObject := ident.(*object.Array)
	idx := index.(*object.Integer).Value
	size := arrayObject.Len - 1
	if idx < 0 {
		idx += (size + 1)
	}
	if idx < 0 || idx > size {
		return newError("Index out of range: len=%d, index=%d.", size, idx)
	}
	arrayObject.Elements[idx] = val
	return nil
}

func evalIfStatement(ie *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTrue(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return object.NONE
	}
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	iterCount := 0
	var body, condition object.Object
	if fs.ID == nil {
		for {
			condition = Eval(fs.Condition, env)
			if !isTrue(condition) {
				break
			}
			if isError(condition) {
				return condition
			}
			if iterCount > 10000 {
				return newError("Max iteration limit reached.")
			}
			iterCount++
			body = Eval(fs.LoopBody, env)
			if body != nil {
				rt := body.Type()
				if rt == object.ErrorObj || rt == object.ReturnObj {
					return body
				}
			}
		}
		return body
	}

	iter := Eval(fs.Iter, env).(*object.Array)
	for _, item := range iter.Elements {
		env.Set(fs.ID.Value, item)
		body = Eval(fs.LoopBody, env)
		if body != nil {
			// the surrounding if is duplicated in isError fn.
			rt := body.Type()
			if rt == object.ErrorObj || rt == object.ReturnObj {
				return body
			}
		}
	}
	return body
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.ErrorObj || rt == object.ReturnObj {
				return result
			}
		}
	}
	return result
}

func isTrue(obj object.Object) bool {
	switch obj {
	case object.NONE:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}
