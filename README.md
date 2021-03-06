Gondl - Command line console for Quandl
=======================================

<!-- [![Build Status](https://travis-ci.org/DannyBen/gondl.svg?branch=master)](https://travis-ci.org/DannyBen/gondl) -->

Gondl provides command line access to the 
[Quandl API](https://www.quandl.com/help/api).

It was developed in Go.

## Features

* Get data for a symbol
* Get a list of symbols in a data source
* Search the entire Quandl database
* Shows/saves JSON, CSV or XML
* Built in local file cache


## Download Windows Binary

[Download the latest build of gondl.exe](https://github.com/DannyBen/gondl/releases)


## Build from Source (All Platforms)

To build from source on Windows, Linux or Mac - 
[Install Go](https://golang.org/doc/install), then:

	$ go get github.com/DannyBen/gondl
	$ cd $GOPATH/src/github.com/DannyBen/gondl
	$ go build

## Examples

Get data for Apple stock:

	gondl get WIKI/AAPL

Get 3 rows of data as JSON, and use an API Key:

	gondl get WIKI/AAPL -r3 -fjson -kYOUR_KEY

Save data as XML to a file:

	gondl get WIKI/CSCO -fxml -oOutFile.txt --rows 10

Get a list of symbols in a source:

	gondl list WIKI --page 1 --per_page 10

Get search results:

	gondl search "crude oil" --page 1 --per_page 10


## Usage:

    gondl --help | -h  
    gondl --version | -v  
    gondl --config [options]  
    gondl --make-config  
    gondl get <symbol> [options]  
    gondl list <source> [options]  
    gondl search <query> [options]  

## Standalone Options:  

    -h, --help                Show this help.  
    -v, --version             Show version details.  
        --config              Show config files location and info.  
        --make-config         Create a default gondl.json file.  

## Global Options:  

    -k, --apikey <key>        Send this api key with the request  
    -f, --format <format>     Output as csv, json or xml (default: csv)  
    -o, --out <file>          Save to file  
    -u, --url                 Show the request URL  
    -d, --debug               Show all registered arguments  
    -D, --cachedir <dir>      Set cache directory (default: ./cache)  
    -C, --cache <mins>        Set cache life to <mins> minutes  
                              0 to disable (default: 240)  

## Get Options:  

    -c, --column <n>          Request data column <n> only  
    -r, --rows <n>            Request <n> rows  
    -t, --trim_start <date>   Start data at <date>, format yyyy-mm-dd  
    -T, --trim_end <date>     End data at <date>, format yyyy-mm-dd  
    -s, --sort_order <order>  Set sort order to asc or desc  
    -x, --exclude_headers     Exclude CSV headers  
        --exclude_data        Get meta data only (JSON/XML format)  
        --collapse <f>        Set frequency to one of: none | daily |  
                              weekly | monthly | quarterly | annual   
        --transformation <t>  Enable data calculation. Set to one of:  
                              diff | rdiff | cumul | normalize  

## Search/List Options:  

    -p, --page <n>            Start at page <n> (default: 1)  
    -P, --per_page <n>        Show <n> results per page (default: 300)  
