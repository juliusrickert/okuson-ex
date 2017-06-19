package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// TODO: make stringArrayFlags support filenames containing a comma
type stringArrayFlags []string

func (f *stringArrayFlags) String() string {
	return strings.Join(*f, ",")
}

func (f *stringArrayFlags) Set(value string) error {
	*f = strings.Split(value, ",")
	return nil
}

var (
	okusonURL    string
	templateFile string
	action       string
	outputFormat string
	inputFiles   stringArrayFlags
)

func usage() {
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.StringVar(&action, "a", "", "Action: 'get' to get tasks or 'combine' to combine tasks")
	flag.StringVar(&outputFormat, "o", "", "Output format: 'json' or 'tex'")
	flag.StringVar(&okusonURL, "url", "", "When action is 'get': URL to OKUSON start page, without index.html and with trailing slash")
	flag.StringVar(&templateFile, "tpl", "", "When output format is 'tex': File containing the template")
	flag.Var(&inputFiles, "i", "When action is 'combine': Comma seperated input files, e.g. \"file1.json,file2.json,file3.json\"")

	flag.Parse()

	if outputFormat != "json" && outputFormat != "tex" {
		usage()
	}

	if outputFormat == "tex" && templateFile == "" {
		usage()
	}

	var err error
	switch action {
	case "get":
		if okusonURL == "" {
			usage()
		}
		err = get()
	case "combine":
		if len(inputFiles) == 0 {
			usage()
		}
		err = combine()
	default:
		usage()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured: %v\n", err)
	}
}
