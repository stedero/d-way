package cfg

import (
	"flag"
	"io/ioutil"
	"os"

	"ibfd.org/d-way/log"
	"ibfd.org/d-way/rule"
)

// defaultPort defines the port to use if not defined
// by the environment or on the command line
const defaultPort = "8080"
const defaultConfigFilePath = "config.json"

var configFilePath string
var logFilePath string
var matcher *rule.Matcher
var configData []byte
var version string

func init() {
	var err error
	flag.Parse()
	configFilePath = flag.Arg(0)
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}
	configData, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("fail to read file %s: %v", configFilePath, err)
	}
	matcher, err = rule.NewMatcher(configData)
	if err != nil {
		log.Fatalf("fail to unmarshal from file %s: %v", configFilePath, err)
	}
}

// GetMatcher returns the rule matcher.
func GetMatcher() *rule.Matcher {
	return matcher
}

// GetPort returns the port to use for this service.
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		flag.Parse()
		port = flag.Arg(1)
		if port == "" {
			port = defaultPort
		}
	}
	return port
}

// GetConfigData provides the content of the config file.
func GetConfigData() []byte {
	return configData
}

// GetLogConfig returns the logging configuration.
func GetLogConfig() *rule.LogConfig {
	if matcher.Logging != nil {
		return matcher.Logging
	}
	return &rule.LogConfig{}
}
