package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var (
		inFileName  = flag.String("input", "", "Location of .apib file as input.")
		outFileName = flag.String("output", "", "Location of Postman file.")

		inFileByte []byte
		outFile    *os.File
		err        error
	)

	flag.Parse()

	if *inFileName == "" {
		inFileByte, err = ioutil.ReadAll(os.Stdin)
	} else {
		inFileByte, err = ioutil.ReadFile(*inFileName)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error on read input:", err)
		os.Exit(1)
	}

	if *outFileName == "" {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create(*outFileName)
		if err != nil {
			os.Exit(1)
		}
		defer outFile.Close()
	}

	_, err = outFile.Write(inFileByte)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error on write output:", err)
		os.Exit(1)
	}
}
