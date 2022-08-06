package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Usage = func() {
		template := `
dedup: Deduplicate identical lines from the provided input.

Usage:
    Deduplicate STDIN and print the results to STDOUT:
        cat file | dedup

    Deduplicate the contents of %s and print the results to STDOUT:
        dedup file

    Deduplicate the contents of %s and overwrite it:
        dedup -w file
`
		msg := fmt.Sprintf(template, "`file`", "`file`")

		fmt.Fprintf(os.Stderr, strings.TrimPrefix(msg, "\n"))
	}

	writePtr := flag.Bool("w", false, "If provided a filename, modify it in-place.")
	flag.Parse()
	write := *writePtr

	err := dedup(os.Args, write)
	if err != nil {
		fmt.Printf("dedup: %s.\n", err.Error())
		os.Exit(1)
	}
}

func dedup(args []string, write bool) error {
	in, out, err := preparePipes(args, write)
	if err != nil {
		return err
	}
	defer in.Close()
	defer out.Close()

	dedupedLines := dedupLines(in)
	for _, line := range dedupedLines {
		if _, err := out.WriteString(line + "\n"); err != nil {
			if closeErr := out.Close(); closeErr != nil {
				return fmt.Errorf("Error when writing the deduplicated lines: %v; error when closing the destination: %v", err, closeErr)
			}
			return fmt.Errorf("Error when writing the deduplicated lines: %v", err)
		}
	}

	if write {
		if err := os.Rename(out.Name(), in.Name()); err != nil {
			return fmt.Errorf("Error when writing the deduplicated file: %v", err)
		}
	}

	return nil
}

func preparePipes(args []string, write bool) (*os.File, *os.File, error) {
	switch len(args) {
	case 1: // Deduplicate STDIN and print the results to STDOUT.
		return os.Stdin, os.Stdout, nil
	case 2: // Deduplicate a file and print the results to STDOUT.
		if write {
			return nil, nil, fmt.Errorf("Given the '-w' flag, but not given a corresponding file to overwrite")
		}

		in, err := os.Open(os.Args[1])
		if err != nil {
			return nil, nil, fmt.Errorf("Unable to open file '%s' for reading: %v", os.Args[1], err)
		}

		return in, os.Stdout, nil
	case 3: // Deduplicate a file in-place.
		if !write {
			// TODO: This could be generalised, I suppose, but it's not fruitful given that you can already cat the files.
			return nil, nil, fmt.Errorf("Given more than one file")
		}

		if os.Args[1] != "-w" {
			return nil, nil, fmt.Errorf("The '-w' flag should be given as the first argument")
		}

		in, err := os.Open(os.Args[2])
		if err != nil {
			return nil, nil, fmt.Errorf("Unable to open file '%s' for reading: %v", os.Args[1], err)
		}

		out, err := os.CreateTemp(os.TempDir(), "dedup-")
		if err != nil {
			return nil, nil, fmt.Errorf("Unable to create a temporary file for deduplication: %v", err)
		}

		return in, out, nil
	default: // Given 4 or more arguments, which is an error case.
		return nil, nil, fmt.Errorf("Unexpectedly given more than 3 arguments")
	}
}

func dedupLines(in *os.File) []string {
	var line string
	dedupedLines := map[string]bool{}
	var out []string

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line = scanner.Text()
		if _, ok := dedupedLines[line]; !ok {
			out = append(out, line)
			dedupedLines[line] = true
		}
	}

	return out
}
