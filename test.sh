#!/bin/bash
assert() {
  expected="$1"
  input="$2"

  ./9cc "$input" > tmp.s
  cc -o tmp tmp.s
  ./tmp
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

assert 0 '0;'
assert 1 '1;'
assert 6 ' 1 + 2 + 3; '
assert 7 '1+2*3;'
assert 3 '3*(2-1);'
assert 3 '(4+3)/2;'
assert 1 '-1+2;'
assert 0 '0 == 1;'
assert 1 '0 == 0;'
assert 0 '0 != 0;'
assert 1 '0 != 1;'
assert 0 '1 <= 0;'
assert 1 '1 <= 1;'
assert 0 '0 >= 1;'
assert 1 '1 >= 1;'
assert 2 'a = 2;'
assert 3 'a = 2; b = 3;'
assert 5 'a = 2; b = a + 3;'
assert 6 'a = b = c = (-2 + 3) * 2; d = a * b + c;'
assert 2 'abc = 2;'

echo OK