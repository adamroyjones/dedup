package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const version = "0.0.8"

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

  Deduplicate the contents of a file, ignoring the casing, and print the results
  to STDOUT:
    dedup -i file

  Deduplicate the contents of a file and overwrite it:
    dedup -w file

  Print version information and exit:
    dedup -v`)
	}

	var i, v, w bool
	flag.BoolVar(&i, "i", false, "Use a case-insensitive match.")
	flag.BoolVar(&v, "v", false, "Print version information and exit.")
	flag.BoolVar(&w, "w", false, "If provided a filename, modify it in-place.")
	flag.Parse()

	if v {
		fmt.Println("dedup version " + version)
		os.Exit(0)
	}

	if err := dedup(flag.Args(), i, w); err != nil {
		fmt.Fprintln(os.Stderr, "dedup: "+err.Error())
		os.Exit(1)
	}
}

func dedup(args []string, ignoreCase, write bool) (err error) {
	in, out, err := preparePipes(args, write)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, in.Close(), out.Close()) }()

	if err := dedupLines(out, in, ignoreCase); err != nil {
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
		// Deduplicate STDIN and print the results to STDOUT.
		if write {
			return nil, nil, errors.New("unable to handle the write (-w) flag without a file")
		}

		return os.Stdin, os.Stdout, nil

	case 1:
		// Deduplicate a file.
		file := args[0]
		in, err := os.Open(file)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to open %q for reading: %w", file, err)
		}

		if !write {
			return in, os.Stdout, nil
		}

		out, err := os.CreateTemp(os.TempDir(), "dedup-")
		if err != nil {
			return nil, nil, fmt.Errorf("unable to create a temporary file for deduplication: %w", err)
		}
		return in, out, nil

	default:
		return nil, nil, errors.New("unexpectedly given more than 1 argument: use cat, if necessary")
	}
}

func dedupLines(out, in *os.File, ignoreCase bool) error {
	dedupedLines := map[string]struct{}{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		lineToMatch := line
		if ignoreCase {
			lineToMatch = strings.ToLower(lineToMatch)
		}
		if _, ok := dedupedLines[lineToMatch]; ok {
			continue
		}

		dedupedLines[lineToMatch] = struct{}{}
		if _, err := out.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("writing a line to the output file: %w", err)
		}
	}

	return nil
}
