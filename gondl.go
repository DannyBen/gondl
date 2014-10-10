package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DannyBen/filecache"
	"github.com/DannyBen/quandl"
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const version = "0.1.0"

func main() {
	arguments, _ := docopt.Parse(usage, nil, true, version, false)

	switch {
	case arguments["get"] != false:
		getSymbol(arguments)
	case arguments["list"] != false:
		getList(arguments)
	case arguments["search"] != false:
		getSearch(arguments)
	}
}

// getSymbol downloads symbol data from Quandl and
// outputs it to stdout or file
func getSymbol(a map[string]interface{}) {
	quandlSetup(a)

	format := a["--format"].(string)
	showUrl := a["--url"].(bool)
	symbols := a["<symbol>"].([]string)

	opts := getOptions(a, "column", "rows", "trim_start",
		"trim_end", "sort_order")

	var result []byte
	var err error
	if len(symbols) == 1 {
		result, err = quandl.GetSymbolRaw(symbols[0], format, opts)
	} else {
		result, err = quandl.GetSymbolsRaw(symbols, format, opts)
	}
	panicon(err)

	output(a, result, format)
	if showUrl {
		fmt.Println("Quandl URL:", quandl.LastUrl)
	}
}

// getList downloads list of symbols for a given source
// and outputs it to stdout or file
func getList(a map[string]interface{}) {
	quandlSetup(a)
	source := a["<source>"].(string)
	format := a["--format"].(string)
	showUrl := a["--url"].(bool)
	page, _ := strconv.Atoi(a["--page"].(string))
	perPage, _ := strconv.Atoi(a["--per_page"].(string))

	result, err := quandl.GetListRaw(source, format, page, perPage)
	panicon(err)

	output(a, result, format)
	if showUrl {
		fmt.Println("Quandl URL:", quandl.LastUrl)
	}
}

// getSearch downloads search results given query and
// outputs it to stdout or file
func getSearch(a map[string]interface{}) {
	quandlSetup(a)
	query := a["<query>"].(string)
	format := a["--format"].(string)
	showUrl := a["--url"].(bool)
	page, _ := strconv.Atoi(a["--page"].(string))
	perPage, _ := strconv.Atoi(a["--per_page"].(string))

	if format == "csv" {
		format = "json"
	}

	result, err := quandl.GetSearchRaw(query, format, page, perPage)
	panicon(err)

	output(a, result, format)
	if showUrl {
		fmt.Println("Quandl URL:", quandl.LastUrl)
	}
}

// getOptions converts command line flags to quandl query string options
func getOptions(a map[string]interface{}, names ...string) quandl.Options {
	opts := quandl.Options{}
	for _, n := range names {
		key := string("--" + n)
		if a[key] != nil {
			opts.Set(n, a[key].(string))
		}
	}

	return opts
}

// output sends formatted output to stdout
func output(a map[string]interface{}, result []byte, format string) {
	outfile := a["--out"]

	var out bytes.Buffer

	if format == "json" {
		json.Indent(&out, result, "", "\t")
	} else {
		out.Write(result)
	}

	if outfile == nil {
		out.WriteTo(os.Stdout)
	} else {
		err := ioutil.WriteFile(outfile.(string), out.Bytes(), 0644)
		panicon(err)
	}
}

// quandlSetup configures the quandl object before each call
func quandlSetup(a map[string]interface{}) {
	if a["--apikey"] != nil {
		quandl.ApiKey = a["--apikey"].(string)
	}

	cacheDir := a["--cachedir"].(string)
	cacheLife, _ := strconv.ParseFloat(a["--cache"].(string), 32)
	quandl.CacheHandler = filecache.Handler{cacheDir, cacheLife}
}

// showArgs shows the command line args for debugging
func showArgs(a map[string]interface{}) {
	for k, v := range a {
		fmt.Println(k, v)
	}
	os.Exit(0)
}

func panicon(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

const usage = `Gondl - Command line console for Quandl

Usage:
  gondl --help | -h  
  gondl --version | -v  
  gondl get <symbol>... [options]  
  gondl list <source> [options]  
  gondl search <query> [options]  

Options:  
  -h, --help                Show this help  
  -v, --version             Show version details  

Global Options:  
  -k, --apikey <key>        Send this api key with the request  
  -f, --format <format>     Output as csv, json or xml [default: csv]  
  -o, --out <file>          Save to file  
  -u, --url                 Show the request URL  
  -C, --cachedir <dir>      Set cache directory [default: ./cache]  
  -c, --cache <mins>        Set cache life to <mins> minutes, 0 to disable   
                            [default: 240]  

Get Options:  
  -n, --column <n>          Request data column <n> only  
  -r, --rows <n>            Request <n> rows  
  -t, --trim_start <date>   Start data at <date>, format yyyy-mm-dd  
  -T, --trim_end <date>     End data at <date>, format yyyy-mm-dd  
  -s, --sort_order <order>  Set sort order to asc or desc  

Search/List Options:  
  -p, --page <n>            Start at page <n> [default: 1]  
  -P, --per_page <n>        Show <n> results per page [default: 300]  

`

// :)
