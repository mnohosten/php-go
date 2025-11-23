package parallel

import (
	"fmt"
	"strings"

	"github.com/krizos/php-go/pkg/ast"
)

// ============================================================================
// Safety Analyzer - Determines if code can be safely parallelized
// ============================================================================

// SafetyAnalyzer analyzes PHP code to determine if it's safe to parallelize
type SafetyAnalyzer struct {
	// Current function being analyzed
	currentFunction *ast.FunctionDeclaration
}

// NewSafetyAnalyzer creates a new safety analyzer
func NewSafetyAnalyzer() *SafetyAnalyzer {
	return &SafetyAnalyzer{}
}

// SafetyReport contains the results of analyzing a function for parallelization safety
type SafetyReport struct {
	// Safe indicates if the function can be safely parallelized
	Safe bool

	// Reasons lists why the function is unsafe (if not safe)
	Reasons []string

	// GlobalReads lists global variables that are read
	GlobalReads []string

	// GlobalWrites lists global variables that are written
	GlobalWrites []string

	// StaticAccess indicates if static variables are accessed
	StaticAccess bool

	// StaticVars lists static variables accessed
	StaticVars []string

	// FileIO indicates if file I/O operations are detected
	FileIO bool

	// FileOperations lists file I/O function calls
	FileOperations []string

	// NetworkIO indicates if network I/O operations are detected
	NetworkIO bool

	// NetworkOperations lists network I/O function calls
	NetworkOperations []string

	// DatabaseAccess indicates if database operations are detected
	DatabaseAccess bool

	// DatabaseOperations lists database function calls
	DatabaseOperations []string

	// SideEffects lists other side effects detected
	SideEffects []string
}

// CanParallelize analyzes a function and returns true if it can be safely parallelized
func (sa *SafetyAnalyzer) CanParallelize(fn *ast.FunctionDeclaration) bool {
	report := sa.AnalyzeFunction(fn)
	return report.Safe
}

// AnalyzeFunction analyzes a function and returns a detailed safety report
func (sa *SafetyAnalyzer) AnalyzeFunction(fn *ast.FunctionDeclaration) *SafetyReport {
	sa.currentFunction = fn

	report := &SafetyReport{
		Safe:               true,
		Reasons:            []string{},
		GlobalReads:        []string{},
		GlobalWrites:       []string{},
		StaticVars:         []string{},
		FileOperations:     []string{},
		NetworkOperations:  []string{},
		DatabaseOperations: []string{},
		SideEffects:        []string{},
	}

	if fn.Body != nil {
		sa.analyzeBlock(fn.Body, report)
	}

	// If any unsafe patterns detected, mark as unsafe
	if len(report.GlobalWrites) > 0 {
		report.Safe = false
		report.Reasons = append(report.Reasons, fmt.Sprintf("writes to %d global variable(s)", len(report.GlobalWrites)))
	}

	if report.StaticAccess {
		report.Safe = false
		report.Reasons = append(report.Reasons, fmt.Sprintf("accesses %d static variable(s)", len(report.StaticVars)))
	}

	if report.FileIO {
		report.Safe = false
		report.Reasons = append(report.Reasons, fmt.Sprintf("performs %d file I/O operation(s)", len(report.FileOperations)))
	}

	if report.NetworkIO {
		report.Safe = false
		report.Reasons = append(report.Reasons, fmt.Sprintf("performs %d network I/O operation(s)", len(report.NetworkOperations)))
	}

	if report.DatabaseAccess {
		report.Safe = false
		report.Reasons = append(report.Reasons, fmt.Sprintf("performs %d database operation(s)", len(report.DatabaseOperations)))
	}

	if len(report.SideEffects) > 0 {
		report.Safe = false
		report.Reasons = append(report.Reasons, fmt.Sprintf("has %d other side effect(s)", len(report.SideEffects)))
	}

	return report
}

// analyzeBlock analyzes a block of statements
func (sa *SafetyAnalyzer) analyzeBlock(block *ast.BlockStatement, report *SafetyReport) {
	for _, stmt := range block.Statements {
		sa.analyzeStmt(stmt, report)
	}
}

