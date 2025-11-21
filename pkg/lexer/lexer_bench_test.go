package lexer

import (
	"testing"
)

// BenchmarkSimpleTokens benchmarks tokenizing simple PHP code
func BenchmarkSimpleTokens(b *testing.B) {
	input := `<?php
$x = 10;
$y = 20;
$z = $x + $y;
echo $z;`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkVariablesAndOperators benchmarks tokenizing variables and operators
func BenchmarkVariablesAndOperators(b *testing.B) {
	input := `<?php
$a = 1 + 2 - 3 * 4 / 5 % 6;
$b = $a << 2 >> 1 & 0xFF | 0x0F ^ 0xF0;
$c = $a == $b || $a != $b && $a < $b;
$d = ($a === $b) ? $a : $b;
$e++; $f--; ++$g; --$h;`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkStringLiterals benchmarks tokenizing string literals
func BenchmarkStringLiterals(b *testing.B) {
	input := `<?php
$str1 = "Hello, World!";
$str2 = 'Single quoted string';
$str3 = "String with $variable interpolation";
$str4 = "String with escape sequences\n\t\r\\";
$str5 = "String with hex \x48\x65\x6C\x6C\x6F";`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkHeredocTokenization benchmarks tokenizing heredoc strings
func BenchmarkHeredocTokenization(b *testing.B) {
	input := `<?php
$text = <<<EOT
This is a heredoc string
with multiple lines
and some content
that spans several rows
for benchmarking purposes
EOT;`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkComplexExpression benchmarks tokenizing complex expressions
func BenchmarkComplexExpression(b *testing.B) {
	input := `<?php
$result = (($a + $b) * ($c - $d)) / (($e % $f) ** ($g & $h));
$array = ["key1" => $value1, "key2" => $value2, "key3" => $value3];
$object->property = $object->method($arg1, $arg2, $arg3);
$class::$staticProp = $class::staticMethod($x, $y, $z);`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkFunctionDeclaration benchmarks tokenizing function declarations
func BenchmarkFunctionDeclaration(b *testing.B) {
	input := `<?php
function calculateTotal(int $price, float $tax, ?string $coupon = null): float {
	$subtotal = $price * 1.0;
	if ($coupon !== null) {
		$subtotal = $subtotal * 0.9;
	}
	return $subtotal + ($subtotal * $tax);
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkClassDeclaration benchmarks tokenizing class declarations
func BenchmarkClassDeclaration(b *testing.B) {
	input := `<?php
class UserController {
	private $database;
	private $cache;

	public function __construct(Database $db) {
		$this->database = $db;
		$this->cache = [];
	}

	public function getUser(int $id): ?User {
		return $this->database->find($id);
	}
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkControlFlow benchmarks tokenizing control flow statements
func BenchmarkControlFlow(b *testing.B) {
	input := `<?php
if ($condition) {
	for ($i = 0; $i < 10; $i++) {
		while ($x > 0) {
			switch ($value) {
				case 1:
					echo "one";
					break;
				case 2:
					echo "two";
					break;
				default:
					echo "other";
			}
			$x--;
		}
	}
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkLargeFile benchmarks tokenizing a larger PHP file
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
				$subtotal = $price * $quantity;
				$total = $total + $subtotal;
			}

			$tax = $total * 0.1;
			$grandTotal = $total + $tax;

			$order->setTotal($grandTotal);
			$order->setStatus("processed");

			$this->db->save($order);
			$this->logger->info("Order processed", ["id" => $orderId]);

			return true;
		} catch (Exception $e) {
			$this->logger->error("Failed to process order", [
				"id" => $orderId,
				"error" => $e->getMessage()
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
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkNumericLiterals benchmarks tokenizing various numeric formats
func BenchmarkNumericLiterals(b *testing.B) {
	input := `<?php
$dec = 12345;
$hex = 0x1A2B3C;
$oct = 0o777;
$bin = 0b11010101;
$float1 = 123.456;
$float2 = 1.23e4;
$float3 = 1.23e-4;`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkComments benchmarks tokenizing code with comments
func BenchmarkComments(b *testing.B) {
	input := `<?php
// Single line comment
$x = 10; // Inline comment

/* Multi-line comment
 * spanning multiple lines
 * with asterisks
 */
$y = 20;

/**
 * DocBlock comment
 * @param int $value
 * @return string
 */
function test($value) {
	return (string)$value;
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input, "bench.php")
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}
