package main

import (
	"fmt"
	"os"
	"testapi"
)

var cases = map[string]func(){
	"node":         testapi.Node,
	"organization": testapi.Organization,
	"user":         testapi.User,
}

// Invoke the requested testapi test function
func main() {
	testcase := os.Args[1]
	fn, ok := cases[testcase]
	if ok {
		fn()
	} else {
		fmt.Fprintf(os.Stderr, "Requested case %+s was not found\n", testcase)
	}
}
