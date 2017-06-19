package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

func getLoginData() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "Enter username: ")
	username, _ := reader.ReadString('\n')
	fmt.Fprintf(os.Stderr, "Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		bytePassword = []byte{}
	}
	fmt.Fprintf(os.Stderr, "\n")
	return username, string(bytePassword)
}

func get() error {
	matnum, password := getLoginData()

	sheetNames, err := getSheetList()
	if err != nil {
		return err
	}

	sheets := make([]sheet, len(sheetNames))

	for i, s := range sheetNames {
		fmt.Fprintf(os.Stderr, "Getting sheet %s\n", s)
		sh, err := getSheet(s, matnum, password)
		if err != nil {
			return err
		}
		sheets[i] = *sh
	}

	var makeFunc outputFunction
	if outputFormat == "tex" {
		makeFunc = outputFunction(makeLaTeX)
	} else if outputFormat == "json" {
		makeFunc = outputFunction(makeJSON)
	}

	err = makeFunc(sheets, os.Stdout)

	if err != nil {
		return err
	}

	return nil
}