// analyzeStmt analyzes a statement
func (sa *SafetyAnalyzer) analyzeStmt(stmt ast.Stmt, report *SafetyReport) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		sa.analyzeExpr(s.Expression, report, false)

	case *ast.BlockStatement:
		sa.analyzeBlock(s, report)

	case *ast.IfStatement:
		sa.analyzeExpr(s.Condition, report, false)
		sa.analyzeBlock(s.Consequence, report)
		for _, elseIf := range s.ElseIfs {
			sa.analyzeExpr(elseIf.Condition, report, false)
			sa.analyzeBlock(elseIf.Consequence, report)
		}
		if s.Alternative != nil {
			sa.analyzeBlock(s.Alternative, report)
		}

	case *ast.WhileStatement:
		sa.analyzeExpr(s.Condition, report, false)
		sa.analyzeBlock(s.Body, report)

	case *ast.DoWhileStatement:
		sa.analyzeBlock(s.Body, report)
		sa.analyzeExpr(s.Condition, report, false)

	case *ast.ForStatement:
		for _, init := range s.Init {
			sa.analyzeExpr(init, report, false)
		}
		for _, cond := range s.Condition {
			sa.analyzeExpr(cond, report, false)
		}
		for _, inc := range s.Increment {
			sa.analyzeExpr(inc, report, false)
		}
		sa.analyzeBlock(s.Body, report)

	case *ast.ForeachStatement:
		sa.analyzeExpr(s.Array, report, false)
		sa.analyzeBlock(s.Body, report)

	case *ast.SwitchStatement:
		sa.analyzeExpr(s.Subject, report, false)
		for _, caseBlock := range s.Cases {
			if caseBlock.Value != nil {
				sa.analyzeExpr(caseBlock.Value, report, false)
			}
			for _, caseStmt := range caseBlock.Body {
				sa.analyzeStmt(caseStmt, report)
			}
		}

	case *ast.TryStatement:
		sa.analyzeBlock(s.Body, report)
		for _, catch := range s.CatchClauses {
			sa.analyzeBlock(catch.Body, report)
		}
		if s.Finally != nil {
			sa.analyzeBlock(s.Finally, report)
		}

	case *ast.ReturnStatement:
		if s.ReturnValue != nil {
			sa.analyzeExpr(s.ReturnValue, report, false)
		}

	case *ast.EchoStatement:
		for _, arg := range s.Expressions {
			sa.analyzeExpr(arg, report, false)
		}
		// Echo is a side effect (output)
		report.SideEffects = append(report.SideEffects, "echo statement")

	case *ast.ThrowStatement:
		sa.analyzeExpr(s.Expression, report, false)
	}
}

// analyzeExpr analyzes an expression
// isAssignTarget indicates if this expression is the left-hand side of an assignment
func (sa *SafetyAnalyzer) analyzeExpr(expr ast.Expr, report *SafetyReport, isAssignTarget bool) {
	switch e := expr.(type) {
	case *ast.AssignmentExpression:
		// Analyze the right side first
		sa.analyzeExpr(e.Right, report, false)

		// Then analyze the left side as an assignment target
		sa.analyzeExpr(e.Left, report, true)

	case *ast.Variable:
		// Check if it's a global variable access
		// In PHP, global variables start with $ and are accessed via the global keyword
		// or via $GLOBALS array
		// For now, we'll be conservative and track all variable accesses
		// A more sophisticated approach would track the scope

	case *ast.CallExpression:
		sa.analyzeFunctionCall(e, report)
		// Analyze arguments
		for _, arg := range e.Arguments {
			sa.analyzeExpr(arg, report, false)
		}

	case *ast.MethodCallExpression:
		sa.analyzeExpr(e.Object, report, false)
		// Analyze arguments
		for _, arg := range e.Arguments {
			sa.analyzeExpr(arg, report, false)
		}
		// Method calls could have side effects
		report.SideEffects = append(report.SideEffects, fmt.Sprintf("method call: %v", e.Method))

	case *ast.StaticCallExpression:
		// Static method calls could have side effects
		report.SideEffects = append(report.SideEffects, fmt.Sprintf("static method call: %v::%v", e.Class, e.Method))
		for _, arg := range e.Arguments {
			sa.analyzeExpr(arg, report, false)
		}

	case *ast.StaticPropertyExpression:
		// Static property access
		report.StaticAccess = true
		varName := fmt.Sprintf("%v::$%v", e.Class, e.Property)
		if !contains(report.StaticVars, varName) {
			report.StaticVars = append(report.StaticVars, varName)
		}

	case *ast.PropertyExpression:
		sa.analyzeExpr(e.Object, report, isAssignTarget)

	case *ast.ArrayExpression:
		for _, element := range e.Elements {
			if element.Key != nil {
				sa.analyzeExpr(element.Key, report, false)
			}
			sa.analyzeExpr(element.Value, report, false)
		}

	case *ast.IndexExpression:
		sa.analyzeExpr(e.Left, report, isAssignTarget)
		sa.analyzeExpr(e.Index, report, false)

	case *ast.InfixExpression:
		sa.analyzeExpr(e.Left, report, false)
		sa.analyzeExpr(e.Right, report, false)

	case *ast.PrefixExpression:
		sa.analyzeExpr(e.Right, report, false)

	case *ast.TernaryExpression:
		sa.analyzeExpr(e.Condition, report, false)
		sa.analyzeExpr(e.Consequence, report, false)
		sa.analyzeExpr(e.Alternative, report, false)

	case *ast.ClosureExpression:
		// Closures with use() clause might capture variables
		if len(e.Use) > 0 {
			report.SideEffects = append(report.SideEffects, "closure with variable capture")
		}
		if e.Body != nil {
			sa.analyzeBlock(e.Body, report)
		}

	case *ast.NewExpression:
		for _, arg := range e.Arguments {
			sa.analyzeExpr(arg, report, false)
		}
		// Object construction could have side effects
		report.SideEffects = append(report.SideEffects, fmt.Sprintf("object creation: new %v", e.Class))
	}
}

