Gondl - Command line console for Quandl
=======================================

Gondl provides command line access to the 
[Quandl API](https://www.quandl.com/help/api).

It was developed in Go.


## Download Windows Binary

TODO


## Build from Source

TODO


## Exampmles

Get data for Apple stock:

	gondl get WIKI/AAPL

Get 3 rows of data as JSON, and use an API Key:

	gondl get WIKI/AAPL -r3 -fjson -kYOUR_KEY

Save multiple symbols as XML to a file:

	gondl get WIKI/AAPL WIKI/CSCO -fxml -oOutFile.txt --rows 10

Get a list of symbols in a source:

	gondl list WIKI --page 1 --per_page 10

Get search results:

	gondl search "crude oil" --page 1 --per_page 10


## Usage:
    gondl --help | -h  
    gondl --version | -v  
    gondl get <symbol>... [options]  
    gondl list <source> [options]  
    gondl search <query> [options]  

## Options:  
    -h, --help                Show this help  
    -v, --version             Show version details  

## Global Options:  
    -k, --apikey <key>        Send this api key with the request  
    -f, --format <format>     Output as csv, json or xml [default: csv]  
    -o, --out <file>          Save to file  
    -u, --url                 Show the request URL  
    -C, --cachedir <dir>      Set cache directory [default: ./cache]  
    -c, --cache <mins>        Set cache life to <min> minutes, 0 to disable   
                              [default: 240]  

## Get Options:  
    -n, --column <n>          Request data column <n> only  
    -r, --rows <n>            Request <n> rows  
    -t, --trim_start <date>   Start data at <date>, format yyyy-mm-dd  
    -T, --trim_end <date>     End data at <date>, format yyyy-mm-dd  
    -s, --sort_order <order>  Set sort order to asc or desc  

## Search/List Options:  
    -p, --page <n>            Start at page <n> [default: 1]  
    -P, --per_page <n>        Show <n> results per page [default: 300]  
