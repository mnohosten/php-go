package ast

// Visitor is an interface for traversing the AST using the visitor pattern
// Each Visit method receives a node and returns a boolean indicating whether
// to continue traversing child nodes (true) or skip them (false)
type Visitor interface {
	// Statement visitors
	VisitExpressionStatement(node *ExpressionStatement) bool
	VisitBlockStatement(node *BlockStatement) bool
	VisitEchoStatement(node *EchoStatement) bool
	VisitReturnStatement(node *ReturnStatement) bool
	VisitBreakStatement(node *BreakStatement) bool
	VisitContinueStatement(node *ContinueStatement) bool
	VisitIfStatement(node *IfStatement) bool
	VisitWhileStatement(node *WhileStatement) bool
	VisitDoWhileStatement(node *DoWhileStatement) bool
	VisitForStatement(node *ForStatement) bool
	VisitForeachStatement(node *ForeachStatement) bool
	VisitSwitchStatement(node *SwitchStatement) bool
	VisitTryStatement(node *TryStatement) bool
	VisitThrowStatement(node *ThrowStatement) bool
	VisitFunctionDeclaration(node *FunctionDeclaration) bool
	VisitClassDeclaration(node *ClassDeclaration) bool
	VisitInterfaceDeclaration(node *InterfaceDeclaration) bool
	VisitTraitDeclaration(node *TraitDeclaration) bool
	VisitPropertyDeclaration(node *PropertyDeclaration) bool
	VisitMethodDeclaration(node *MethodDeclaration) bool
	VisitClassConstantDeclaration(node *ClassConstantDeclaration) bool
	VisitTraitUse(node *TraitUse) bool

	// Expression visitors
	VisitIdentifier(node *Identifier) bool
	VisitIntegerLiteral(node *IntegerLiteral) bool
	VisitFloatLiteral(node *FloatLiteral) bool
	VisitStringLiteral(node *StringLiteral) bool
	VisitBooleanLiteral(node *BooleanLiteral) bool
	VisitNullLiteral(node *NullLiteral) bool
	VisitVariable(node *Variable) bool
	VisitArrayExpression(node *ArrayExpression) bool
	VisitPrefixExpression(node *PrefixExpression) bool
	VisitInfixExpression(node *InfixExpression) bool
	VisitAssignmentExpression(node *AssignmentExpression) bool
	VisitTernaryExpression(node *TernaryExpression) bool
	VisitIndexExpression(node *IndexExpression) bool
	VisitPropertyExpression(node *PropertyExpression) bool
	VisitNullsafePropertyExpression(node *NullsafePropertyExpression) bool
	VisitStaticPropertyExpression(node *StaticPropertyExpression) bool
	VisitCallExpression(node *CallExpression) bool
	VisitMethodCallExpression(node *MethodCallExpression) bool
	VisitStaticCallExpression(node *StaticCallExpression) bool
	VisitNewExpression(node *NewExpression) bool
	VisitInstanceofExpression(node *InstanceofExpression) bool
	VisitCastExpression(node *CastExpression) bool
	VisitGroupedExpression(node *GroupedExpression) bool
	VisitMatchExpression(node *MatchExpression) bool
	VisitNullableType(node *NullableType) bool
	VisitUnionType(node *UnionType) bool
	VisitIntersectionType(node *IntersectionType) bool
}

