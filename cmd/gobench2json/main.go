// This program reads Go benchmark output from standard input, parses it,
// and prints a JSON representation of the benchmarks grouped by package to standard output.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Benchmark holds the data for a single benchmark result.
type Benchmark struct {
	Name        string  `json:"name"`
	Runs        int64   `json:"runs"`
	NsPerOp     float64 `json:"ns_per_op,omitempty"`
	BPerOp      float64 `json:"b_per_op,omitempty"`
	AllocsPerOp int64   `json:"allocs_per_op,omitempty"`
	MBPerS      float64 `json:"mb_per_s,omitempty"`
}

func main() {
	// benchmarksByPackage will store the final data, grouped by package name.
	benchmarksByPackage := make(map[string][]Benchmark)
	// currentPackage holds the package name for the benchmarks currently being processed.
	var currentPackage string

    // The regular expression is designed to parse a line of benchmark output.
    // It is composed of the following parts:
    // 1. `^`: Asserts that the match must start at the beginning of the line.
    // 2. `(Benchmark[^ ]+)`: This is the first capturing group. It captures the full benchmark name.
    //    - `Benchmark`: Matches the literal string "Benchmark".
    //    - `[^ ]+`: Matches one or more characters that are not a tab. This is key to capturing
    //      the entire benchmark name, including any special characters, without including the
    //      trailing spaces or the tab separator.
    // 3. ` +`: Matches one or more whitespace characters that separate the name from the iteration count.
    // 4. `([0-9]+)`: This is the second capturing group. It captures the number of iterations.
    //
    // Example line from `go test -bench` output:
    // BenchmarkAuthenticator_HappyPath-8                         101954         13603 ns/op
    //
    // When applied to the example line:
    // - The first capturing group will match: "BenchmarkAuthenticator_HappyPath-8                      "
    // - The second capturing group will match: "101954"
    // The code then trims the whitespace from the first group to get the clean name.
	re := regexp.MustCompile(`^(Benchmark[^\t]+)\s+([0-9]+)`)
	scanner := bufio.NewScanner(os.Stdin)

	// Process the input line by line.
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line defines the package.
		if strings.HasPrefix(line, "pkg: ") {
			currentPackage = strings.TrimSpace(strings.TrimPrefix(line, "pkg: "))
			continue
		}

		// Use the regex to check if the line is a benchmark result.
		matches := re.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		runs, err := strconv.ParseInt(matches[2], 10, 64)
		if err != nil {
			continue
		}

		b := Benchmark{
			Name: strings.TrimSpace(matches[1]),
			Runs: runs,
		}

		// Parse the metrics for the benchmark.
		fields := strings.Fields(line)
		for i := 2; i < len(fields)-1; i += 2 {
			value, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				continue
			}
			unit := fields[i+1]
			switch unit {
			case "ns/op":
				b.NsPerOp = value
			case "B/op":
				b.BPerOp = value
			case "allocs/op":
				b.AllocsPerOp = int64(value)
			case "MB/s":
				b.MBPerS = value
			}
		}

		// Add the benchmark to the current package's slice.
		if currentPackage != "" {
			benchmarksByPackage[currentPackage] = append(benchmarksByPackage[currentPackage], b)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		os.Exit(1)
	}

	// If no benchmarks were found, output an empty JSON object.
	if len(benchmarksByPackage) == 0 {
		fmt.Println("{}")
		return
	}

	// Marshal the map into a JSON string.
	jsonData, err := json.MarshalIndent(benchmarksByPackage, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling json: %v\n", err)
		os.Exit(1)
	}

	// Print the final JSON output.
	fmt.Println(string(jsonData))
}
