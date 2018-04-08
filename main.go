package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	count = flag.Int("n", 0, "Display the first count lines (default 10)")
	bytes = flag.Int("c", 0, "Display the first bytes")
)

func init() {
	flag.Usage = func() {
		fmt.Printf(`Usage:
  %s [-n count | -b bytes] [FILE ...]

`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *count > 0 && *bytes > 0 {
		fmt.Printf("%s: can't combine line and byte counts\n", os.Args[0])
		os.Exit(1)
	}
	if *count == 0 && *bytes == 0 {
		*count = 10
	}
}

func main() {
	files := flag.Args()

	if len(files) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for i, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			fmt.Printf("%s: %s: No such file or directory\n", os.Args[0], f)
			continue
		}
		if len(files) > 1 {
			displayFileName(f, i)
		}
		file, err := os.Open(f)
		defer file.Close()
		handleError(err)

		rd := bufio.NewReader(file)

		if *count > 0 {
			readCount(rd)
		} else if *bytes > 0 {
			readBytes(rd)
		}
	}
}

func displayFileName(f string, i int) {
	if i > 0 {
		// Display break line after the second or more files
		fmt.Println("")
	}
	fmt.Printf("==> %s <==\n", f)
}

func readCount(rd *bufio.Reader) {
	i := 1
	for {
		line, _, err := rd.ReadLine()
		fmt.Println(string(line))
		if err == io.EOF {
			break
		} else if i >= *count {
			break
		} else if err != nil {
			handleError(err)
		}
		i++
	}
}

func readBytes(rd *bufio.Reader) {
	b := make([]byte, *bytes) // TODO: large byte count is specified
	_, err := rd.Read(b)
	handleError(err)
	fmt.Println(string(b))
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
