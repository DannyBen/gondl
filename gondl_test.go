package main

import (
	"encoding/json"
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

func Example_GetMulti() {
	run([]string{"get", "WIKI/AAPL.4", "WIKI/MSFT.4", "-t2014-01-01", "-T2014-01-04", "-x"})
	// Output:
	// 2014-01-03,540.98,36.91
	// 2014-01-02,553.13,37.16
}

func Example_GetCollapse() {
	run([]string{"get", "WIKI/AAPL", "-t2014-01-01", "-T2014-02-01", "-x", "-n4", "--collapse=weekly"})
	// Output:
	// 2014-02-02,500.6
	// 2014-01-26,546.07
	// 2014-01-19,540.67
	// 2014-01-12,532.94
	// 2014-01-05,540.98
}

func Example_GetTransformation() {
	run([]string{"get", "WIKI/MSFT", "-t2014-01-01", "-T2014-01-10", "-sasc", "-n4", "--transformation=normalize"})
	// Output:
	// Date,Close
	// 2014-01-02,100.0
	// 2014-01-03,99.327233584499
	// 2014-01-06,97.228202368138
	// 2014-01-07,97.981700753498
	// 2014-01-08,96.232508073197
	// 2014-01-09,95.613562970936
	// 2014-01-10,96.986006458558
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

func Example_List() {
	run([]string{"list", "WIKI", "-P2", "-fjson", "-otmp"})
	data, _ := ioutil.ReadFile("tmp")
	os.Remove("tmp")
	var r map[string]interface{}
	json.Unmarshal(data, &r)

	countTest := r["total_count"].(float64) > 2000
	lenTest := len(r["docs"].([]interface{}))
	sourceTest := r["sources"].([]interface{})[0].(map[string]interface{})["id"]

	fmt.Println(countTest, lenTest, sourceTest)
	// Output:
	// true 2 4922
}

func Example_Search() {
	run([]string{"search", "nasdaq composite", "-P3", "-fjson", "-otmp"})
	data, _ := ioutil.ReadFile("tmp")
	os.Remove("tmp")
	var r map[string]interface{}
	json.Unmarshal(data, &r)

	countTest := r["total_count"].(float64) > 570000
	lenTest := len(r["docs"].([]interface{}))

	fmt.Println(countTest, lenTest)
	// Output:
	// true 3
}
