<?php
// Simplified macro-benchmarks for php-go (using supported syntax)

// Simple loop benchmark
function simple_loop() {
  $a = 0;
  for ($i = 0; $i < 10000; $i = $i + 1) {
    $a = $a + 1;
  }
  return $a;
}

// Function call benchmark
function hallo($a) {
  return $a;
}

function function_calls() {
  for ($i = 0; $i < 10000; $i = $i + 1) {
    hallo("test");
  }
}

// Recursive Fibonacci
function fibo($n){
    if ($n < 2) {
        return 1;
    }
    return fibo($n - 2) + fibo($n - 1);
}

// Array operations
function array_ops($n) {
  $X = array();
  for ($i=0; $i<$n; $i = $i + 1) {
    $X[$i] = $i;
  }
  $Y = array();
  for ($i=$n-1; $i>=0; $i = $i - 1) {
    $Y[$i] = $X[$i];
  }
  return $Y;
}

// Hash/associative array operations
function hash_ops($n) {
  $X = array();
  for ($i = 1; $i <= $n; $i = $i + 1) {
    $key = "key_" . $i;
    $X[$key] = $i;
  }
  $c = 0;
  for ($i = $n; $i > 0; $i = $i - 1) {
    $key = "key_" . $i;
    if ($X[$key]) {
      $c = $c + 1;
    }
  }
  return $c;
}

// String concatenation
function string_concat($n) {
  $str = "";
  for ($i = 0; $i < $n; $i = $i + 1) {
    $str = $str . "hello\n";
  }
  return strlen($str);
}

// Matrix operations
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

function matrix($n) {
  $SIZE = 10;
  $m1 = mkmatrix($SIZE, $SIZE);
  $m2 = mkmatrix($SIZE, $SIZE);
  $mm = array();
  for ($i = 0; $i < $n; $i = $i + 1) {
    $mm = mmult($SIZE, $SIZE, $m1, $m2);
  }
  return $mm;
}

// Nested loops
function nested_loops($n) {
  $x = 0;
  for ($a=0; $a<$n; $a = $a + 1) {
    for ($b=0; $b<$n; $b = $b + 1) {
      for ($c=0; $c<$n; $c = $c + 1) {
        $x = $x + 1;
      }
    }
  }
  return $x;
}

// Mandelbrot-like computation
function mandel_simple() {
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
  return $count;
}

// Run all benchmarks
echo "=== PHP-Go Macro Benchmarks ===\n";

$result = simple_loop();
echo "simple_loop: " . $result . "\n";

function_calls();
echo "function_calls: OK\n";

$result = fibo(15);
echo "fibo(15): " . $result . "\n";

$result = array_ops(1000);
echo "array_ops(1000): " . count($result) . " elements\n";

$result = hash_ops(1000);
echo "hash_ops(1000): " . $result . "\n";

$result = string_concat(100);
echo "string_concat(100): " . $result . " bytes\n";

$result = matrix(3);
echo "matrix(3): OK\n";

$result = nested_loops(10);
echo "nested_loops(10): " . $result . "\n";

$result = mandel_simple();
echo "mandel_simple: " . $result . "\n";

echo "=== All benchmarks completed ===\n";
?>
