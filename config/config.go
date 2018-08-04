package cfg

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"ibfd.org/d-way/rule"
)

// defaultPort defines the port to use if not defined
// by the environment or on the command line
const defaultPort = "8080"
const defaultConfigFilePath = "config.json"

var configFilePath string
var matcher *rule.Matcher

func init() {
	flag.Parse()
	configFilePath = flag.Arg(0)
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("fail to read file %s: %v", configFilePath, err)
	}
	matcher, err = rule.NewMatcher(data)
	if err != nil {
		log.Fatalf("fail to unmarshal from file %s: %v", configFilePath, err)
	}
	log.Printf("loaded configuration from %s", configFilePath)
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
