package vm

// Opcode represents a single VM instruction type
type Opcode uint8

// String returns the name of the opcode for debugging
func (op Opcode) String() string {
	if int(op) < len(opcodeNames) {
		return opcodeNames[op]
	}
	return "UNKNOWN"
}

// VM Opcodes - Directly from PHP's Zend Engine (php-src/Zend/zend_vm_opcodes.h)
// These opcodes define all operations the VM can perform
const (
	// ========================================
	// Basic Operations
	// ========================================

	// OpNop - No operation, used as placeholder
	OpNop Opcode = 0

	// ========================================
	// Arithmetic Operations (1-12)
	// ========================================

	// OpAdd - Add two values: result = op1 + op2
	OpAdd Opcode = 1

	// OpSub - Subtract two values: result = op1 - op2
	OpSub Opcode = 2

	// OpMul - Multiply two values: result = op1 * op2
	OpMul Opcode = 3

	// OpDiv - Divide two values: result = op1 / op2
	OpDiv Opcode = 4

	// OpMod - Modulo operation: result = op1 % op2
	OpMod Opcode = 5

	// OpSL - Bitwise shift left: result = op1 << op2
	OpSL Opcode = 6

	// OpSR - Bitwise shift right: result = op1 >> op2
	OpSR Opcode = 7

	// OpConcat - String concatenation: result = op1 . op2
	OpConcat Opcode = 8

	// OpBWOr - Bitwise OR: result = op1 | op2
	OpBWOr Opcode = 9

	// OpBWAnd - Bitwise AND: result = op1 & op2
	OpBWAnd Opcode = 10

	// OpBWXor - Bitwise XOR: result = op1 ^ op2
	OpBWXor Opcode = 11

	// OpPow - Power/Exponentiation: result = op1 ** op2
	OpPow Opcode = 12

	// ========================================
	// Unary Operations (13-15)
	// ========================================

	// OpBWNot - Bitwise NOT: result = ~op1
	OpBWNot Opcode = 13

	// OpBoolNot - Boolean NOT: result = !op1
	OpBoolNot Opcode = 14

	// OpBoolXor - Boolean XOR: result = op1 xor op2
	OpBoolXor Opcode = 15

	// ========================================
	// Comparison Operations (16-21)
	// ========================================

	// OpIsIdentical - Identical comparison: result = op1 === op2
	OpIsIdentical Opcode = 16

	// OpIsNotIdentical - Not identical comparison: result = op1 !== op2
	OpIsNotIdentical Opcode = 17

	// OpIsEqual - Equal comparison: result = op1 == op2
	OpIsEqual Opcode = 18

	// OpIsNotEqual - Not equal comparison: result = op1 != op2
	OpIsNotEqual Opcode = 19

	// OpIsSmaller - Less than comparison: result = op1 < op2
	OpIsSmaller Opcode = 20

	// OpIsSmallerOrEqual - Less than or equal comparison: result = op1 <= op2
	OpIsSmallerOrEqual Opcode = 21

	// ========================================
	// Assignment Operations (22-33)
	// ========================================

	// OpAssign - Variable assignment: op1 = op2
	OpAssign Opcode = 22

	// OpAssignDim - Array element assignment: op1[op2] = result
	OpAssignDim Opcode = 23

	// OpAssignObj - Object property assignment: op1->op2 = result
	OpAssignObj Opcode = 24

	// OpAssignStaticProp - Static property assignment: Class::$prop = result
	OpAssignStaticProp Opcode = 25

	// OpAssignOp - Compound assignment: op1 += op2, op1 -= op2, etc.
	OpAssignOp Opcode = 26

	// OpAssignDimOp - Compound array assignment: op1[op2] += result
	OpAssignDimOp Opcode = 27

	// OpAssignObjOp - Compound object assignment: op1->op2 += result
	OpAssignObjOp Opcode = 28

	// OpAssignStaticPropOp - Compound static property assignment: Class::$prop += result
	OpAssignStaticPropOp Opcode = 29

	// OpAssignRef - Assignment by reference: op1 =& op2
	OpAssignRef Opcode = 30

	// OpQMAssign - Quick assign (no side effects): result = op1
	OpQMAssign Opcode = 31

	// OpAssignObjRef - Object property assignment by reference: op1->op2 =& result
	OpAssignObjRef Opcode = 32

	// OpAssignStaticPropRef - Static property assignment by reference: Class::$prop =& result
	OpAssignStaticPropRef Opcode = 33

	// ========================================
	// Increment/Decrement Operations (34-41)
	// ========================================

	// OpPreInc - Pre-increment: result = ++op1
	OpPreInc Opcode = 34

	// OpPreDec - Pre-decrement: result = --op1
	OpPreDec Opcode = 35

	// OpPostInc - Post-increment: result = op1++
	OpPostInc Opcode = 36

	// OpPostDec - Post-decrement: result = op1--
	OpPostDec Opcode = 37

	// OpPreIncStaticProp - Pre-increment static property: result = ++Class::$prop
	OpPreIncStaticProp Opcode = 38

	// OpPreDecStaticProp - Pre-decrement static property: result = --Class::$prop
	OpPreDecStaticProp Opcode = 39

	// OpPostIncStaticProp - Post-increment static property: result = Class::$prop++
	OpPostIncStaticProp Opcode = 40

	// OpPostDecStaticProp - Post-decrement static property: result = Class::$prop--
	OpPostDecStaticProp Opcode = 41

	// ========================================
	// Control Flow - Jumps (42-48)
	// ========================================

	// OpJmp - Unconditional jump to address
	OpJmp Opcode = 42

	// OpJmpZ - Jump if zero (false): if (!op1) goto address
	OpJmpZ Opcode = 43

	// OpJmpNZ - Jump if not zero (true): if (op1) goto address
	OpJmpNZ Opcode = 44

	// 45 is missing in PHP opcodes

	// OpJmpZEx - Jump if zero with assignment: result = op1; if (!result) goto address
	OpJmpZEx Opcode = 46

	// OpJmpNZEx - Jump if not zero with assignment: result = op1; if (result) goto address
	OpJmpNZEx Opcode = 47

	// OpCase - Switch case comparison: if (op1 == op2) goto address
	OpCase Opcode = 48

	// ========================================
	// Variable/Type Operations (49-52)
	// ========================================

	// OpCheckVar - Check if variable is defined
	OpCheckVar Opcode = 49

	// OpSendVarNoRefEx - Send variable by value (no reference allowed)
	OpSendVarNoRefEx Opcode = 50

	// OpCast - Type cast: result = (type)op1
	OpCast Opcode = 51

	// OpBool - Convert to boolean: result = (bool)op1
	OpBool Opcode = 52

	// ========================================
	// String Operations (53-56)
	// ========================================

	// OpFastConcat - Optimized string concatenation
	OpFastConcat Opcode = 53

	// OpRopeInit - Initialize rope for multiple concatenations
	OpRopeInit Opcode = 54

	// OpRopeAdd - Add to rope
	OpRopeAdd Opcode = 55

	// OpRopeEnd - Finalize rope into string
	OpRopeEnd Opcode = 56

	// ========================================
	// Error Suppression (57-58)
	// ========================================

	// OpBeginSilence - Begin error suppression (@)
	OpBeginSilence Opcode = 57

	// OpEndSilence - End error suppression
	OpEndSilence Opcode = 58

	// ========================================
	// Function Call Operations (59-67)
	// ========================================

	// OpInitFcallByName - Initialize function call by name
	OpInitFcallByName Opcode = 59

	// OpDoFcall - Execute function call
	OpDoFcall Opcode = 60

	// OpInitFcall - Initialize optimized function call
	OpInitFcall Opcode = 61

	// OpReturn - Return from function
	OpReturn Opcode = 62

	// OpRecv - Receive function parameter
	OpRecv Opcode = 63

	// OpRecvInit - Receive parameter with default value
	OpRecvInit Opcode = 64

	// OpSendVal - Send value argument to function
	OpSendVal Opcode = 65

	// OpSendVarEx - Send variable argument to function
	OpSendVarEx Opcode = 66

	// OpSendRef - Send argument by reference
	OpSendRef Opcode = 67

	// ========================================
	// Object Operations (68-70)
	// ========================================

	// OpNew - Create new object instance: result = new Class()
	OpNew Opcode = 68

	// OpInitNsFcallByName - Initialize namespaced function call
	OpInitNsFcallByName Opcode = 69

	// OpFree - Free temporary variable
	OpFree Opcode = 70

	// ========================================
	// Array Operations (71-72)
	// ========================================

	// OpInitArray - Initialize array: result = []
	OpInitArray Opcode = 71

	// OpAddArrayElement - Add element to array: result[] = op1 or result[op2] = op1
	OpAddArrayElement Opcode = 72

	// ========================================
	// Include/Eval Operations (73)
	// ========================================

	// OpIncludeOrEval - Execute include, require, eval
	OpIncludeOrEval Opcode = 73

	// ========================================
	// Unset Operations (74-76)
	// ========================================

	// OpUnsetVar - Unset variable: unset($var)
	OpUnsetVar Opcode = 74

	// OpUnsetDim - Unset array element: unset($arr[$key])
	OpUnsetDim Opcode = 75

	// OpUnsetObj - Unset object property: unset($obj->prop)
	OpUnsetObj Opcode = 76

	// ========================================
	// Foreach Operations (77-78)
	// ========================================

	// OpFeResetR - Reset foreach for read: foreach ($arr as $val)
	OpFeResetR Opcode = 77

	// OpFeFetchR - Fetch next foreach element for read
	OpFeFetchR Opcode = 78

	// 79 is missing in PHP opcodes

	// ========================================
	// Fetch Operations - Read (80-82)
	// ========================================

	// OpFetchR - Fetch variable for read: result = $var
	OpFetchR Opcode = 80

	// OpFetchDimR - Fetch array element for read: result = $arr[$key]
	OpFetchDimR Opcode = 81

	// OpFetchObjR - Fetch object property for read: result = $obj->prop
	OpFetchObjR Opcode = 82

	// ========================================
	// Fetch Operations - Write (83-85)
	// ========================================

	// OpFetchW - Fetch variable for write
	OpFetchW Opcode = 83

	// OpFetchDimW - Fetch array element for write: $arr[$key] = ...
	OpFetchDimW Opcode = 84

	// OpFetchObjW - Fetch object property for write: $obj->prop = ...
	OpFetchObjW Opcode = 85

	// ========================================
	// Fetch Operations - Read-Write (86-88)
	// ========================================

	// OpFetchRW - Fetch variable for read-write: $var += 1
	OpFetchRW Opcode = 86

	// OpFetchDimRW - Fetch array element for read-write: $arr[$key] += 1
	OpFetchDimRW Opcode = 87

	// OpFetchObjRW - Fetch object property for read-write: $obj->prop += 1
	OpFetchObjRW Opcode = 88

	// ========================================
	// Fetch Operations - IsSet (89-91)
	// ========================================

	// OpFetchIs - Fetch variable for isset/empty check
	OpFetchIs Opcode = 89

	// OpFetchDimIs - Fetch array element for isset/empty check
	OpFetchDimIs Opcode = 90

	// OpFetchObjIs - Fetch object property for isset/empty check
	OpFetchObjIs Opcode = 91

	// ========================================
	// Fetch Operations - Function Argument (92-94)
	// ========================================

	// OpFetchFuncArg - Fetch variable as function argument
	OpFetchFuncArg Opcode = 92

	// OpFetchDimFuncArg - Fetch array element as function argument
	OpFetchDimFuncArg Opcode = 93

	// OpFetchObjFuncArg - Fetch object property as function argument
	OpFetchObjFuncArg Opcode = 94

	// ========================================
	// Fetch Operations - Unset (95-97)
	// ========================================

	// OpFetchUnset - Fetch variable for unset
	OpFetchUnset Opcode = 95

	// OpFetchDimUnset - Fetch array element for unset
	OpFetchDimUnset Opcode = 96

	// OpFetchObjUnset - Fetch object property for unset
	OpFetchObjUnset Opcode = 97

	// ========================================
	// Fetch Operations - List (98-99)
	// ========================================

	// OpFetchListR - Fetch for list() assignment
	OpFetchListR Opcode = 98

	// OpFetchConstant - Fetch constant value
	OpFetchConstant Opcode = 99

	// ========================================
	// Function/Extension Operations (100-106)
	// ========================================

	// OpCheckFuncArg - Check function argument type
	OpCheckFuncArg Opcode = 100

	// OpExtStmt - Extension statement hook
	OpExtStmt Opcode = 101

	// OpExtFcallBegin - Extension function call begin hook
	OpExtFcallBegin Opcode = 102

	// OpExtFcallEnd - Extension function call end hook
	OpExtFcallEnd Opcode = 103

	// OpExtNop - Extension no-operation hook
	OpExtNop Opcode = 104

	// OpTicks - Execute tick function
	OpTicks Opcode = 105

	// OpSendVarNoRef - Send variable without reference
	OpSendVarNoRef Opcode = 106

	// ========================================
	// Exception Operations (107-108)
	// ========================================

	// OpCatch - Catch exception
	OpCatch Opcode = 107

	// OpThrow - Throw exception
	OpThrow Opcode = 108

	// ========================================
	// Class Operations (109-113)
	// ========================================

	// OpFetchClass - Fetch class entry by name
	OpFetchClass Opcode = 109

	// OpClone - Clone object: result = clone $obj
	OpClone Opcode = 110

	// OpReturnByRef - Return by reference
	OpReturnByRef Opcode = 111

	// OpInitMethodCall - Initialize method call: $obj->method()
	OpInitMethodCall Opcode = 112

	// OpInitStaticMethodCall - Initialize static method call: Class::method()
	OpInitStaticMethodCall Opcode = 113

	// ========================================
	// IsSet/Empty Operations (114-115)
	// ========================================

	// OpIssetIsemptyVar - Check isset/empty on variable
	OpIssetIsemptyVar Opcode = 114

	// OpIssetIsemptyDimObj - Check isset/empty on array element or object property
	OpIssetIsemptyDimObj Opcode = 115

	// ========================================
	// More Function Call Operations (116-120)
	// ========================================

	// OpSendValEx - Send value argument (extended)
	OpSendValEx Opcode = 116

	// OpSendVar - Send variable argument
	OpSendVar Opcode = 117

	// OpInitUserCall - Initialize user function call
	OpInitUserCall Opcode = 118

	// OpSendArray - Send array as argument list
	OpSendArray Opcode = 119

	// OpSendUser - Send user argument
	OpSendUser Opcode = 120

	// ========================================
	// Built-in Functions (121-123)
	// ========================================

	// OpStrlen - Get string length: result = strlen(op1)
	OpStrlen Opcode = 121

	// OpDefined - Check if constant is defined
	OpDefined Opcode = 122

	// OpTypeCheck - Check variable type
	OpTypeCheck Opcode = 123

	// OpVerifyReturnType - Verify function return type matches declaration
	OpVerifyReturnType Opcode = 124

	// ========================================
	// Foreach Operations - Write (125-127)
	// ========================================

	// OpFeResetRW - Reset foreach for read-write: foreach ($arr as &$val)
	OpFeResetRW Opcode = 125

	// OpFeFetchRW - Fetch next foreach element for read-write
	OpFeFetchRW Opcode = 126

	// OpFeFree - Free foreach iterator
	OpFeFree Opcode = 127

	// ========================================
	// Dynamic Call Operations (128-131)
	// ========================================

	// OpInitDynamicCall - Initialize dynamic function call: $func()
	OpInitDynamicCall Opcode = 128

	// OpDoIcall - Execute internal function call
	OpDoIcall Opcode = 129

	// OpDoUcall - Execute user function call
	OpDoUcall Opcode = 130

	// OpDoFcallByName - Execute function call by name
	OpDoFcallByName Opcode = 131

	// ========================================
	// Object Inc/Dec Operations (132-135)
	// ========================================

	// OpPreIncObj - Pre-increment object property: result = ++$obj->prop
	OpPreIncObj Opcode = 132

	// OpPreDecObj - Pre-decrement object property: result = --$obj->prop
	OpPreDecObj Opcode = 133

	// OpPostIncObj - Post-increment object property: result = $obj->prop++
	OpPostIncObj Opcode = 134

	// OpPostDecObj - Post-decrement object property: result = $obj->prop--
	OpPostDecObj Opcode = 135

	// ========================================
	// Echo Operation (136-137)
	// ========================================

	// OpEcho - Output string: echo $str
	OpEcho Opcode = 136

	// OpOpData - Additional operand data for previous instruction
	OpOpData Opcode = 137

	// ========================================
	// Type Operations (138)
	// ========================================

	// OpInstanceof - Check if object is instance of class: result = $obj instanceof Class
	OpInstanceof Opcode = 138

	// ========================================
	// Generator Operations (139-142)
	// ========================================

	// OpGeneratorCreate - Create generator object
	OpGeneratorCreate Opcode = 139

	// OpMakeRef - Make variable a reference
	OpMakeRef Opcode = 140

	// OpDeclareFunction - Declare function
	OpDeclareFunction Opcode = 141

	// OpDeclareLambdaFunction - Declare anonymous function/closure
	OpDeclareLambdaFunction Opcode = 142

	// ========================================
	// Declaration Operations (143-146)
	// ========================================

	// OpDeclareConst - Declare constant
	OpDeclareConst Opcode = 143

	// OpDeclareClass - Declare class
	OpDeclareClass Opcode = 144

	// OpDeclareClassDelayed - Declare class (delayed binding)
	OpDeclareClassDelayed Opcode = 145

	// OpDeclareAnonClass - Declare anonymous class
	OpDeclareAnonClass Opcode = 146

	// ========================================
	// Array Operations (147)
	// ========================================

	// OpAddArrayUnpack - Add unpacked array elements: [...$arr]
	OpAddArrayUnpack Opcode = 147

	// ========================================
	// Property Operations (148)
	// ========================================

	// OpIssetIsemptyPropObj - Check isset/empty on object property
	OpIssetIsemptyPropObj Opcode = 148

	// ========================================
	// Exception Handling (149-150)
	// ========================================

	// OpHandleException - Handle exception in catch block
	OpHandleException Opcode = 149

	// OpUserOpcode - User-defined opcode handler
	OpUserOpcode Opcode = 150

	// ========================================
	// Assertion Operations (151)
	// ========================================

	// OpAssertCheck - Check assertion
	OpAssertCheck Opcode = 151

	// ========================================
	// Control Flow - Special Jumps (152)
	// ========================================

	// OpJmpSet - Jump and set: if (op1) { result = op1; goto address; }
	OpJmpSet Opcode = 152

	// ========================================
	// CV Operations (153-154)
	// ========================================

	// OpUnsetCV - Unset compiled variable
	OpUnsetCV Opcode = 153

	// OpIssetIsemptyCV - Check isset/empty on compiled variable
	OpIssetIsemptyCV Opcode = 154

	// ========================================
	// List Operations (155)
	// ========================================

	// OpFetchListW - Fetch for list() assignment (write)
	OpFetchListW Opcode = 155

	// ========================================
	// Variable Operations (156-157)
	// ========================================

	// OpSeparate - Separate variable (copy-on-write)
	OpSeparate Opcode = 156

	// OpFetchClassName - Fetch class name from object
	OpFetchClassName Opcode = 157

	// ========================================
	// Trampoline Operations (158-159)
	// ========================================

	// OpCallTrampoline - Call trampoline function
	OpCallTrampoline Opcode = 158

	// OpDiscardException - Discard exception
	OpDiscardException Opcode = 159

	// ========================================
	// Generator/Yield Operations (160-163)
	// ========================================

	// OpYield - Yield value from generator
	OpYield Opcode = 160

	// OpGeneratorReturn - Return from generator
	OpGeneratorReturn Opcode = 161

	// OpFastCall - Fast call for internal use
	OpFastCall Opcode = 162

	// OpFastRet - Fast return for internal use
	OpFastRet Opcode = 163

	// ========================================
	// Variadic Operations (164-165)
	// ========================================

	// OpRecvVariadic - Receive variadic parameters: ...$args
	OpRecvVariadic Opcode = 164

	// OpSendUnpack - Send unpacked arguments: ...$ args
	OpSendUnpack Opcode = 165

	// ========================================
	// Generator Operations (166-167)
	// ========================================

	// OpYieldFrom - Yield from another generator: yield from $generator
	OpYieldFrom Opcode = 166

	// OpCopyTmp - Copy temporary variable
	OpCopyTmp Opcode = 167

	// ========================================
	// Global/Static Operations (168-169)
	// ========================================

	// OpBindGlobal - Bind global variable
	OpBindGlobal Opcode = 168

	// OpCoalesce - Null coalescing: result = op1 ?? op2
	OpCoalesce Opcode = 169

	// ========================================
	// Special Operations (170-172)
	// ========================================

	// OpSpaceship - Spaceship operator: result = op1 <=> op2
	OpSpaceship Opcode = 170

	// OpFuncNumArgs - Get number of function arguments: func_num_args()
	OpFuncNumArgs Opcode = 171

	// OpFuncGetArgs - Get function arguments: func_get_args()
	OpFuncGetArgs Opcode = 172

	// ========================================
	// Static Property Operations (173-180)
	// ========================================

	// OpFetchStaticPropR - Fetch static property for read
	OpFetchStaticPropR Opcode = 173

	// OpFetchStaticPropW - Fetch static property for write
	OpFetchStaticPropW Opcode = 174

	// OpFetchStaticPropRW - Fetch static property for read-write
	OpFetchStaticPropRW Opcode = 175

	// OpFetchStaticPropIs - Fetch static property for isset/empty check
	OpFetchStaticPropIs Opcode = 176

	// OpFetchStaticPropFuncArg - Fetch static property as function argument
	OpFetchStaticPropFuncArg Opcode = 177

	// OpFetchStaticPropUnset - Fetch static property for unset
	OpFetchStaticPropUnset Opcode = 178

	// OpUnsetStaticProp - Unset static property
	OpUnsetStaticProp Opcode = 179

	// OpIssetIsemptyStaticProp - Check isset/empty on static property
	OpIssetIsemptyStaticProp Opcode = 180

	// ========================================
	// Class Constant Operations (181)
	// ========================================

	// OpFetchClassConstant - Fetch class constant: Class::CONST
	OpFetchClassConstant Opcode = 181

	// ========================================
	// Closure Operations (182-184)
	// ========================================

	// OpBindLexical - Bind lexical variables for closure: use ($var)
	OpBindLexical Opcode = 182

	// OpBindStatic - Bind static variables
	OpBindStatic Opcode = 183

	// OpFetchThis - Fetch $this variable
	OpFetchThis Opcode = 184

	// ========================================
	// Function Argument Operations (185-186)
	// ========================================

	// OpSendFuncArg - Send function argument
	OpSendFuncArg Opcode = 185

	// OpIssetIsemptyThis - Check isset/empty on $this
	OpIssetIsemptyThis Opcode = 186

	// ========================================
	// Switch Operations (187-188)
	// ========================================

	// OpSwitchLong - Optimized switch for integer values
	OpSwitchLong Opcode = 187

	// OpSwitchString - Optimized switch for string values
	OpSwitchString Opcode = 188

	// ========================================
	// Array Operations (189-190)
	// ========================================

	// OpInArray - Check if value is in array: in_array($needle, $haystack)
	OpInArray Opcode = 189

	// OpCount - Count array elements: count($arr)
	OpCount Opcode = 190

	// ========================================
	// Class Information Operations (191-193)
	// ========================================

	// OpGetClass - Get class name: get_class($obj)
	OpGetClass Opcode = 191

	// OpGetCalledClass - Get called class: get_called_class()
	OpGetCalledClass Opcode = 192

	// OpGetType - Get variable type: get_type($var)
	OpGetType Opcode = 193

	// ========================================
	// Array Key Operations (194)
	// ========================================

	// OpArrayKeyExists - Check if array key exists: array_key_exists($key, $arr)
	OpArrayKeyExists Opcode = 194

	// ========================================
	// Match Operations (195-197)
	// ========================================

	// OpMatch - Match expression (PHP 8.0+)
	OpMatch Opcode = 195

	// OpCaseStrict - Strict case comparison for match
	OpCaseStrict Opcode = 196

	// OpMatchError - Throw UnhandledMatchError
	OpMatchError Opcode = 197

	// ========================================
	// Null Operations (198)
	// ========================================

	// OpJmpNull - Jump if null: if (op1 === null) goto address
	OpJmpNull Opcode = 198

	// ========================================
	// Argument Operations (199)
	// ========================================

	// OpCheckUndefArgs - Check for undefined arguments
	OpCheckUndefArgs Opcode = 199

	// ========================================
	// Global Operations (200)
	// ========================================

	// OpFetchGlobals - Fetch global variables
	OpFetchGlobals Opcode = 200

	// ========================================
	// Type Verification (201)
	// ========================================

	// OpVerifyNeverType - Verify never type (PHP 8.1+)
	OpVerifyNeverType Opcode = 201

	// ========================================
	// Callable Operations (202)
	// ========================================

	// OpCallableConvert - Convert to first-class callable (PHP 8.1+)
	OpCallableConvert Opcode = 202

	// ========================================
	// Static Initialization (203)
	// ========================================

	// OpBindInitStaticOrJmp - Bind and initialize static variable or jump
	OpBindInitStaticOrJmp Opcode = 203

	// ========================================
	// Frameless Internal Calls (204-207)
	// ========================================

	// OpFramelessIcall0 - Frameless internal call with 0 arguments
	OpFramelessIcall0 Opcode = 204

	// OpFramelessIcall1 - Frameless internal call with 1 argument
	OpFramelessIcall1 Opcode = 205

	// OpFramelessIcall2 - Frameless internal call with 2 arguments
	OpFramelessIcall2 Opcode = 206

	// OpFramelessIcall3 - Frameless internal call with 3 arguments
	OpFramelessIcall3 Opcode = 207

	// ========================================
	// Jump Operations (208)
	// ========================================

	// OpJmpFrameless - Frameless jump
	OpJmpFrameless Opcode = 208

	// ========================================
	// Property Hook Operations (209)
	// ========================================

	// OpInitParentPropertyHookCall - Initialize parent property hook call (PHP 8.4+)
	OpInitParentPropertyHookCall Opcode = 209

	// ========================================
	// Declaration Operations (210)
	// ========================================

	// OpDeclareAttributedConst - Declare constant with attributes (PHP 8.4+)
	OpDeclareAttributedConst Opcode = 210
)