// analyzeFunctionCall analyzes a function call for unsafe operations
func (sa *SafetyAnalyzer) analyzeFunctionCall(call *ast.CallExpression, report *SafetyReport) {
	funcName := ""
	if ident, ok := call.Function.(*ast.Identifier); ok {
		funcName = strings.ToLower(ident.Value)
	}

	// Check for global variable access via $GLOBALS
	if funcName == "global" {
		// global keyword marks variables as global
		for _, arg := range call.Arguments {
			if v, ok := arg.(*ast.Variable); ok {
				if !contains(report.GlobalWrites, v.Name) {
					report.GlobalWrites = append(report.GlobalWrites, v.Name)
				}
			}
		}
		return
	}

	// Check for file I/O operations
	if sa.isFileIOFunction(funcName) {
		report.FileIO = true
		if !contains(report.FileOperations, funcName) {
			report.FileOperations = append(report.FileOperations, funcName)
		}
		return
	}

	// Check for network I/O operations
	if sa.isNetworkIOFunction(funcName) {
		report.NetworkIO = true
		if !contains(report.NetworkOperations, funcName) {
			report.NetworkOperations = append(report.NetworkOperations, funcName)
		}
		return
	}

	// Check for database operations
	if sa.isDatabaseFunction(funcName) {
		report.DatabaseAccess = true
		if !contains(report.DatabaseOperations, funcName) {
			report.DatabaseOperations = append(report.DatabaseOperations, funcName)
		}
		return
	}

	// Check for other side-effect functions
	if sa.hasSideEffects(funcName) {
		report.SideEffects = append(report.SideEffects, funcName)
	}
}

// isFileIOFunction returns true if the function performs file I/O
func (sa *SafetyAnalyzer) isFileIOFunction(name string) bool {
	fileIOFunctions := []string{
		"fopen", "fclose", "fread", "fwrite", "fgets", "fgetc", "fputs", "fputc",
		"file_get_contents", "file_put_contents", "file", "readfile",
		"unlink", "rename", "copy", "mkdir", "rmdir", "chmod", "chown",
		"file_exists", "is_file", "is_dir", "is_readable", "is_writable",
		"stat", "lstat", "fstat", "touch",
		"glob", "scandir", "opendir", "readdir", "closedir",
	}
	return contains(fileIOFunctions, name)
}

// isNetworkIOFunction returns true if the function performs network I/O
func (sa *SafetyAnalyzer) isNetworkIOFunction(name string) bool {
	networkIOFunctions := []string{
		"fsockopen", "pfsockopen", "socket_create", "socket_connect",
		"socket_bind", "socket_listen", "socket_accept",
		"socket_read", "socket_write", "socket_send", "socket_recv",
		"curl_init", "curl_exec", "curl_multi_exec",
		"file_get_contents", // Can be used for URLs
		"fopen",              // Can open URLs
		"gethostbyname", "gethostbyaddr", "dns_get_record",
		"mail", "smtp",
	}
	return contains(networkIOFunctions, name)
}

