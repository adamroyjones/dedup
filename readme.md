# dedup

This program removes duplicate lines from standard input or a file.

- [Current status](#current-status)
- [Why does this program exist?](#why-does-this-program-exist)
- [Installing](#installing)
- [Usage](#usage)

## Current status

It seems fine. There are currently no tests.

## Why does this program exist?

[uniq](https://en.wikipedia.org/wiki/Uniq) takes a collection of lines and
removes its adjacent duplicates. To remove _all_ duplicates with uniq, you must
first sort the input. This is not always desirable.

The (trivial) program in this repository removes the adjacency condition: all
duplicate lines are removed. The order of the input is preserved.

It also support case-insensitive deduplication with the `-i` flag.

## Installing

This requires a Go toolchain.

```sh
go install github.com/adamroyjones/dedup@latest
```

## Usage

```sh
dedup -h
```
