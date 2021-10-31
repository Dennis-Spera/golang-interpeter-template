package lexical

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// Analysis - syntactical analysis
//   tokens - {[A-Z|a-z]*[A-Z|a-z|0-9|_]*}
//   statement terminator - ;
//   token terminator - ';' ,' '
//   operand -
//
//   statements
//    declare a variable
//       decl  var  nbr; # defaults to 0
//       decl  var  nbr = 12;
//       decl  var  str; # defaults to 'undefined'
//       decl  var  str = 'hello';
//   store to array
//
//   next
//
//     ** come up with grammars for statements for syntactly correct program
//
//   to do:
//     add variable using ptr and line feeds to
//     determine position in file in case of an
//     error
//
//     account for operator, assignments, operands,
//     constants
//
//     what special characters to reject
//---------------------------------------------------------------
//
//  types of tokens
//  ---------------
//  identifier,keyword    ^\w+          state=1
//  delimeter             space, ;      state=2
//  integer               0-9+          state=3
//  operator              ==,=<,=>      state=5
//  assignment            =             state=6
//  condition             (...)         state=7
//
//  to do:
//     code for state 7
//     state diagrams pass1
//     syntax trees pass2
//----------------------------------------------------------------
// finite state machine
//
//  state transitions
//
//  state = 0 -- initial state or null state
//  state 0 -> (all other states)
//  state 1 -> 2  -- identifier to delimeter
//  state 2 -> 2  -- delimeter  to delimeter
//  state 2 -> 1  -- delimeter  to identifier
//  state 3 -> 2  -- integer to delimeter
//  state 2 -> 3  -- delimeter to integer
//  state 1 -> 5  -- identifier to operator
//  state 2 -> 5  -- delimeter to operator
//  state 5 -> 2  -- operator to delimeter
//  state 3 -> 5  -- integer to operator
//  state 1 -> 6  -- identifier to assignment
//  state 6 -> 1  -- assignment to identifier
//  state 2 -> 6  -- delimeter  to assignment
//  state 6 -> 2  -- assignment to delimeter
//
//
//  add var that gets set when a transition is made
//  so if it doesn't get set we know we have a situation
//  with unhandled transition
//
//----------------------------------------------------------------
//
func GrammarCompiler(tokens []string, emap map[string]string) string {
	var finiteStateMachine string
	var state string

	// BNF - BAckus NAur Form
	//
	// program :: literal :: terminator
	// datatype :: identifier ::  equal :: literal :: terminator
	// const :: identifier :: literal :: terminator
	// identifier :: assignment :: identifier || literal

	state = "un-parsed"
	a := tokens
	for i, s := range a {

		if i == 0 {
			if strings.ToLower(s) != "program" {
				fmt.Println("GRM-1005: ", emap["GRM-1005"])
				os.Exit(1005)
			}

			if strings.ToLower(s) == "program" || finiteStateMachine == "program" {
				if state == "un-parsed" {
					finiteStateMachine = "program"
					state = "parsing"

				} else if state == "parsing" {
					// checked if reserved word
					// check if identifier / literal
					state = "literal"

				} else if state == "literal" {
					// checked if terminator
					state = "un-parsed"

				}
			}

			fmt.Println("[statement]: ", finiteStateMachine, " ", "[state]: ", state)
		}
		//fmt.Println(i, s)
	}
	return (finiteStateMachine)
}

