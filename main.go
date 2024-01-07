package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
)

const version = "0.0.7"

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, `dedup: Remove duplicate lines from standard input or a file.

This program is an analogue of uniq without the adjacency condition on its
lines. See github.com/adamroyjones/dedup.

Usage:
  Deduplicate STDIN and print the results to STDOUT:
    cat file | dedup

  Deduplicate the contents of a file and print the results to STDOUT:
    dedup file

  Deduplicate the contents of a file and overwrite it:
    dedup -w file

  Print version information and exit
    dedup -v`)
	}

	var v, w bool
	flag.BoolVar(&v, "v", false, "Print version information and exit.")
	flag.BoolVar(&w, "w", false, "If provided a filename, modify it in-place.")
	flag.Parse()

	if v {
		fmt.Println("dedup version " + version)
		os.Exit(0)
	}

	if err := dedup(os.Args, w); err != nil {
		fmt.Fprintln(os.Stderr, "dedup: "+err.Error())
		os.Exit(1)
	}
}

func dedup(args []string, write bool) (err error) {
	in, out, err := preparePipes(args, write)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, in.Close(), out.Close()) }()

	if err := dedupLines(out, in); err != nil {
		return fmt.Errorf("writing the deduplicated lines: %w", err)
	}

	if write {
		if err := os.Rename(out.Name(), in.Name()); err != nil {
			return fmt.Errorf("writing the deduplicated file: %w", err)
		}
	}

	return nil
}

func preparePipes(args []string, write bool) (*os.File, *os.File, error) {
	switch len(args) {
	case 0:
		panic("unreachable")

	case 1:
		// Deduplicate STDIN and print the results to STDOUT.
		return os.Stdin, os.Stdout, nil

	case 2:
		// Deduplicate a file and print the results to STDOUT.
		if write {
			return nil, nil, errors.New("given the -w flag, but not given a corresponding file to overwrite")
		}

		in, err := os.Open(os.Args[1])
		if err != nil {
			return nil, nil, fmt.Errorf("unable to open %q for reading: %w", os.Args[1], err)
		}

		return in, os.Stdout, nil

	case 3:
		// Deduplicate a file in-place.
		if !write {
			return nil, nil, errors.New("given more than one file: use cat, if necessary")
		}

		if os.Args[1] != "-w" {
			return nil, nil, errors.New("the -w flag must be given as the first argument")
		}

		in, err := os.Open(os.Args[2])
		if err != nil {
			return nil, nil, fmt.Errorf("unable to open file %q for reading: %w", os.Args[1], err)
		}

		out, err := os.CreateTemp(os.TempDir(), "dedup-")
		if err != nil {
			return nil, nil, fmt.Errorf("unable to create a temporary file for deduplication: %w", err)
		}

		return in, out, nil

	default:
		// Given 4 or more arguments. This is an error case.
		return nil, nil, errors.New("unexpectedly given more than 3 arguments")
	}
}

func dedupLines(out, in *os.File) error {
	dedupedLines := map[string]struct{}{}

	scanner := bufio.NewScanner(in)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if _, ok := dedupedLines[line]; ok {
			continue
		}

		if _, err := out.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("writing a line to the output file: %w", err)
		}
		dedupedLines[line] = struct{}{}
	}

	return nil
}
