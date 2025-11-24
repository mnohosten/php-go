<?php
// Macro-benchmark suite for php-go
// Based on Zend bench.php - simplified for php-go compatibility

function simple() {
  $a = 0;
  for ($i = 0; $i < 1000000; $i++)
    $a++;

  $thisisanotherlongname = 0;
  for ($thisisalongname = 0; $thisisalongname < 1000000; $thisisalongname++)
    $thisisanotherlongname++;
}

function simplecall() {
  for ($i = 0; $i < 1000000; $i++)
    strlen("hallo");
}

function hallo($a) {
}

function simpleucall() {
  for ($i = 0; $i < 1000000; $i++)
    hallo("hallo");
}

function simpleudcall() {
  for ($i = 0; $i < 1000000; $i++)
    hallo2("hallo");
}

function hallo2($a) {
}

function mandel() {
  $w1=50;
  $h1=150;
  $recen=-.45;
  $imcen=0.0;
  $r=0.7;
  $s=0;  $rec=0;  $imc=0;  $re=0;  $im=0;  $re2=0;  $im2=0;
  $x=0;  $y=0;  $w2=0;  $h2=0;  $color=0;
  $s=2*$r/$w1;
  $w2=40;
  $h2=12;
  for ($y=0 ; $y<=$w1; $y=$y+1) {
    $imc=$s*($y-$h2)+$imcen;
    for ($x=0 ; $x<=$h1; $x=$x+1) {
      $rec=$s*($x-$w2)+$recen;
      $re=$rec;
      $im=$imc;
      $color=1000;
      $re2=$re*$re;
      $im2=$im*$im;
      while( ((($re2+$im2)<1000000) && $color>0)) {
        $im=$re*$im*2+$imc;
        $re=$re2-$im2+$rec;
        $re2=$re*$re;
        $im2=$im*$im;
        $color=$color-1;
      }
    }
  }
}

function Ack($m, $n){
  if($m == 0) return $n+1;
  if($n == 0) return Ack($m-1, 1);
  return Ack($m - 1, Ack($m, ($n - 1)));
}

function ackermann($n) {
  $r = Ack(3,$n);
  echo "Ack(3,$n): $r\n";
}

function ary($n) {
  for ($i=0; $i<$n; $i++) {
    $X[$i] = $i;
  }
  for ($i=$n-1; $i>=0; $i--) {
    $Y[$i] = $X[$i];
  }
  $last = $n-1;
  echo "$Y[$last]\n";
}

function ary2($n) {
  for ($i=0; $i<$n;) {
    $X[$i] = $i; ++$i;
    $X[$i] = $i; ++$i;
    $X[$i] = $i; ++$i;
    $X[$i] = $i; ++$i;
    $X[$i] = $i; ++$i;

    $X[$i] = $i; ++$i;
    $X[$i] = $i; ++$i;
    $X[$i] = $i; ++$i;
    $X[$i] = $i; ++$i;
    $X[$i] = $i; ++$i;
  }
  for ($i=$n-1; $i>=0;) {
    $Y[$i] = $X[$i]; --$i;
    $Y[$i] = $X[$i]; --$i;
    $Y[$i] = $X[$i]; --$i;
    $Y[$i] = $X[$i]; --$i;
    $Y[$i] = $X[$i]; --$i;

    $Y[$i] = $X[$i]; --$i;
    $Y[$i] = $X[$i]; --$i;
    $Y[$i] = $X[$i]; --$i;
    $Y[$i] = $X[$i]; --$i;
    $Y[$i] = $X[$i]; --$i;
  }
  $last = $n-1;
  echo "$Y[$last]\n";
}

function ary3($n) {
  for ($i=0; $i<$n; $i++) {
    $X[$i] = $i + 1;
    $Y[$i] = 0;
  }
  for ($k=0; $k<1000; $k++) {
    for ($i=$n-1; $i>=0; $i--) {
      $Y[$i] += $X[$i];
    }
  }
  $last = $n-1;
  echo "$Y[0] $Y[$last]\n";
}

function fibo_r($n){
    return(($n < 2) ? 1 : fibo_r($n - 2) + fibo_r($n - 1));
}

function fibo($n) {
  $r = fibo_r($n);
  echo "$r\n";
}

function hash1($n) {
  for ($i = 1; $i <= $n; $i++) {
    $X[dechex($i)] = $i;
  }
  $c = 0;
  for ($i = $n; $i > 0; $i--) {
    if ($X[dechex($i)]) { $c++; }
  }
  echo "$c\n";
}

function hash2($n) {
  for ($i = 0; $i < $n; $i++) {
    $hash1["foo_$i"] = $i;
    $hash2["foo_$i"] = 0;
  }
  for ($i = $n; $i > 0; $i--) {
    foreach($hash1 as $key => $value) $hash2[$key] += $value;
  }
  $first = "foo_0";
  $last  = "foo_".($n-1);
  echo "$hash1[$first] $hash1[$last] $hash2[$first] $hash2[$last]\n";
}

function mkmatrix ($rows, $cols) {
    $count = 1;
    $mx = array();
    for ($i=0; $i<$rows; $i++) {
    for ($j=0; $j<$cols; $j++) {
        $mx[$i][$j] = $count++;
    }
    }
    return($mx);
}

function mmult ($rows, $cols, $m1, $m2) {
    $m3 = array();
    for ($i=0; $i<$rows; $i++) {
    for ($j=0; $j<$cols; $j++) {
        $x = 0;
        for ($k=0; $k<$cols; $k++) {
        $x += $m1[$i][$k] * $m2[$k][$j];
        }
        $m3[$i][$j] = $x;
    }
    }
    return($m3);
}

function matrix($n) {
  $SIZE = 30;
  $m1 = mkmatrix($SIZE, $SIZE);
  $m2 = mkmatrix($SIZE, $SIZE);
  while ($n--) {
    $mm = mmult($SIZE, $SIZE, $m1, $m2);
  }
  echo "{$mm[0][0]} {$mm[2][3]} {$mm[3][2]} {$mm[4][4]}\n";
}

function nestedloop($n) {
  $x = 0;
  for ($a=0; $a<$n; $a++)
    for ($b=0; $b<$n; $b++)
      for ($c=0; $c<$n; $c++)
        for ($d=0; $d<$n; $d++)
          for ($e=0; $e<$n; $e++)
            for ($f=0; $f<$n; $f++)
             $x++;
  echo "$x\n";
}

function strcat($n) {
  $str = "";
  while ($n-- > 0) {
    $str .= "hello\n";
  }
  $len = strlen($str);
  echo "$len\n";
}

// Run benchmarks
echo "Running macro-benchmarks...\n";

simple();
echo "simple: OK\n";

simplecall();
echo "simplecall: OK\n";

simpleucall();
echo "simpleucall: OK\n";

simpleudcall();
echo "simpleudcall: OK\n";

mandel();
echo "mandel: OK\n";

ackermann(7);
echo "ackermann(7): OK\n";

ary(50000);
echo "ary(50000): OK\n";

ary2(50000);
echo "ary2(50000): OK\n";

ary3(2000);
echo "ary3(2000): OK\n";

fibo(30);
echo "fibo(30): OK\n";

hash1(50000);
echo "hash1(50000): OK\n";

hash2(500);
echo "hash2(500): OK\n";

matrix(20);
echo "matrix(20): OK\n";

nestedloop(12);
echo "nestedloop(12): OK\n";

strcat(200000);
echo "strcat(200000): OK\n";

echo "All benchmarks completed!\n";
?>
