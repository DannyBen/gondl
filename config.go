package main

import (
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const configName = "gondl.json"

// Type config holds a key=value configuration map
type Config struct {
	Values map[string]interface{}
}

// NewConfig returns a ready to work with configuration object
func NewConfig(usage string, args []string, version string) Config {
	config := Config{}
	config.init(usage, args, version)
	config.ammend(workingDirConfig())
	config.ammend(homeDirConfig())
	config.ammend(localDirConfig())
	config.fill("--cachedir", "./cache")
	config.fill("--cache", "240")
	config.fill("--page", "1")
	config.fill("--per_page", "300")
	config.fill("--format", "csv")
	return config
}

// init loads values from the docopt into the config object
func (c *Config) init(usage string, args []string, version string) {
	c.Values, _ = docopt.Parse(usage, args, true, version, false)
}

// fill populates config fields with values, if they are nil
func (c *Config) fill(key, value string) {
	if c.Values[key] == nil {
		c.Values[key] = value
	}
}

// ammend reads a config file and adds its vales to the config
// object. Formerly existing values will prevail.
func (c *Config) ammend(path string) {
	loaded := loadConfig(path)
	c.Values = merge(c.Values, loaded)
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
		log.Fatalf("Error in %v:\n%v", filename, err.Error())
	}

	return result
}

// merge combines two Configs.
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

// localDirConfig returns the path to the local config file
func localDirConfig() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	panicon(err)
	return dir + string(filepath.Separator) + configName
}

// workingDirConfig returns the path working directory config file
func workingDirConfig() string {
	dir, err := os.Getwd()
	panicon(err)
	return dir + string(filepath.Separator) + configName
}

// homeDirConfig returns the path to the user home directory config
// file
func homeDirConfig() string {
	usr, err := user.Current()
	panicon(err)
	return usr.HomeDir + string(filepath.Separator) + configName
}

// makeConfig creates a sample gondl.json template file, if it
// doesn't already exist
func makeConfig() {
	file := workingDirConfig()
	msg := "Sample config file created here"
	if _, err := os.Stat(file); err == nil {
		msg = "Config file already exists"
	} else {
		err := ioutil.WriteFile(file, []byte(configTemplate), 0644)
		if err != nil {
			msg = "Error - cannot create config file"
		}
	}
	fmt.Printf(configHelp, msg, file, file, homeDirConfig(), localDirConfig())
}

const configTemplate = `{
	"--apikey": "YOUR_KEY",
	"--trim_start": "2014-01-01",
	"--per_page": "10",
	"--url": true
}
`

const configHelp = `%v:
%v

You may edit it and use any of the long-form options (--options) in
it. The arguments you provide in the command line override any argument
in any of the config files.

Gondl will look for config files in this order:

Working Directory: %v
User Directory   : %v
Local Directory  : %v

`

//;D
