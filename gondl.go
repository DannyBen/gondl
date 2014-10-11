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
	"sort"
	"strconv"
)

const version = "0.1.2"

func main() {
	run(nil)
}

// run is injectable main
func run(args []string) {
	arguments, _ := docopt.Parse(usage, args, true, version, false)
	config := loadConfig("config.json")
	arguments = merge(arguments, config)

	switch {
	case arguments["--config"] != false:
		makeConfig()
	case arguments["get"] != false:
		getSymbol(arguments)
	case arguments["list"] != false:
		getList(arguments)
	case arguments["search"] != false:
		getSearch(arguments)
	}

	if arguments["--debug"].(bool) {
		showArgs(arguments)
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
		"trim_end", "sort_order", "collapse", "transformation",
		"exclude_headers")

	var result []byte
	var err error
	if len(symbols) == 1 {
		result, err = quandl.GetSymbolRaw(symbols[0], format, opts)
	} else {
		result, err = quandl.GetSymbolsRaw(symbols, format, opts)
	}
	panicon(err)

	output(a, result, format)
	showQuandlUrl(showUrl, quandl.LastUrl)
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
	showQuandlUrl(showUrl, quandl.LastUrl)
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

	// TODO: Remove this patch when Quandl guys fix the bug
	//       Also remove from quandl library
	if format == "csv" {
		format = "json"
	}

	result, err := quandl.GetSearchRaw(query, format, page, perPage)
	panicon(err)

	output(a, result, format)
	showQuandlUrl(showUrl, quandl.LastUrl)
}

// makeConfig creates a default config.json template file
func makeConfig() {
	ioutil.WriteFile("config.json", []byte(configTemplate), 0644)
}

// getOptions converts command line flags to quandl query string options
func getOptions(a map[string]interface{}, names ...string) quandl.Options {
	opts := quandl.Options{}
	for _, n := range names {
		key := string("--" + n)
		if a[key] != nil {
			if v, ok := a[key].(string); ok {
				opts.Set(n, v)
			} else if v, ok := a[key].(bool); ok {
				if v {
					opts.Set(n, "true")
				}
			}
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
	if cacheLife > 0 {
		quandl.CacheHandler = filecache.Handler{cacheDir, cacheLife}
	}
}

// showArgs shows the command line args (--debug)
func showArgs(a map[string]interface{}) {
	fmt.Println("\nRegistered Arguments:")
	var keys []string
	for k := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("  %-18s %v\n", k, a[k])
	}
}

// showQuandlUrl shows the url used in the request (--url)
func showQuandlUrl(show bool, url string) {
	if show {
		fmt.Printf("\nQuandl URL:\n%s\n\n", url)
	}
}

// loadConfig loads a JSON config file if available
func loadConfig(filename string) map[string]interface{} {
	var result map[string]interface{}
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return result
	}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		log.Fatal("Error in config.json: " + err.Error())
	}
	return result
}

// merge combines two maps.
// truthiness takes priority over falsiness
// mapA takes priority over mapB
func merge(mapA, mapB map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range mapA {
		result[k] = v
	}
	for k, v := range mapB {
		if _, ok := result[k]; !ok || result[k] == nil || result[k] == false {
			result[k] = v
		}
	}
	return result
}

func panicon(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

const configTemplate = `{
	"--apikey": "YOUR_KEY",
	"--trim_start": "2014-01-01",
	"--per_page": "10",
	"--url": true
}
`

const usage = `Gondl - Command line console for Quandl

Usage:
  gondl --help | -h  
  gondl --version | -v  
  gondl --config  
  gondl get <symbol>... [options]  
  gondl list <source> [options]  
  gondl search <query> [options]  

Standalone Options:  
  -h, --help                Show this help  
  -v, --version             Show version details  
      --config              Create a default config.json file  
                            You may place any of the --options in it  

Global Options:  
  -k, --apikey <key>        Send this api key with the request  
  -f, --format <format>     Output as csv, json or xml [default: csv]  
  -o, --out <file>          Save to file  
  -u, --url                 Show the request URL  
  -d, --debug               Show all registered arguments  
  -C, --cachedir <dir>      Set cache directory [default: ./cache]  
  -c, --cache <mins>        Set cache life to <mins> minutes  
                            0 to disable [default: 240]  

Get Options:  
  -n, --column <n>          Request data column <n> only  
  -r, --rows <n>            Request <n> rows  
  -t, --trim_start <date>   Start data at <date>, format yyyy-mm-dd  
  -T, --trim_end <date>     End data at <date>, format yyyy-mm-dd  
  -s, --sort_order <order>  Set sort order to asc or desc  
  -x, --exclude_headers     Exclude CSV headers  
      --collapse <f>        Set frequency to one of: none | daily |  
                            weekly | monthly | quarterly | annual   
      --transformation <t>  Enable data calculation. Set to one of:  
                            diff | rdiff | cumul | normalize  

Search/List Options:  
  -p, --page <n>            Start at page <n> [default: 1]  
  -P, --per_page <n>        Show <n> results per page [default: 300]  

`

//;D