// Analysis - syntactical analysis
func Analysis(file string, emap map[string]string) []string {
	var codeString []string
	var token string
	//var lastCh string
	var ptr int
	var state int
	var lineNumber int
	var tokens []string
	//__BLANK__ := " "

	filebuffer, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	inputdata := string(filebuffer)
	data := bufio.NewScanner(strings.NewReader(inputdata))
	data.Split(bufio.ScanRunes)

	for data.Scan() {
		codeString = append(codeString, data.Text())
	}

	sizeOfCodeString := len(codeString)

	// fmt.Println("sizeOfCodeString=", sizeOfCodeString)
	ptr = 0
	state = 0
	lineNumber = 1
	for ptr < sizeOfCodeString {

		if codeString[ptr] == "\n" {
			lineNumber++
		}

		re := regexp.MustCompile("[a-zA-Z]")
		if re.MatchString(codeString[ptr]) {

			if state == 0 {
				// state transitiion 0 -> 1
				// initial state starting with identifier (var, reserved word)
				state = 1
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 1 {
				// state transitiion none
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 2 {
				// state transitiion delimeter to identifier
				state = 1
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			}

		}

		re = regexp.MustCompile("[0-9]")
		if re.MatchString(codeString[ptr]) {

			if state == 1 {
				// state transitiion none identifier to identifier
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 2 {
				// state transitiion delimeter to integer
				state = 3
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 3 {
				// state transitiion none identifier to identifier
				state = 3
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			}

		}

		re = regexp.MustCompile("[=]")
		if re.MatchString(codeString[ptr]) {

			if state == 1 && codeString[ptr+1] == " " {
				// state transitiion identifier to assignment
				state = 6
				token += codeString[ptr]
				//lastCh = codeString[ptr]

			} else if state == 1 && codeString[ptr] == "=" && codeString[ptr+1] == "=" {
				// state transitiion  whitespace to assignment
				tokens = append(tokens, token)
				token = ""
				state = 5
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 1 && codeString[ptr] == "=" && codeString[ptr+1] == "<" {
				// state transitiion  whitespace to assignment
				state = 5
				tokens = append(tokens, token)
				token = ""
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 1 && codeString[ptr] == "=" && codeString[ptr+1] == ">" {
				// state transitiion  whitespace to assignment
				state = 5
				tokens = append(tokens, token)
				token = ""
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 2 && codeString[ptr] == "=" && codeString[ptr+1] == "=" {
				// state transitiion  whitespace to assignment
				state = 5
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 2 && codeString[ptr+1] == " " {
				// state transitiion  whitespace to assignment
				state = 6
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 2 && codeString[ptr] == "=" && codeString[ptr+1] == ">" {
				// state transitiion  whitespace to assignment
				state = 5
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 2 && codeString[ptr] == "=" && codeString[ptr+1] == "<" {
				// state transitiion  whitespace to assignment
				state = 5
				token += codeString[ptr]
				//lastCh = codeString[ptr]
			} else if state == 5 && codeString[ptr] == "=" {
				// state transitiion  whitespace to assignment
				state = 5
				token += codeString[ptr]
				tokens = append(tokens, token)
				token = ""
				//lastCh = codeString[ptr]
			}

		}

		re = regexp.MustCompile("[>]")
		if re.MatchString(codeString[ptr]) {
			if state == 5 && codeString[ptr] == ">" {
				// state transitiion  whitespace to assignment
				state = 0
				token += codeString[ptr]
				tokens = append(tokens, token)
				token = ""
				//lastCh = codeString[ptr]
			}

		}

		re = regexp.MustCompile("[<]")
		if re.MatchString(codeString[ptr]) {
			if state == 5 && codeString[ptr] == "<" {
				// state transitiion  whitespace to assignment
				state = 0
				token += codeString[ptr]
				tokens = append(tokens, token)
				token = ""
				//lastCh = codeString[ptr]
			}

		}

		re = regexp.MustCompile("[ ;\n]")
		if re.MatchString(codeString[ptr]) {

			if state == 1 {
				// state transitiion identifier to delimeter
				state = 2
				tokens = append(tokens, token)
				if codeString[ptr] == ";" {
					token = codeString[ptr]
					tokens = append(tokens, token)
				}
				token = ""

				//lastCh = codeString[ptr]
			}

			if state == 3 {
				// state transitiion integer to delimeter
				state = 2
				tokens = append(tokens, token)
				token = ""
				//lastCh = codeString[ptr]
			}

			if state == 6 {
				// state transitiion assignment to delimeter
				state = 2
				tokens = append(tokens, token)
				token = ""
				//lastCh = codeString[ptr]
			}

			if state == 5 {
				// state transitiion assignment to delimeter
				state = 2
				tokens = append(tokens, token)
				token = ""
				//lastCh = codeString[ptr]
			}

		}

		ptr++
	}

	// fmt.Println("token : " + token + "\n")
	tokens = append(tokens, token)

	// a := tokens
	// for i, s := range a {
	//	fmt.Println(i, s)
	//}
	return (tokens)
}

