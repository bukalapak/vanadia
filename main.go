package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bukalapak/vanadia/blueprint"
	"github.com/bukalapak/vanadia/config"
	"github.com/bukalapak/vanadia/postman"
)

const (
	defaultConfigFileName = "vanadia.yml"
)

func main() {
	var (
		inFileName     = flag.String("input", "", "Location of .apib file as input.")
		outFileName    = flag.String("output", "", "Location of Postman file.")
		configFileName = flag.String("config", defaultConfigFileName, "Location of vanadia.yml.")

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

	cfg, err := config.FromFile(*configFileName)
	if err != nil {
		func() {
			switch err.(type) {
			case *os.PathError:
				if *configFileName == defaultConfigFileName {
					// If no default config file is defined, then we just
					// use default Config value
					cfg = config.DefaultConfig
					return
				} else {
					fmt.Fprintln(os.Stderr, "error reading config:", err)
					os.Exit(1)
				}
			}
			fmt.Fprintln(os.Stderr, "error reading config:", err)
			os.Exit(1)
		}()
	}

	bp, err := blueprint.GetStructure(inFileByte)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing blueprint:", err)
		os.Exit(1)
	}

	collection, err := postman.CreateCollection(bp)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating collection:", err)
		os.Exit(1)
	}

	if cfg.SchemeToEnv != "" {
		postman.SchemeToEnv(&collection, cfg.SchemeToEnv)
	}
	if cfg.HostToEnv.Segments > 0 {
		postman.HostToEnv(&collection, cfg.HostToEnv.Segments)
	}
	if cfg.AuthTokenToEnv {
		postman.AuthTokenToEnv(&collection)
	}
	postman.AddGlobalHeaders(&collection, cfg.GlobalHeaders)

	json, err := json.MarshalIndent(collection, "", "\t")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error exporting to JSON:", err)
		os.Exit(1)
	}

	_, err = outFile.Write(json)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error on write output:", err)
		os.Exit(1)
	}
}
