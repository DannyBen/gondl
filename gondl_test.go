package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Example_GetSimple() {
	run([]string{"get", "WIKI/CSCO", "-n4", "-t2014-01-01", "-T2014-01-04"})
	// Output:
	// Date,Close
	// 2014-01-03,21.98
	// 2014-01-02,22.0
}

func Example_GetFile() {
	run([]string{"get", "WIKI/AAPL", "-n4", "-t2014-01-01", "-T2014-01-04", "-otmp", "-sasc"})
	data, _ := ioutil.ReadFile("tmp")
	os.Remove("tmp")
	fmt.Println(string(data))
	// Output:
	// Date,Close
	// 2014-01-02,553.13
	// 2014-01-03,540.98
}
