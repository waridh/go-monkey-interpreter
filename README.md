# go-monkey-interpreter

## Overview

This project is the interpreter for the Monkey Programming language, the spec
of which is written by [Thorsten Ball](https://thorstenball.com/), following
the wonderfully written [Writing An Interpreter In Go](https://interpreterbook.com/).

Monkey is an expression based programming language with support for basic
data types, such as booleans, integers, strings, arrays, and hash maps.
It also treats functions as a first class citizen, allowing for higher
order functions.

## Example Code

The following is an implementation of the Fibonacci sequence.

```monkey
let fib_aux = fn(target, sub_two, sub_one, counter) {
  if (target == counter) {
    sub_two + sub_one
  } else {
    fib_aux(target, sub_one, sub_one + sub_two, counter + 1)
  }
};
let fib = fn(x){
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      1
    } else {
      fib_aux(x, 0, 1, 2)
    }
  }
};
```
