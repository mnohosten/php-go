package parser

import (
	"testing"

	"github.com/krizos/php-go/pkg/lexer"
)

// BenchmarkSimpleExpression benchmarks parsing simple expressions
func BenchmarkSimpleExpression(b *testing.B) {
	input := `<?php
$x = 10;
$y = 20;
$z = $x + $y;`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkComplexExpression benchmarks parsing complex expressions
func BenchmarkComplexExpression(b *testing.B) {
	input := `<?php
$result = (($a + $b) * ($c - $d)) / (($e % $f) ** ($g & $h));
$comparison = ($x == $y) && ($a != $b) || ($c < $d) && ($e > $f);
$ternary = ($condition) ? $trueValue : $falseValue;`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkArrayLiteral benchmarks parsing array literals
func BenchmarkArrayLiteral(b *testing.B) {
	input := `<?php
$simple = [1, 2, 3, 4, 5];
$assoc = ["name" => "John", "age" => 30, "city" => "New York"];
$nested = [
	[1, 2, 3],
	[4, 5, 6],
	[7, 8, 9]
];`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkIfStatement benchmarks parsing if statements
func BenchmarkIfStatement(b *testing.B) {
	input := `<?php
if ($x > 0) {
	echo "positive";
} elseif ($x < 0) {
	echo "negative";
} else {
	echo "zero";
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkForLoop benchmarks parsing for loops
func BenchmarkForLoop(b *testing.B) {
	input := `<?php
for ($i = 0; $i < 10; $i++) {
	echo $i;
	for ($j = 0; $j < 5; $j++) {
		echo $j;
	}
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkWhileLoop benchmarks parsing while loops
func BenchmarkWhileLoop(b *testing.B) {
	input := `<?php
while ($x > 0) {
	$x--;
	echo $x;
}

do {
	$y++;
	echo $y;
} while ($y < 100);`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkSwitchStatement benchmarks parsing switch statements
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkTryCatchFinally benchmarks parsing try-catch-finally
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkFunctionDeclaration benchmarks parsing function declarations
func BenchmarkFunctionDeclaration(b *testing.B) {
	input := `<?php
function calculateTotal(
	int $price,
	float $tax,
	?string $coupon = null
): float {
	$subtotal = $price * 1.0;
	if ($coupon !== null) {
		$subtotal = $subtotal * 0.9;
	}
	return $subtotal + ($subtotal * $tax);
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkSimpleClass benchmarks parsing simple class declarations
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkComplexClass benchmarks parsing complex class declarations
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

	public function show($id) {
		$user = $this->userService->find($id);
		if ($user === null) {
			throw new Exception("User not found");
		}
		return $this->render("users/show", $user);
	}
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkInterface benchmarks parsing interface declarations
func BenchmarkInterface(b *testing.B) {
	input := `<?php
interface Repository {
	public function find($id);
	public function findAll();
	public function save($entity);
	public function delete($entity);
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkTrait benchmarks parsing trait declarations
func BenchmarkTrait(b *testing.B) {
	input := `<?php
trait Timestamped {
	private $createdAt;
	private $updatedAt;

	public function setCreatedAt($timestamp) {
		$this->createdAt = $timestamp;
	}

	public function setUpdatedAt($timestamp) {
		$this->updatedAt = $timestamp;
	}

	public function getCreatedAt() {
		return $this->createdAt;
	}
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkMethodChaining benchmarks parsing method chaining
func BenchmarkMethodChaining(b *testing.B) {
	input := `<?php
$result = $query
	->select("*")
	->from("users")
	->where("active", true)
	->orderBy("name")
	->limit(10)
	->get();`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkLargeFile benchmarks parsing a large realistic PHP file
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkTypeHints benchmarks parsing type hints
func BenchmarkTypeHints(b *testing.B) {
	input := `<?php
function process(
	int $id,
	string $name,
	?array $options,
	callable $callback
): bool {
	return true;
}

function getResult(): int|string|null {
	return 42;
}

function nullable(?User $user): ?string {
	return $user->getName();
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input, "bench.php")
		p := New(l)
		p.ParseProgram()
	}
}