// Walk traverses the AST starting from the given node using the visitor pattern
func Walk(v Visitor, node Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	// Statements
	case *ExpressionStatement:
		if v.VisitExpressionStatement(n) {
			Walk(v, n.Expression)
		}
	case *BlockStatement:
		if v.VisitBlockStatement(n) {
			for _, stmt := range n.Statements {
				Walk(v, stmt)
			}
		}
	case *EchoStatement:
		if v.VisitEchoStatement(n) {
			for _, expr := range n.Expressions {
				Walk(v, expr)
			}
		}
	case *ReturnStatement:
		if v.VisitReturnStatement(n) {
			Walk(v, n.ReturnValue)
		}
	case *BreakStatement:
		if v.VisitBreakStatement(n) {
			Walk(v, n.Depth)
		}
	case *ContinueStatement:
		if v.VisitContinueStatement(n) {
			Walk(v, n.Depth)
		}
	case *IfStatement:
		if v.VisitIfStatement(n) {
			Walk(v, n.Condition)
			Walk(v, n.Consequence)
			for _, elseif := range n.ElseIfs {
				Walk(v, elseif.Condition)
				Walk(v, elseif.Consequence)
			}
			Walk(v, n.Alternative)
		}
	case *WhileStatement:
		if v.VisitWhileStatement(n) {
			Walk(v, n.Condition)
			Walk(v, n.Body)
		}
	case *DoWhileStatement:
		if v.VisitDoWhileStatement(n) {
			Walk(v, n.Body)
			Walk(v, n.Condition)
		}
	case *ForStatement:
		if v.VisitForStatement(n) {
			for _, expr := range n.Init {
				Walk(v, expr)
			}
			for _, expr := range n.Condition {
				Walk(v, expr)
			}
			for _, expr := range n.Increment {
				Walk(v, expr)
			}
			Walk(v, n.Body)
		}
	case *ForeachStatement:
		if v.VisitForeachStatement(n) {
			Walk(v, n.Array)
			Walk(v, n.Key)
			Walk(v, n.Value)
			Walk(v, n.Body)
		}
	case *SwitchStatement:
		if v.VisitSwitchStatement(n) {
			Walk(v, n.Subject)
			for _, c := range n.Cases {
				Walk(v, c.Value)
				for _, stmt := range c.Body {
					Walk(v, stmt)
				}
			}
		}
	case *TryStatement:
		if v.VisitTryStatement(n) {
			Walk(v, n.Body)
			for _, c := range n.CatchClauses {
				for _, t := range c.Types {
					Walk(v, t)
				}
				Walk(v, c.Variable)
				Walk(v, c.Body)
			}
			Walk(v, n.Finally)
		}
	case *ThrowStatement:
		if v.VisitThrowStatement(n) {
			Walk(v, n.Expression)
		}
	case *FunctionDeclaration:
		if v.VisitFunctionDeclaration(n) {
			Walk(v, n.Name)
			for _, p := range n.Parameters {
				Walk(v, p.Type)
				Walk(v, p.Name)
				Walk(v, p.DefaultValue)
			}
			Walk(v, n.ReturnType)
			Walk(v, n.Body)
		}
	case *ClassDeclaration:
		if v.VisitClassDeclaration(n) {
			Walk(v, n.Name)
			Walk(v, n.Extends)
			for _, i := range n.Implements {
				Walk(v, i)
			}
			for _, stmt := range n.Body {
				Walk(v, stmt)
			}
		}
	case *InterfaceDeclaration:
		if v.VisitInterfaceDeclaration(n) {
			Walk(v, n.Name)
			for _, e := range n.Extends {
				Walk(v, e)
			}
			for _, m := range n.Body {
				Walk(v, m.Name)
				for _, p := range m.Parameters {
					Walk(v, p.Type)
					Walk(v, p.Name)
					Walk(v, p.DefaultValue)
				}
				Walk(v, m.ReturnType)
			}
		}
	case *TraitDeclaration:
		if v.VisitTraitDeclaration(n) {
			Walk(v, n.Name)
			for _, stmt := range n.Body {
				Walk(v, stmt)
			}
		}
	case *PropertyDeclaration:
		if v.VisitPropertyDeclaration(n) {
			Walk(v, n.Type)
			for _, p := range n.Properties {
				Walk(v, p.Name)
				Walk(v, p.DefaultValue)
			}
		}
	case *MethodDeclaration:
		if v.VisitMethodDeclaration(n) {
			Walk(v, n.Name)
			for _, p := range n.Parameters {
				Walk(v, p.Type)
				Walk(v, p.Name)
				Walk(v, p.DefaultValue)
			}
			Walk(v, n.ReturnType)
			Walk(v, n.Body)
		}
	case *ClassConstantDeclaration:
		if v.VisitClassConstantDeclaration(n) {
			for _, c := range n.Constants {
				Walk(v, c.Name)
				Walk(v, c.Value)
			}
		}
	case *TraitUse:
		if v.VisitTraitUse(n) {
			for _, t := range n.Traits {
				Walk(v, t)
			}
		}

	// Expressions
	case *Identifier:
		v.VisitIdentifier(n)
	case *IntegerLiteral:
		v.VisitIntegerLiteral(n)
	case *FloatLiteral:
		v.VisitFloatLiteral(n)
	case *StringLiteral:
		v.VisitStringLiteral(n)
	case *BooleanLiteral:
		v.VisitBooleanLiteral(n)
	case *NullLiteral:
		v.VisitNullLiteral(n)
	case *Variable:
		v.VisitVariable(n)
	case *ArrayExpression:
		if v.VisitArrayExpression(n) {
			for _, elem := range n.Elements {
				Walk(v, elem.Key)
				Walk(v, elem.Value)
			}
		}
	case *PrefixExpression:
		if v.VisitPrefixExpression(n) {
			Walk(v, n.Right)
		}
	case *InfixExpression:
		if v.VisitInfixExpression(n) {
			Walk(v, n.Left)
			Walk(v, n.Right)
		}
	case *AssignmentExpression:
		if v.VisitAssignmentExpression(n) {
			Walk(v, n.Left)
			Walk(v, n.Right)
		}
	case *TernaryExpression:
		if v.VisitTernaryExpression(n) {
			Walk(v, n.Condition)
			Walk(v, n.Consequence)
			Walk(v, n.Alternative)
		}
	case *IndexExpression:
		if v.VisitIndexExpression(n) {
			Walk(v, n.Left)
			Walk(v, n.Index)
		}
	case *PropertyExpression:
		if v.VisitPropertyExpression(n) {
			Walk(v, n.Object)
			Walk(v, n.Property)
		}
	case *NullsafePropertyExpression:
		if v.VisitNullsafePropertyExpression(n) {
			Walk(v, n.Object)
			Walk(v, n.Property)
		}
	case *StaticPropertyExpression:
		if v.VisitStaticPropertyExpression(n) {
			Walk(v, n.Class)
			Walk(v, n.Property)
		}
	case *CallExpression:
		if v.VisitCallExpression(n) {
			Walk(v, n.Function)
			for _, arg := range n.Arguments {
				Walk(v, arg)
			}
		}
	case *MethodCallExpression:
		if v.VisitMethodCallExpression(n) {
			Walk(v, n.Object)
			Walk(v, n.Method)
			for _, arg := range n.Arguments {
				Walk(v, arg)
			}
		}
	case *StaticCallExpression:
		if v.VisitStaticCallExpression(n) {
			Walk(v, n.Class)
			Walk(v, n.Method)
			for _, arg := range n.Arguments {
				Walk(v, arg)
			}
		}
	case *NewExpression:
		if v.VisitNewExpression(n) {
			Walk(v, n.Class)
			for _, arg := range n.Arguments {
				Walk(v, arg)
			}
		}
	case *InstanceofExpression:
		if v.VisitInstanceofExpression(n) {
			Walk(v, n.Left)
			Walk(v, n.Right)
		}
	case *CastExpression:
		if v.VisitCastExpression(n) {
			Walk(v, n.Expr)
		}
	case *GroupedExpression:
		if v.VisitGroupedExpression(n) {
			Walk(v, n.Expr)
		}
	case *MatchExpression:
		if v.VisitMatchExpression(n) {
			Walk(v, n.Subject)
			for _, arm := range n.Arms {
				for _, cond := range arm.Conditions {
					Walk(v, cond)
				}
				Walk(v, arm.Body)
			}
		}
	case *NullableType:
		if v.VisitNullableType(n) {
			Walk(v, n.Type)
		}
	case *UnionType:
		if v.VisitUnionType(n) {
			for _, t := range n.Types {
				Walk(v, t)
			}
		}
	case *IntersectionType:
		if v.VisitIntersectionType(n) {
			for _, t := range n.Types {
				Walk(v, t)
			}
		}
	}
}