// Total number of opcodes
const OpcodeLast = 210

// opcodeNames maps opcodes to their string names for debugging
var opcodeNames = [OpcodeLast + 1]string{
	OpNop:                            "NOP",
	OpAdd:                            "ADD",
	OpSub:                            "SUB",
	OpMul:                            "MUL",
	OpDiv:                            "DIV",
	OpMod:                            "MOD",
	OpSL:                             "SL",
	OpSR:                             "SR",
	OpConcat:                         "CONCAT",
	OpBWOr:                           "BW_OR",
	OpBWAnd:                          "BW_AND",
	OpBWXor:                          "BW_XOR",
	OpPow:                            "POW",
	OpBWNot:                          "BW_NOT",
	OpBoolNot:                        "BOOL_NOT",
	OpBoolXor:                        "BOOL_XOR",
	OpIsIdentical:                    "IS_IDENTICAL",
	OpIsNotIdentical:                 "IS_NOT_IDENTICAL",
	OpIsEqual:                        "IS_EQUAL",
	OpIsNotEqual:                     "IS_NOT_EQUAL",
	OpIsSmaller:                      "IS_SMALLER",
	OpIsSmallerOrEqual:               "IS_SMALLER_OR_EQUAL",
	OpAssign:                         "ASSIGN",
	OpAssignDim:                      "ASSIGN_DIM",
	OpAssignObj:                      "ASSIGN_OBJ",
	OpAssignStaticProp:               "ASSIGN_STATIC_PROP",
	OpAssignOp:                       "ASSIGN_OP",
	OpAssignDimOp:                    "ASSIGN_DIM_OP",
	OpAssignObjOp:                    "ASSIGN_OBJ_OP",
	OpAssignStaticPropOp:             "ASSIGN_STATIC_PROP_OP",
	OpAssignRef:                      "ASSIGN_REF",
	OpQMAssign:                       "QM_ASSIGN",
	OpAssignObjRef:                   "ASSIGN_OBJ_REF",
	OpAssignStaticPropRef:            "ASSIGN_STATIC_PROP_REF",
	OpPreInc:                         "PRE_INC",
	OpPreDec:                         "PRE_DEC",
	OpPostInc:                        "POST_INC",
	OpPostDec:                        "POST_DEC",
	OpPreIncStaticProp:               "PRE_INC_STATIC_PROP",
	OpPreDecStaticProp:               "PRE_DEC_STATIC_PROP",
	OpPostIncStaticProp:              "POST_INC_STATIC_PROP",
	OpPostDecStaticProp:              "POST_DEC_STATIC_PROP",
	OpJmp:                            "JMP",
	OpJmpZ:                           "JMPZ",
	OpJmpNZ:                          "JMPNZ",
	OpJmpZEx:                         "JMPZ_EX",
	OpJmpNZEx:                        "JMPNZ_EX",
	OpCase:                           "CASE",
	OpCheckVar:                       "CHECK_VAR",
	OpSendVarNoRefEx:                 "SEND_VAR_NO_REF_EX",
	OpCast:                           "CAST",
	OpBool:                           "BOOL",
	OpFastConcat:                     "FAST_CONCAT",
	OpRopeInit:                       "ROPE_INIT",
	OpRopeAdd:                        "ROPE_ADD",
	OpRopeEnd:                        "ROPE_END",
	OpBeginSilence:                   "BEGIN_SILENCE",
	OpEndSilence:                     "END_SILENCE",
	OpInitFcallByName:                "INIT_FCALL_BY_NAME",
	OpDoFcall:                        "DO_FCALL",
	OpInitFcall:                      "INIT_FCALL",
	OpReturn:                         "RETURN",
	OpRecv:                           "RECV",
	OpRecvInit:                       "RECV_INIT",
	OpSendVal:                        "SEND_VAL",
	OpSendVarEx:                      "SEND_VAR_EX",
	OpSendRef:                        "SEND_REF",
	OpNew:                            "NEW",
	OpInitNsFcallByName:              "INIT_NS_FCALL_BY_NAME",
	OpFree:                           "FREE",
	OpInitArray:                      "INIT_ARRAY",
	OpAddArrayElement:                "ADD_ARRAY_ELEMENT",
	OpIncludeOrEval:                  "INCLUDE_OR_EVAL",
	OpUnsetVar:                       "UNSET_VAR",
	OpUnsetDim:                       "UNSET_DIM",
	OpUnsetObj:                       "UNSET_OBJ",
	OpFeResetR:                       "FE_RESET_R",
	OpFeFetchR:                       "FE_FETCH_R",
	OpFetchR:                         "FETCH_R",
	OpFetchDimR:                      "FETCH_DIM_R",
	OpFetchObjR:                      "FETCH_OBJ_R",
	OpFetchW:                         "FETCH_W",
	OpFetchDimW:                      "FETCH_DIM_W",
	OpFetchObjW:                      "FETCH_OBJ_W",
	OpFetchRW:                        "FETCH_RW",
	OpFetchDimRW:                     "FETCH_DIM_RW",
	OpFetchObjRW:                     "FETCH_OBJ_RW",
	OpFetchIs:                        "FETCH_IS",
	OpFetchDimIs:                     "FETCH_DIM_IS",
	OpFetchObjIs:                     "FETCH_OBJ_IS",
	OpFetchFuncArg:                   "FETCH_FUNC_ARG",
	OpFetchDimFuncArg:                "FETCH_DIM_FUNC_ARG",
	OpFetchObjFuncArg:                "FETCH_OBJ_FUNC_ARG",
	OpFetchUnset:                     "FETCH_UNSET",
	OpFetchDimUnset:                  "FETCH_DIM_UNSET",
	OpFetchObjUnset:                  "FETCH_OBJ_UNSET",
	OpFetchListR:                     "FETCH_LIST_R",
	OpFetchConstant:                  "FETCH_CONSTANT",
	OpCheckFuncArg:                   "CHECK_FUNC_ARG",
	OpExtStmt:                        "EXT_STMT",
	OpExtFcallBegin:                  "EXT_FCALL_BEGIN",
	OpExtFcallEnd:                    "EXT_FCALL_END",
	OpExtNop:                         "EXT_NOP",
	OpTicks:                          "TICKS",
	OpSendVarNoRef:                   "SEND_VAR_NO_REF",
	OpCatch:                          "CATCH",
	OpThrow:                          "THROW",
	OpFetchClass:                     "FETCH_CLASS",
	OpClone:                          "CLONE",
	OpReturnByRef:                    "RETURN_BY_REF",
	OpInitMethodCall:                 "INIT_METHOD_CALL",
	OpInitStaticMethodCall:           "INIT_STATIC_METHOD_CALL",
	OpIssetIsemptyVar:                "ISSET_ISEMPTY_VAR",
	OpIssetIsemptyDimObj:             "ISSET_ISEMPTY_DIM_OBJ",
	OpSendValEx:                      "SEND_VAL_EX",
	OpSendVar:                        "SEND_VAR",
	OpInitUserCall:                   "INIT_USER_CALL",
	OpSendArray:                      "SEND_ARRAY",
	OpSendUser:                       "SEND_USER",
	OpStrlen:                         "STRLEN",
	OpDefined:                        "DEFINED",
	OpTypeCheck:                      "TYPE_CHECK",
	OpVerifyReturnType:               "VERIFY_RETURN_TYPE",
	OpFeResetRW:                      "FE_RESET_RW",
	OpFeFetchRW:                      "FE_FETCH_RW",
	OpFeFree:                         "FE_FREE",
	OpInitDynamicCall:                "INIT_DYNAMIC_CALL",
	OpDoIcall:                        "DO_ICALL",
	OpDoUcall:                        "DO_UCALL",
	OpDoFcallByName:                  "DO_FCALL_BY_NAME",
	OpPreIncObj:                      "PRE_INC_OBJ",
	OpPreDecObj:                      "PRE_DEC_OBJ",
	OpPostIncObj:                     "POST_INC_OBJ",
	OpPostDecObj:                     "POST_DEC_OBJ",
	OpEcho:                           "ECHO",
	OpOpData:                         "OP_DATA",
	OpInstanceof:                     "INSTANCEOF",
	OpGeneratorCreate:                "GENERATOR_CREATE",
	OpMakeRef:                        "MAKE_REF",
	OpDeclareFunction:                "DECLARE_FUNCTION",
	OpDeclareLambdaFunction:          "DECLARE_LAMBDA_FUNCTION",
	OpDeclareConst:                   "DECLARE_CONST",
	OpDeclareClass:                   "DECLARE_CLASS",
	OpDeclareClassDelayed:            "DECLARE_CLASS_DELAYED",
	OpDeclareAnonClass:               "DECLARE_ANON_CLASS",
	OpAddArrayUnpack:                 "ADD_ARRAY_UNPACK",
	OpIssetIsemptyPropObj:            "ISSET_ISEMPTY_PROP_OBJ",
	OpHandleException:                "HANDLE_EXCEPTION",
	OpUserOpcode:                     "USER_OPCODE",
	OpAssertCheck:                    "ASSERT_CHECK",
	OpJmpSet:                         "JMP_SET",
	OpUnsetCV:                        "UNSET_CV",
	OpIssetIsemptyCV:                 "ISSET_ISEMPTY_CV",
	OpFetchListW:                     "FETCH_LIST_W",
	OpSeparate:                       "SEPARATE",
	OpFetchClassName:                 "FETCH_CLASS_NAME",
	OpCallTrampoline:                 "CALL_TRAMPOLINE",
	OpDiscardException:               "DISCARD_EXCEPTION",
	OpYield:                          "YIELD",
	OpGeneratorReturn:                "GENERATOR_RETURN",
	OpFastCall:                       "FAST_CALL",
	OpFastRet:                        "FAST_RET",
	OpRecvVariadic:                   "RECV_VARIADIC",
	OpSendUnpack:                     "SEND_UNPACK",
	OpYieldFrom:                      "YIELD_FROM",
	OpCopyTmp:                        "COPY_TMP",
	OpBindGlobal:                     "BIND_GLOBAL",
	OpCoalesce:                       "COALESCE",
	OpSpaceship:                      "SPACESHIP",
	OpFuncNumArgs:                    "FUNC_NUM_ARGS",
	OpFuncGetArgs:                    "FUNC_GET_ARGS",
	OpFetchStaticPropR:               "FETCH_STATIC_PROP_R",
	OpFetchStaticPropW:               "FETCH_STATIC_PROP_W",
	OpFetchStaticPropRW:              "FETCH_STATIC_PROP_RW",
	OpFetchStaticPropIs:              "FETCH_STATIC_PROP_IS",
	OpFetchStaticPropFuncArg:         "FETCH_STATIC_PROP_FUNC_ARG",
	OpFetchStaticPropUnset:           "FETCH_STATIC_PROP_UNSET",
	OpUnsetStaticProp:                "UNSET_STATIC_PROP",
	OpIssetIsemptyStaticProp:         "ISSET_ISEMPTY_STATIC_PROP",
	OpFetchClassConstant:             "FETCH_CLASS_CONSTANT",
	OpBindLexical:                    "BIND_LEXICAL",
	OpBindStatic:                     "BIND_STATIC",
	OpFetchThis:                      "FETCH_THIS",
	OpSendFuncArg:                    "SEND_FUNC_ARG",
	OpIssetIsemptyThis:               "ISSET_ISEMPTY_THIS",
	OpSwitchLong:                     "SWITCH_LONG",
	OpSwitchString:                   "SWITCH_STRING",
	OpInArray:                        "IN_ARRAY",
	OpCount:                          "COUNT",
	OpGetClass:                       "GET_CLASS",
	OpGetCalledClass:                 "GET_CALLED_CLASS",
	OpGetType:                        "GET_TYPE",
	OpArrayKeyExists:                 "ARRAY_KEY_EXISTS",
	OpMatch:                          "MATCH",
	OpCaseStrict:                     "CASE_STRICT",
	OpMatchError:                     "MATCH_ERROR",
	OpJmpNull:                        "JMP_NULL",
	OpCheckUndefArgs:                 "CHECK_UNDEF_ARGS",
	OpFetchGlobals:                   "FETCH_GLOBALS",
	OpVerifyNeverType:                "VERIFY_NEVER_TYPE",
	OpCallableConvert:                "CALLABLE_CONVERT",
	OpBindInitStaticOrJmp:            "BIND_INIT_STATIC_OR_JMP",
	OpFramelessIcall0:                "FRAMELESS_ICALL_0",
	OpFramelessIcall1:                "FRAMELESS_ICALL_1",
	OpFramelessIcall2:                "FRAMELESS_ICALL_2",
	OpFramelessIcall3:                "FRAMELESS_ICALL_3",
	OpJmpFrameless:                   "JMP_FRAMELESS",
	OpInitParentPropertyHookCall:     "INIT_PARENT_PROPERTY_HOOK_CALL",
	OpDeclareAttributedConst:         "DECLARE_ATTRIBUTED_CONST",
}
