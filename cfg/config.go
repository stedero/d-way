package cfg

import (
	"flag"
	"io"
	"io/ioutil"
	"os"

	log "ibfd.org/d-way/log4u"
	"ibfd.org/d-way/rule"
)

// defaultPort defines the port to use if not defined
// by the environment or on the command line
const defaultPort = "8080"
const defaultConfigFilePath = "config.json"
const defaultLogLevel = "DEBUG"

var configFilePath string
var logFilePath string
var matcher *rule.Matcher
var configData []byte
var version string
var logFile *os.File

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
	logFile = configureLogging(matcher.Logging)
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

func configureLogging(logConfig *rule.LogConfig) *os.File {
	var logFile *os.File
	var err error
	if logConfig == nil || logConfig.Filename == "" {
		log.SetLevel(defaultLogLevel)
	} else {
		logFile, err = os.Create(logConfig.Filename)
		if err != nil {
			log.Fatalf("failed to create file %s: %v", logConfig.Filename, err)
		}
		logger := io.MultiWriter(os.Stderr, logFile)
		log.SetLevel(logConfig.Level)
		log.SetOutput(logger)
	}
	return logFile
}

// CloseLog closes the log file
func CloseLog() {
	if logFile != nil {
		logFile.Close()
	}
}
