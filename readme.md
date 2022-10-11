# dedup

Remove duplicate lines from standard input or a file.

- [Why does this program exist?](#why-does-this-program-exist)
- [Installing](#installing)
- [Usage](#usage)

## Why does this program exist?

The Unix utility `uniq` takes a collection of lines and removes its adjacent
duplicates. To remove _all_ duplicates with `uniq`, you must first sort the
input. This is not always desirable.

The program in this repository, `dedup`, removes the adjacency condition: all
duplicate lines are removed. The program preserves the order of the input.

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
