package main

import (
	"bitbucket.org/kardianos/osext"
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
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, err := osext.ExecutableFolder()
	panicon(err)
	return dir + configName
}

// workingDirConfig returns the path to the working directory config
// file
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
	if exist(file) {
		msg = "Config file already exists"
	} else {
		err := ioutil.WriteFile(file, []byte(configTemplate), 0644)
		if err != nil {
			msg = "Error - cannot create config file"
		}
	}
	fmt.Printf(configHelp, msg, file)
}

// showConfig shows information about the config files.
func showConfig() {
	f1, f2, f3 := "Not Found", "Not Found", "Not Found"
	if exist(workingDirConfig()) {
		f1 = "Found"
	}
	if exist(homeDirConfig()) {
		f2 = "Found"
	}
	if exist(localDirConfig()) {
		f3 = "Found"
	}
	fmt.Printf(configInfo,
		f1, workingDirConfig(),
		f2, homeDirConfig(),
		f3, localDirConfig(),
	)
}

func exist(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
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

Run 'gondl --config' for more information about config files.
`

const configInfo = `
Gondl will look for config files in three folders. The working 
directory, the user's home directory and the local (executable's) 
directory.

Values in the working directory will have precedence over values in 
the home directory, and values in the home directory 
will have precedence over values in the local directory.

  Working Directory: (%v)
  %v

  User Directory: (%v)
  %v

  Local Directory: (%v)
  %v

`

//;D