// BaseVisitor provides default implementations for all visitor methods
// Embed this in your visitor to only override the methods you need
type BaseVisitor struct{}

func (bv *BaseVisitor) VisitExpressionStatement(node *ExpressionStatement) bool       { return true }
func (bv *BaseVisitor) VisitBlockStatement(node *BlockStatement) bool                 { return true }
func (bv *BaseVisitor) VisitEchoStatement(node *EchoStatement) bool                   { return true }
func (bv *BaseVisitor) VisitReturnStatement(node *ReturnStatement) bool               { return true }
func (bv *BaseVisitor) VisitBreakStatement(node *BreakStatement) bool                 { return true }
func (bv *BaseVisitor) VisitContinueStatement(node *ContinueStatement) bool           { return true }
func (bv *BaseVisitor) VisitIfStatement(node *IfStatement) bool                       { return true }
func (bv *BaseVisitor) VisitWhileStatement(node *WhileStatement) bool                 { return true }
func (bv *BaseVisitor) VisitDoWhileStatement(node *DoWhileStatement) bool             { return true }
func (bv *BaseVisitor) VisitForStatement(node *ForStatement) bool                     { return true }
func (bv *BaseVisitor) VisitForeachStatement(node *ForeachStatement) bool             { return true }
func (bv *BaseVisitor) VisitSwitchStatement(node *SwitchStatement) bool               { return true }
func (bv *BaseVisitor) VisitTryStatement(node *TryStatement) bool                     { return true }
func (bv *BaseVisitor) VisitThrowStatement(node *ThrowStatement) bool                 { return true }
func (bv *BaseVisitor) VisitFunctionDeclaration(node *FunctionDeclaration) bool       { return true }
func (bv *BaseVisitor) VisitClassDeclaration(node *ClassDeclaration) bool             { return true }
func (bv *BaseVisitor) VisitInterfaceDeclaration(node *InterfaceDeclaration) bool     { return true }
func (bv *BaseVisitor) VisitTraitDeclaration(node *TraitDeclaration) bool             { return true }
func (bv *BaseVisitor) VisitPropertyDeclaration(node *PropertyDeclaration) bool       { return true }
func (bv *BaseVisitor) VisitMethodDeclaration(node *MethodDeclaration) bool           { return true }
func (bv *BaseVisitor) VisitClassConstantDeclaration(node *ClassConstantDeclaration) bool {
	return true
}
func (bv *BaseVisitor) VisitTraitUse(node *TraitUse) bool                             { return true }
func (bv *BaseVisitor) VisitIdentifier(node *Identifier) bool                         { return true }
func (bv *BaseVisitor) VisitIntegerLiteral(node *IntegerLiteral) bool                 { return true }
func (bv *BaseVisitor) VisitFloatLiteral(node *FloatLiteral) bool                     { return true }
func (bv *BaseVisitor) VisitStringLiteral(node *StringLiteral) bool                   { return true }
func (bv *BaseVisitor) VisitBooleanLiteral(node *BooleanLiteral) bool                 { return true }
func (bv *BaseVisitor) VisitNullLiteral(node *NullLiteral) bool                       { return true }
func (bv *BaseVisitor) VisitVariable(node *Variable) bool                             { return true }
func (bv *BaseVisitor) VisitArrayExpression(node *ArrayExpression) bool               { return true }
func (bv *BaseVisitor) VisitPrefixExpression(node *PrefixExpression) bool             { return true }
func (bv *BaseVisitor) VisitInfixExpression(node *InfixExpression) bool               { return true }
func (bv *BaseVisitor) VisitAssignmentExpression(node *AssignmentExpression) bool     { return true }
func (bv *BaseVisitor) VisitTernaryExpression(node *TernaryExpression) bool           { return true }
func (bv *BaseVisitor) VisitIndexExpression(node *IndexExpression) bool               { return true }
func (bv *BaseVisitor) VisitPropertyExpression(node *PropertyExpression) bool         { return true }
func (bv *BaseVisitor) VisitNullsafePropertyExpression(node *NullsafePropertyExpression) bool {
	return true
}
func (bv *BaseVisitor) VisitStaticPropertyExpression(node *StaticPropertyExpression) bool {
	return true
}
func (bv *BaseVisitor) VisitCallExpression(node *CallExpression) bool             { return true }
func (bv *BaseVisitor) VisitMethodCallExpression(node *MethodCallExpression) bool { return true }
func (bv *BaseVisitor) VisitStaticCallExpression(node *StaticCallExpression) bool { return true }
func (bv *BaseVisitor) VisitNewExpression(node *NewExpression) bool               { return true }
func (bv *BaseVisitor) VisitInstanceofExpression(node *InstanceofExpression) bool { return true }
func (bv *BaseVisitor) VisitCastExpression(node *CastExpression) bool             { return true }
func (bv *BaseVisitor) VisitGroupedExpression(node *GroupedExpression) bool       { return true }
func (bv *BaseVisitor) VisitMatchExpression(node *MatchExpression) bool           { return true }
func (bv *BaseVisitor) VisitNullableType(node *NullableType) bool                 { return true }
func (bv *BaseVisitor) VisitUnionType(node *UnionType) bool                       { return true }
func (bv *BaseVisitor) VisitIntersectionType(node *IntersectionType) bool         { return true }
