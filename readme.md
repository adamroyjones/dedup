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
```
