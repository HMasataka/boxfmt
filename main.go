package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	overwrite := flag.Bool("w", false, "overwrite the input file")
	output := flag.String("o", "", "output file path")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: mdfmt [options] <file>")
		os.Exit(1)
	}

	if *overwrite && *output != "" {
		fmt.Fprintln(os.Stderr, "error: -w and -o cannot be used together")
		os.Exit(1)
	}

	inputPath := flag.Arg(0)

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	content := string(data)
	result := processFile(content)

	if *overwrite {
		if err := os.WriteFile(inputPath, []byte(result), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	} else if *output != "" {
		if err := os.WriteFile(*output, []byte(result), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Print(result)
	}
}
