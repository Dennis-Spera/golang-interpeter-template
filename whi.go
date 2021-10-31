package main

// tutorial for go
// export GOPATH=/c/Users/User/Desktop/Dropbox/projects/golang:/C/Users/User/go
// 	"tutorial/condif"
//	"tutorial/help"
//---------------------------------------------------------------
// Translated to Go, the saying is as goes:
//
// Accept interfaces, return structs
//
// cd C:\Users\Admin.DIGITALSTORM-PC\Dropbox\projects\golang\whi
// cd /C/Users/Admin.DIGITALSTORM-PC/Dropbox/projects/golang/whi
//----------------------------------------------------------------
//
//  setup directory structure
//
//  mkdir foo
//  create file foo/foo.go
//
//  go mod init main
//
//  "main/foo"
//
//  in main.go or main go program
//
//  "main/lexical"
//
//
//------------------------------------------------------------------

import (
	"flag"
	"fmt"
	"main/lexical"
	"os"
	"regexp"
)

func returnErrorMap() map[string]string {

	return map[string]string{
		"INP-1000": "improper input file extention",
		"INP-1001": "no input file specified",
		"INP-1002": "file name cannot begin with a numeric",
		"INP-1003": "I/O error reading file",
		"LEX-1004": "misplaced literal",
		"GRM-1005": "first word of a program should be PROGRAM",
	}
}

func validateCommandLine(file string, emap map[string]string) {
	if file == "null_file" {
		fmt.Println("INP-1001: ", emap["INP-1001"])
		os.Exit(1001)
	}

	matched1, _ := regexp.MatchString("(.*).(whi)", file)
	if !matched1 {
		fmt.Println("INP-1000: ", emap["INP-1000"])
		os.Exit(1000)
	}

	matched2, _ := regexp.MatchString("^[0-9]+", file)
	if matched2 {
		fmt.Println("INP-1002: ", emap["INP-1002"])
		os.Exit(1002)
	}

	_, err := os.Stat(file)
	if err != nil {
		fmt.Println("INP-1003: ", emap["INP-1003"]+" "+file)
		os.Exit(1003)
	}

	// does file exist

}

func main() {

	infilePtr := flag.String("c", "sample.whi", "an input file to interpet")
	flag.Parse()

	//------------------------
	// get system error map
	//------------------------
	emap := returnErrorMap()
	//------------------------
	// command line validation
	//------------------------
	validateCommandLine(*infilePtr, emap)
	//------------------------
	// syntax checking
	//------------------------
	tokens := lexical.Analysis(*infilePtr, emap)
	status := lexical.GrammarCompiler(tokens, emap)

	if status == "ok" {
		fmt.Println("Grammer analysis complete...")
	}

}
