# dedup

The Unix utility `uniq` removes adjacent duplicates from a `\n`-separated
collection of input strings. This repository contains a program that removes
_all_ duplicates from such a collection. The program preserves the order of the
input.

The implementation is naive and simplistic but, for my current purposes, it is
good enough.

## Installing

```sh
go install gitlab.com/adamroyjones/dedup@latest
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
