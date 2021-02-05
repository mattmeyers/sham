# sham

`sham` is a random data generator that uses a custom DSL for defining the data's shape. 

## Install

```
go get -u github.com/mattmeyers/sham/cmd/sham/...
```

## Usage

```
sham is a tool for generating random data

Usage:

    sham [options] <schema>

Options:

    -pretty     pretty print the result
```

To ensure the schema is not affected by any shell escaping, it is recommended that the schema be surrounded by single quotes.

### Example

The following schema

```
{
    "name": name,
    "friends": [(1,5),
        {
            "name": name,
            "age": (20,30)
        }
    ]
}
```

will produce data such as

```
{
    "name": "John Doe",
    "friends": [
        {
            "name: "Bob Smith",
            "age": 21
        },
        {
            "name: "Matt Doe",
            "age": 28
        }
    ]
}
```

## Sham Language

The Sham language defines the structure of the random data. This language is a superset of JSON that adds integer ranges and generator functions. For the full grammar, refer to `doc/sham.ebnf`. For the base JSON grammar, refer to [RFC 8259](https://tools.ietf.org/html/rfc8259). Sham adds two structures to this grammar:

### Ranges

A range is an inclusive range of integers defined by the production

```ebnf
range : '(' INTEGER ',' INTEGER ')' ;
```

where the first integer is the min and the second is the max. This range includes both the min and max. If a range appears at the beginning of an array, the a random number of elements will be generated in the array. In any other position, a range will evaluate to a random integer in the range.

### Terminal Generators

A terminal generator is a function identifier defined by the production

```ebnf
generator : [a-zA-Z][a-zA-Z]* ;
```

In the generated data, the terminal generator will be replaced by a single value. Generators must match a function defined by the `sham` CLI tool. Unkown generators will return an error.

