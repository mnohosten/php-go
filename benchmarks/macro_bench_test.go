package benchmarks

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/krizos/php-go/pkg/compiler"
	"github.com/krizos/php-go/pkg/lexer"
	"github.com/krizos/php-go/pkg/parser"
	"github.com/krizos/php-go/pkg/vm"
)

// execPHP executes a PHP script using php-go
func execPHP(t testing.TB, script string) (string, error) {
	t.Helper()

	// Lex and Parse
	l := lexer.New(script, "bench.php")
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return "", &benchError{msg: p.Errors()[0]}
	}

	// Compile
	c := compiler.New()
	if err := c.Compile(program); err != nil {
		return "", err
	}

	bytecode := c.Bytecode()

	// Execute
	v := vm.NewWithBytecode(bytecode.Instructions, bytecode.Constants)

	if err := v.Execute(bytecode.Instructions); err != nil {
		return "", err
	}

	return v.GetOutput(), nil
}

// benchError wraps a string as an error
type benchError struct {
	msg string
}

func (e *benchError) Error() string {
	return e.msg
}

// BenchmarkFullScript benchmarks end-to-end execution of the complete simple_bench.php script
func BenchmarkFullScript(b *testing.B) {
	scriptPath := filepath.Join("scripts", "simple_bench.php")
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		b.Skipf("Skipping full script benchmark: %v", err)
		return
	}

	script := string(content)

	// Warmup
	_, err = execPHP(b, script)
	if err != nil {
		b.Logf("Warning: warmup failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkSimpleLoop benchmarks simple loop operations
func BenchmarkSimpleLoop(b *testing.B) {
	script := `<?php
$a = 0;
for ($i = 0; $i < 10000; $i = $i + 1) {
  $a = $a + 1;
}
echo $a;
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkFunctionCalls benchmarks function call overhead
func BenchmarkFunctionCalls(b *testing.B) {
	script := `<?php
function test($a) {
  return $a;
}

for ($i = 0; $i < 5000; $i = $i + 1) {
  test("hello");
}
echo "OK";
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkRecursion benchmarks recursive function calls (Fibonacci)
func BenchmarkRecursion(b *testing.B) {
	script := `<?php
function fibo($n){
    if ($n < 2) {
        return 1;
    }
    return fibo($n - 2) + fibo($n - 1);
}

echo fibo(15);
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkArrayOperations benchmarks array creation and manipulation
func BenchmarkArrayOperations(b *testing.B) {
	script := `<?php
$X = array();
for ($i=0; $i<1000; $i = $i + 1) {
  $X[$i] = $i;
}
$Y = array();
for ($i=999; $i>=0; $i = $i - 1) {
  $Y[$i] = $X[$i];
}
echo count($Y);
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkHashOperations benchmarks associative array (hash) operations
func BenchmarkHashOperations(b *testing.B) {
	script := `<?php
$X = array();
for ($i = 1; $i <= 1000; $i = $i + 1) {
  $key = "key_" . $i;
  $X[$key] = $i;
}
$c = 0;
for ($i = 1000; $i > 0; $i = $i - 1) {
  $key = "key_" . $i;
  if ($X[$key]) {
    $c = $c + 1;
  }
}
echo $c;
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkMatrixMultiplication benchmarks matrix operations
func BenchmarkMatrixMultiplication(b *testing.B) {
	script := `<?php
function mkmatrix ($rows, $cols) {
    $count = 1;
    $mx = array();
    for ($i=0; $i<$rows; $i = $i + 1) {
      for ($j=0; $j<$cols; $j = $j + 1) {
        $mx[$i][$j] = $count;
        $count = $count + 1;
      }
    }
    return $mx;
}

function mmult ($rows, $cols, $m1, $m2) {
    $m3 = array();
    for ($i=0; $i<$rows; $i = $i + 1) {
      for ($j=0; $j<$cols; $j = $j + 1) {
        $x = 0;
        for ($k=0; $k<$cols; $k = $k + 1) {
          $x = $x + $m1[$i][$k] * $m2[$k][$j];
        }
        $m3[$i][$j] = $x;
      }
    }
    return $m3;
}

$SIZE = 10;
$m1 = mkmatrix($SIZE, $SIZE);
$m2 = mkmatrix($SIZE, $SIZE);
$mm = mmult($SIZE, $SIZE, $m1, $m2);
echo "OK";
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkStringConcatenation benchmarks string operations
func BenchmarkStringConcatenation(b *testing.B) {
	script := `<?php
$str = "";
for ($i = 0; $i < 500; $i = $i + 1) {
  $str = $str . "hello\n";
}
echo strlen($str);
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkNestedLoops benchmarks nested loop performance
func BenchmarkNestedLoops(b *testing.B) {
	script := `<?php
$x = 0;
for ($a=0; $a<10; $a = $a + 1) {
  for ($b=0; $b<10; $b = $b + 1) {
    for ($c=0; $c<10; $c = $c + 1) {
      $x = $x + 1;
    }
  }
}
echo $x;
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkFibonacci benchmarks recursive Fibonacci
func BenchmarkFibonacci(b *testing.B) {
	script := `<?php
function fibo($n){
    if ($n < 2) {
        return 1;
    }
    return fibo($n - 2) + fibo($n - 1);
}

echo fibo(20);
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// BenchmarkMandelbrot benchmarks computational workload
func BenchmarkMandelbrot(b *testing.B) {
	script := `<?php
$w1=10;
$h1=20;
$recen=-0.45;
$imcen=0.0;
$r=0.7;
$s=2*$r/$w1;
$w2=40;
$h2=12;
$count = 0;
for ($y=0; $y<=$w1; $y = $y + 1) {
  $imc=$s*($y-$h2)+$imcen;
  for ($x=0; $x<=$h1; $x = $x + 1) {
    $rec=$s*($x-$w2)+$recen;
    $re=$rec;
    $im=$imc;
    $color=100;
    $re2=$re*$re;
    $im2=$im*$im;
    while ((($re2+$im2)<1000000) && $color>0) {
      $im=$re*$im*2+$imc;
      $re=$re2-$im2+$rec;
      $re2=$re*$re;
      $im2=$im*$im;
      $color = $color - 1;
    }
    if ($color == 0) {
      $count = $count + 1;
    }
  }
}
echo $count;
?>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := execPHP(b, script)
		if err != nil {
			b.Fatalf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// Helper function to compare with actual PHP if available
func ComparePHPPerformance(t *testing.T, script string, name string) {
	phpPath, err := exec.LookPath("php")
	if err != nil {
		t.Skip("PHP not found in PATH, skipping comparison")
	}

	// Write script to temp file
	tmpfile, err := os.CreateTemp("", "bench-*.php")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Run with PHP (3 iterations for average)
	phpTimes := make([]time.Duration, 3)
	for i := 0; i < 3; i++ {
		start := time.Now()
		cmd := exec.Command(phpPath, tmpfile.Name())
		phpOutput, err := cmd.CombinedOutput()
		phpTimes[i] = time.Since(start)
		if err != nil {
			t.Logf("PHP execution failed (iteration %d): %v\nOutput: %s", i, err, phpOutput)
		}
	}

	// Calculate average PHP time
	var avgPHPDuration time.Duration
	for _, d := range phpTimes {
		avgPHPDuration += d
	}
	avgPHPDuration /= time.Duration(len(phpTimes))

	// Run with php-go (3 iterations for average)
	phpgoTimes := make([]time.Duration, 3)
	for i := 0; i < 3; i++ {
		start := time.Now()
		phpgoOutput, err := execPHP(t, script)
		phpgoTimes[i] = time.Since(start)
		if err != nil {
			t.Fatalf("php-go execution failed (iteration %d): %v", i, err)
		}
		if i == 0 {
			t.Logf("php-go output length: %d bytes", len(phpgoOutput))
		}
	}

	// Calculate average php-go time
	var avgPhpgoDuration time.Duration
	for _, d := range phpgoTimes {
		avgPhpgoDuration += d
	}
	avgPhpgoDuration /= time.Duration(len(phpgoTimes))

	t.Logf("=== %s ===", name)
	t.Logf("PHP avg execution time: %v", avgPHPDuration)
	t.Logf("php-go avg execution time: %v", avgPhpgoDuration)

	if avgPHPDuration > 0 {
		ratio := float64(avgPhpgoDuration) / float64(avgPHPDuration)
		t.Logf("php-go/PHP ratio: %.2fx", ratio)
		if ratio < 1.0 {
			t.Logf("✓ php-go is FASTER than PHP by %.1f%%", (1-ratio)*100)
		} else {
			t.Logf("php-go is %.1f%% slower than PHP", (ratio-1)*100)
		}
	}
}

// TestComparison is a test (not benchmark) that compares outputs and performance
func TestComparisonSimple(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comparison test in short mode")
	}

	script := `<?php
function fibo($n){
    if ($n < 2) {
        return 1;
    }
    return fibo($n - 2) + fibo($n - 1);
}
echo fibo(20);
?>`

	ComparePHPPerformance(t, script, "Fibonacci(20)")
}

func TestComparisonArray(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comparison test in short mode")
	}

	script := `<?php
$X = array();
for ($i=0; $i<1000; $i = $i + 1) {
  $X[$i] = $i;
}
echo count($X);
?>`

	ComparePHPPerformance(t, script, "Array Operations")
}