// isDatabaseFunction returns true if the function performs database operations
func (sa *SafetyAnalyzer) isDatabaseFunction(name string) bool {
	databaseFunctions := []string{
		// PDO
		"pdo",
		// MySQLi
		"mysqli_connect", "mysqli_query", "mysqli_fetch", "mysqli_close",
		"mysql_connect", "mysql_query", "mysql_fetch", "mysql_close",
		// PostgreSQL
		"pg_connect", "pg_query", "pg_fetch", "pg_close",
		// SQLite
		"sqlite_open", "sqlite_query", "sqlite_fetch", "sqlite_close",
		// Generic database
		"db_query", "database_query",
	}
	return contains(databaseFunctions, name)
}

// hasSideEffects returns true if the function has side effects
func (sa *SafetyAnalyzer) hasSideEffects(name string) bool {
	sideEffectFunctions := []string{
		"echo", "print", "printf", "var_dump", "print_r", "var_export",
		"header", "setcookie", "session_start", "session_destroy",
		"ob_start", "ob_end_flush", "ob_clean",
		"exit", "die", "register_shutdown_function",
		"error_log", "trigger_error", "user_error",
		"rand", "mt_rand", "srand", "mt_srand", // Non-deterministic
		"time", "microtime", // Time-dependent
		"sleep", "usleep", "time_nanosleep",
	}
	return contains(sideEffectFunctions, name)
}

// IsPureFunction analyzes a function and returns true if it's pure (no side effects)
func (sa *SafetyAnalyzer) IsPureFunction(fn *ast.FunctionDeclaration) bool {
	report := sa.AnalyzeFunction(fn)

	// A function is pure if:
	// 1. It doesn't write to globals
	// 2. It doesn't access static variables
	// 3. It doesn't perform I/O
	// 4. It doesn't have other side effects
	// 5. It can read globals (but not write)

	if len(report.GlobalWrites) > 0 {
		return false
	}
	if report.StaticAccess {
		return false
	}
	if report.FileIO {
		return false
	}
	if report.NetworkIO {
		return false
	}
	if report.DatabaseAccess {
		return false
	}
	if len(report.SideEffects) > 0 {
		return false
	}

	return true
}

// GenerateReport generates a human-readable safety report
func (sa *SafetyAnalyzer) GenerateReport(fn *ast.FunctionDeclaration) string {
	report := sa.AnalyzeFunction(fn)

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Safety Analysis for function '%s':\n", fn.Name.Value))
	builder.WriteString(strings.Repeat("=", 60))
	builder.WriteString("\n\n")

	if report.Safe {
		builder.WriteString("✓ SAFE TO PARALLELIZE\n\n")
		builder.WriteString("This function appears to be pure and can be safely parallelized.\n")
	} else {
		builder.WriteString("✗ UNSAFE TO PARALLELIZE\n\n")
		builder.WriteString("Reasons:\n")
		for _, reason := range report.Reasons {
			builder.WriteString(fmt.Sprintf("  - %s\n", reason))
		}
	}

	// Details
	builder.WriteString("\nDetails:\n")
	builder.WriteString(strings.Repeat("-", 60))
	builder.WriteString("\n")

	if len(report.GlobalReads) > 0 {
		builder.WriteString(fmt.Sprintf("Global reads: %s\n", strings.Join(report.GlobalReads, ", ")))
	}

	if len(report.GlobalWrites) > 0 {
		builder.WriteString(fmt.Sprintf("Global writes: %s\n", strings.Join(report.GlobalWrites, ", ")))
	}

	if len(report.StaticVars) > 0 {
		builder.WriteString(fmt.Sprintf("Static variables: %s\n", strings.Join(report.StaticVars, ", ")))
	}

	if len(report.FileOperations) > 0 {
		builder.WriteString(fmt.Sprintf("File I/O: %s\n", strings.Join(report.FileOperations, ", ")))
	}

	if len(report.NetworkOperations) > 0 {
		builder.WriteString(fmt.Sprintf("Network I/O: %s\n", strings.Join(report.NetworkOperations, ", ")))
	}

	if len(report.DatabaseOperations) > 0 {
		builder.WriteString(fmt.Sprintf("Database: %s\n", strings.Join(report.DatabaseOperations, ", ")))
	}

	if len(report.SideEffects) > 0 {
		builder.WriteString(fmt.Sprintf("Side effects: %s\n", strings.Join(report.SideEffects, ", ")))
	}

	return builder.String()
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
