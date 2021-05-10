# CornyLang
`CornyLang` is an expressive programming language built from scrach using a Top [!https://interpreterbook.com/]Down Parser algorithm. In this repository I will host all implementation of `CornyLang` in every programming language I learn.

`CornyLang` has the following features:

- [C-like syntax](#syntax)
- [Variable Bindings](#variable-bindings)
- [Primitive Data Types](#primitive-data-types))
- [Compound Data Types](#compound-data-types)
- [Expressions](#expressions)
- [Built-in funcions](#built-in-funcions)
- [First Class and High-order Functions (Inspired from Lisp)](#first-class-high-order-functions)
- [Closures (Inpired from List)](#closures)
- [Control Flow](#control-flow)

`CornyLang` disadvantages

- No comment support
- No iterator control flow (while, for, etc)

## C-like syntax

```Javascript
>>> let a = 10;
>>> let b = 20;
>>> let add = fn(n1, n2) { return n1 + n2 };
>>> add(a, b); // 30
```

## Variable Bindings

Bindings are made by using the `let` keyword followed by any valid expression.

```Javascript
>>> let integer = 10; // integer variable
>>> let string = "John"; // string variable
>>> let boolean = true; // boolean variable
>>> let fruits = ["Apple", "Grappe", "Water Melon"]; // array variable
>>> let json = {"name": "John", "last": "Doe"}; // hash or dictionary variable
>>> let greets = fn(name) { return "Hello, " + name }; // function variable
```

## Primitive Data Types

```Javascript
>>> let year = 1985; // integer variable
>>> let name = "Irwin"; // string variable
>>> let married = true; // boolean variable
>>> let country = null; // null variable
```

## Compound Data Types

```Javascript
>>> let hobbies = ["Play Guitar", "Game Development", "Study Chinese"]; // array variable
>>> let info = {"name": "Irwin", "sports": ["Basket", "Skate"]}; // hash or dictionary variable
```

## Expressions

#### Arithmetic Expressions
```Javascript
>>> 5 + 2
>>> 5 - 2
>>> 5 * 2
>>> 5 / 2
>>> (5 + 2) * 3
```

#### Relational Operators
```Javascript
>>> 5 < 2
>>> 5 <= 2
>>> 5 > 2
>>> 5 >= 2
>>> ((5 + 2) * 3) == 25
>>> 25 != 52
```

#### Logical Operators *(Operands must be resolved to boolean types)*
```Javascript
>>> true and false
>>> true or false
>>> 5 > 2 and 3 < 4
```

#### Prefix Operators *(Minus '-' and Not '!' operators)*
```Javascript
>>> -5
>>> !true
>>> ---5
>>> !!!false
```

#### Assignment Operator

If cases when your left side variable is repeated inmediately on the right side, you can write it with this compressed form:

```Javascript
>>> let a = 10;
>>> a = a + 5; // repeated twice then...
>>> a += 5; // use this format for the same result
>>> a -= 2;
>>> a *= 3;
>>> a /= 1;
```

## Built-in funcions
- Len: get the size of any string, array or dictionary
    ```Javascript
    >>> let text = "This is a sample text";
    >>> len(text); // 21
    ```
- Type: returns the type of the argument passed.
    ```Javascript
    >>> let text = "This is a sample text";
    >>> type(text); // "string"
    ```
- First: returns an array with the first element of the array passed.
    ```Javascript
    >>> let hobbies = ["Play Guitar", "Game Development", "Study Chinese"];
    >>> first(hobbies); // ["Play Guitar"]
    ```
- Last: returns an array with the last element of the array passed.
    ```Javascript
    >>> let hobbies = ["Play Guitar", "Game Development", "Study Chinese"];
    >>> last(hobbies); // ["Study Chinese"]
    ```
- Rest: returns an array with the `tail` of the array passed. (2..n avoid the head)
    ```Javascript
    >>> let hobbies = ["Play Guitar", "Game Development", "Study Chinese"];
    >>> rest(hobbies); // ["Game Development", "Study Chinese"]
    ```

## First Class and High-order Functions (Inspired from Lisp)
If you come from *imperative* languages like C then you must be familiar with function declaration which must be written somewhere in your code and then get called. `CornyLang` will not let you write static functions because they are `First Class Citizens` which means they are treated like any other variable that you can bind somewhere and take it with you.

```Javascript
>>> let square = fn(a) { return a * a };
>>> square(5); // 25
```

Functions in `CornyLang` are also in `High Order`. These are functions that take other functions as arguments. Please refer to the following example and take your time to understand it.

```Javascript
>>> let square = fn(base) { return base * base };
>>> let cube = fn(base, squareFunc) { return squareFunc(base) * base; };
>>> cube(3, square); // 27
```

The `cube()` function takes two arguments: the base and another function called `square()` who gets called inside of the function body. This is done because in `CornyLang` functions are just values like strings or booleans and that feature exists in `Functional Programming Languages` lile `Lisp` or `Haskell`.

## Closures (Inpired from List)
`Closures` or `Lexical Scope` is another feature from `Functional Programming Languages` *(yep! I was inspired from functional programming languages)* that allow write a function inside another function.

Let's see another solution for the `cube()` function example. This time with `Closures`

#### Closure # 1: using binding and calling the closure separately
```Javascript
>>> let cube = fn(base) { let square = fn(base) { return base * base; }; return square(base) * base; };
>>> cube(3); // 27
```

#### Closure # 2: using binding and calling the closure at a time
```Javascript
>>> let cube = fn(base) { let square = fn(base) { return base * base; }(base) * base; };
>>> cube(3); // 27
```

#### Closure # 3: using anonymous function
```Javascript
>>> let cube = fn(base) { fn(base) { return base * base; }(base) * base; };
>>> cube(3); // 27
```

## Control Flow
#### If Expression
In `Imperative` languages the `IF` is a `Statement` but in `CornyLang` it is an `Expression` too but it usage is pretty similar:

```Javascript
>>> let lowest = fn(x, y) { if (x < y) { return x; } else { return y; }; };
>>> lowest(5, 10); // 5
```

## License

`CornyLang` is released under the MIT Licence.
