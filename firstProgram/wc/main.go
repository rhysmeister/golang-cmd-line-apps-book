package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func count(r io.Reader, countLines bool, countBytes bool) int {
	scanner := bufio.NewScanner(r)
	// This is a bit wonky but it's just for the exercise: supplying both flags messes this up
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}
	if countBytes {
		scanner.Split(bufio.ScanBytes)
	}
	wc := 0
	for scanner.Scan() {
		wc++
	}
	return wc
}

func main() {
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count bytes")
	file := flag.String("file", "", "File to read text from")
	flag.Parse()
	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			fmt.Printf("%s does not exist", *file)
			os.Exit(1)
		}
		defer f.Close()
		fmt.Println(count(f, *lines, *bytes))
	} else {
		fmt.Println(count(os.Stdin, *lines, *bytes))
	}
}
