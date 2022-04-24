package dedup

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	dedupedLines := map[string]bool{}
	output := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	var line string

	for scanner.Scan() {
		line = scanner.Text()
		if _, ok := dedupedLines[line]; !ok {
			output = append(output, line)
			dedupedLines[line] = true
		}
	}

	for _, line := range output {
		fmt.Println(line)
	}
}
