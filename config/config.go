package cfg

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

// defaultPort defines the port to use if not defined
// by the environment or on the command line
const defaultPort = "8080"
const defaultConfigFilePath = "config.json"

// Rule defines process rules
type Rule struct {
	Regex     string `json:"regex"`
	Regexc    *regexp.Regexp
	Processes []string `json:"process"`
}

// Config defines the structure of the config.json file
type Config struct {
	Comment string  `json:"comment"`
	Rules   []*Rule `json:"rules"`
}

var configFilePath string
var config Config

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
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("fail to unmarshal from file %s: %v", configFilePath, err)
	}
	for _, rule := range config.Rules {
		rule.Regexc = regexp.MustCompile(rule.Regex)
	}
	log.Printf("loaded configuration from %s", configFilePath)
}

// GetConfig return the configuration
func GetConfig() *Config {
	return &config
}

// GetPort returns the port to use for this service
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

// Match finds a rule that matches path
func (config *Config) Match(path string) *Rule {
	for _, rule := range config.Rules {
		if rule.Regexc.Match([]byte(path)) {
			return rule
		}
	}
	return config.Rules[len(config.Rules)-1]
}
