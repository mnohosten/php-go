package compiler

import (
	"testing"

	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
)

// BenchmarkSimpleExpression benchmarks compiling simple expressions
func BenchmarkSimpleExpression(b *testing.B) {
	input := `<?php
$x = 10;
$y = 20;
$z = $x + $y;`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkArithmeticOperations benchmarks compiling arithmetic operations
func BenchmarkArithmeticOperations(b *testing.B) {
	input := `<?php
$a = 1 + 2;
$b = 3 - 4;
$c = 5 * 6;
$d = 7 / 8;
$e = 9 % 10;
$f = 2 ** 3;`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkComplexExpression benchmarks compiling complex expressions
func BenchmarkComplexExpression(b *testing.B) {
	input := `<?php
$result = (($a + $b) * ($c - $d)) / (($e % $f) ** ($g & $h));
$comparison = ($x == $y) && ($a != $b) || ($c < $d) && ($e > $f);
$ternary = ($condition) ? $trueValue : $falseValue;`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkArrayLiteral benchmarks compiling array literals
func BenchmarkArrayLiteral(b *testing.B) {
	input := `<?php
$simple = [1, 2, 3, 4, 5];
$assoc = ["name" => "John", "age" => 30, "city" => "New York"];
$nested = [
	[1, 2, 3],
	[4, 5, 6],
	[7, 8, 9]
];`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkIfStatement benchmarks compiling if statements
func BenchmarkIfStatement(b *testing.B) {
	input := `<?php
if ($x > 0) {
	echo "positive";
} elseif ($x < 0) {
	echo "negative";
} else {
	echo "zero";
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkForLoop benchmarks compiling for loops
func BenchmarkForLoop(b *testing.B) {
	input := `<?php
for ($i = 0; $i < 10; $i++) {
	echo $i;
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkWhileLoop benchmarks compiling while loops
func BenchmarkWhileLoop(b *testing.B) {
	input := `<?php
while ($x > 0) {
	$x--;
	echo $x;
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkForeachLoop benchmarks compiling foreach loops
func BenchmarkForeachLoop(b *testing.B) {
	input := `<?php
foreach ($array as $key => $value) {
	echo $key . ": " . $value;
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkSwitchStatement benchmarks compiling switch statements
func BenchmarkSwitchStatement(b *testing.B) {
	input := `<?php
switch ($value) {
	case 1:
		echo "one";
		break;
	case 2:
		echo "two";
		break;
	case 3:
		echo "three";
		break;
	default:
		echo "other";
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkTryCatchFinally benchmarks compiling try-catch-finally
func BenchmarkTryCatchFinally(b *testing.B) {
	input := `<?php
try {
	$result = riskyOperation();
	echo $result;
} catch (Exception $e) {
	echo $e->getMessage();
} finally {
	cleanup();
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkFunctionDeclaration benchmarks compiling function declarations
func BenchmarkFunctionDeclaration(b *testing.B) {
	input := `<?php
function calculateTotal(int $price, float $tax, ?string $coupon = null): float {
	$subtotal = $price * 1.0;
	if ($coupon !== null) {
		$subtotal = $subtotal * 0.9;
	}
	return $subtotal + ($subtotal * $tax);
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkFunctionCall benchmarks compiling function calls
func BenchmarkFunctionCall(b *testing.B) {
	input := `<?php
$result = strlen($str);
$formatted = sprintf("%s: %d", $name, $count);
$filtered = array_filter($items, $callback);`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkSimpleClass benchmarks compiling simple class declarations
func BenchmarkSimpleClass(b *testing.B) {
	input := `<?php
class User {
	private $id;
	private $name;
	private $email;

	public function __construct($id, $name, $email) {
		$this->id = $id;
		$this->name = $name;
		$this->email = $email;
	}

	public function getId() {
		return $this->id;
	}
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkComplexClass benchmarks compiling complex class declarations
func BenchmarkComplexClass(b *testing.B) {
	input := `<?php
abstract class BaseController {
	protected $db;
	protected $cache;

	public function __construct($db) {
		$this->db = $db;
		$this->cache = [];
	}

	abstract public function index();

	protected function render($view, $data) {
		return $view . $data;
	}
}

final class UserController extends BaseController {
	private $userService;

	public function __construct($db, $userService) {
		parent::__construct($db);
		$this->userService = $userService;
	}

	public function index() {
		$users = $this->userService->getAll();
		return $this->render("users/index", $users);
	}
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkObjectInstantiation benchmarks compiling object instantiation
func BenchmarkObjectInstantiation(b *testing.B) {
	input := `<?php
$user = new User($id, $name, $email);
$controller = new UserController($db, $service);
$instance = new $className($arg1, $arg2);`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkPropertyAccess benchmarks compiling property access
func BenchmarkPropertyAccess(b *testing.B) {
	input := `<?php
$name = $user->name;
$user->email = "test@example.com";
$value = $object->property->nested->deep;`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkMethodCall benchmarks compiling method calls
func BenchmarkMethodCall(b *testing.B) {
	input := `<?php
$result = $user->getName();
$user->setEmail("test@example.com");
$chain = $query->select("*")->where("id", 1)->get();`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkStaticAccess benchmarks compiling static access
func BenchmarkStaticAccess(b *testing.B) {
	input := `<?php
$value = User::$staticProperty;
User::$staticProperty = $newValue;
$result = User::staticMethod($arg);
$className = User::class;`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkArrayAccess benchmarks compiling array access
func BenchmarkArrayAccess(b *testing.B) {
	input := `<?php
$item = $array[$key];
$array[$key] = $value;
$nested = $data["user"]["profile"]["name"];`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkClosure benchmarks compiling closures
func BenchmarkClosure(b *testing.B) {
	input := `<?php
$add = function($a, $b) {
	return $a + $b;
};

$multiplier = function($factor) use ($base) {
	return function($x) use ($factor, $base) {
		return $x * $factor + $base;
	};
};`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkLargeFile benchmarks compiling a large realistic PHP file
func BenchmarkLargeFile(b *testing.B) {
	input := `<?php
class OrderProcessor {
	private $db;
	private $logger;
	private $cache;

	public function __construct($db, $logger) {
		$this->db = $db;
		$this->logger = $logger;
		$this->cache = [];
	}

	public function processOrder($orderId) {
		try {
			$order = $this->getOrder($orderId);
			if ($order === null) {
				throw new Exception("Order not found");
			}

			$items = $order->getItems();
			$total = 0.0;

			foreach ($items as $item) {
				$price = $item->getPrice();
				$quantity = $item->getQuantity();

				if ($quantity <= 0) {
					throw new Exception("Invalid quantity");
				}

				$subtotal = $price * $quantity;

				if ($item->hasDiscount()) {
					$discount = $item->getDiscount();
					$subtotal = $subtotal - ($subtotal * $discount / 100);
				}

				$total = $total + $subtotal;
			}

			$tax = $this->calculateTax($total);
			$shipping = $this->calculateShipping($order);
			$grandTotal = $total + $tax + $shipping;

			$order->setSubtotal($total);
			$order->setTax($tax);
			$order->setShipping($shipping);
			$order->setTotal($grandTotal);
			$order->setStatus("processed");
			$order->setProcessedAt(time());

			$this->db->beginTransaction();
			$this->db->save($order);
			$this->db->commit();

			$this->cache[$orderId] = $order;
			$this->logger->info("Order processed successfully", [
				"order_id" => $orderId,
				"total" => $grandTotal
			]);

			return true;
		} catch (Exception $e) {
			$this->db->rollback();
			$this->logger->error("Failed to process order", [
				"order_id" => $orderId,
				"error" => $e->getMessage(),
				"trace" => $e->getTraceAsString()
			]);
			return false;
		}
	}

	private function getOrder($id) {
		if ($this->cache[$id]) {
			return $this->cache[$id];
		}
		$order = $this->db->find($id);
		$this->cache[$id] = $order;
		return $order;
	}

	private function calculateTax($subtotal) {
		$rate = 0.1;
		return $subtotal * $rate;
	}

	private function calculateShipping($order) {
		$weight = 0;
		foreach ($order->getItems() as $item) {
			$weight = $weight + $item->getWeight();
		}

		if ($weight < 1) {
			return 5.0;
		} elseif ($weight < 5) {
			return 10.0;
		} elseif ($weight < 10) {
			return 15.0;
		} else {
			return 20.0;
		}
	}
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkConstantFolding benchmarks constant folding optimization
func BenchmarkConstantFolding(b *testing.B) {
	input := `<?php
$a = 2 + 3;
$b = 10 * 5;
$c = 100 / 2;
$d = 15 - 5;
$e = "Hello" . " " . "World";`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkSymbolTableOperations benchmarks symbol table operations
func BenchmarkSymbolTableOperations(b *testing.B) {
	input := `<?php
$var1 = 1;
$var2 = 2;
$var3 = 3;
$var4 = 4;
$var5 = 5;
$result = $var1 + $var2 + $var3 + $var4 + $var5;`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkConstantTableOperations benchmarks constant table operations
func BenchmarkConstantTableOperations(b *testing.B) {
	input := `<?php
$a = 123;
$b = 456;
$c = "test";
$d = "example";
$e = 3.14;
$f = 2.71;
$g = true;
$h = false;
$i = null;`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}

// BenchmarkJumpPatching benchmarks jump patching in control flow
func BenchmarkJumpPatching(b *testing.B) {
	input := `<?php
if ($a > 0) {
	if ($b > 0) {
		if ($c > 0) {
			echo "all positive";
		}
	}
}

for ($i = 0; $i < 10; $i++) {
	if ($i % 2 == 0) {
		continue;
	}
	echo $i;
}`

	l := lexer.New(input, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := New()
		c.Compile(program)
	}
}
