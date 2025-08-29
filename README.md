# Pirat Programming Language Interpreter

## What is this?

Pirat is a fully interpreted pirate-themed programming language with syntax loosely inspired by python. This implementation uses a recursive top-down parser (Pratt), and tree-walking evaluator. And before you ask, there are no other implementations and the language itself is not very well defined.

## Features

#### Functions
```
yar greeting be f(x):
  gives 'Hello ' + x + "!".
.
greeting('world!').
```

#### For (4) loops
```
yar i be 0.
4 i < 10:
  if i = 6:
    break.
  .
  i be i + 1
.
```

#### Arrays
```
yar arrrrr be [1, "2", (1+2), [1, 2, 3]].
arrrrr[0] be 2.
```

#### Hash maps
```
yar map be {"key1": 1, "key2": 2}.
map["key3"] be "three".
```

#### Control flow
```
if nay <> ay:
  ahoy("false is not equal to true").
lsif i = 1:
  ahoy("i equals 1").
ls:
  ahoy("Else block hit").
.
```

#### Chests (structs)
```
chest myChestType|foo, bar|.
yar instance be myChestType|"fooVal", f(): gives "barval"..|.
instance|foo be instance|bar().
```

## How to run locally (assuming you are not using the release executables)
You should have golang and make installed

For repl: `make repl`

With a pir file: `make run FILE=<YOUR_PIR_FILE>`
