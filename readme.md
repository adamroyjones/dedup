# dedup

The Unix utility `uniq` takes a collection of lines and removes its adjacent
duplicates. This repository contains a program, `dedup`, that acts like `uniq`,
only it removes the adjacency condition; that is, it removes _all_ duplicates
from such a collection. The program preserves the order of the input.

The implementation is naive and simplistic but, for my current purposes, it is
good enough.

## Installing

```sh
go install github.com/adamroyjones/dedup@latest
```

## Usage

```sh
# Deduplicate STDIN and print the result to STDOUT.
cat file | dedup

# Deduplicate the contents of `file` and print the result to STDOUT.
dedup file

# Deduplicate the contents of `file` and write the results to `file`.
dedup -w file
```
