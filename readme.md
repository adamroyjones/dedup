# dedup

This program removes duplicate lines from standard input or a file.

- [Why does this program exist?](#why-does-this-program-exist)
- [Installing](#installing)
- [Usage](#usage)

## Why does this program exist?

[uniq](https://en.wikipedia.org/wiki/Uniq) takes a collection of lines and
removes its adjacent duplicates. To remove _all_ duplicates with uniq, you must
first sort the input. This is not always desirable.

The (trivial) program in this repository removes the adjacency condition: all
duplicate lines are removed. The order of the input is preserved.

## Installing

This requires a Go toolchain.

```sh
go install github.com/adamroyjones/dedup@latest
```

## Usage

```console
$ dedup -h
dedup: Deduplicate identical lines from the input.

This program is uniq without the adjacency condition on its lines. See
github.com/adamroyjones/dedup for more.

Usage:
    Deduplicate STDIN and print the results to STDOUT:
        cat file | dedup

    Deduplicate the contents of a file and print the results to STDOUT:
        dedup [file]

    Deduplicate the contents of a file and overwrite it:
        dedup -w [file]
```
