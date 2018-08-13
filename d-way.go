package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"ibfd.org/d-way/config"
	"ibfd.org/d-way/doc"
	"ibfd.org/d-way/log"
	"ibfd.org/d-way/prc"
	"ibfd.org/d-way/rule"
)

const pathPrefix = "/d-way"

func main() {
	logConfig := cfg.GetLogConfig()
	if logConfig.Filename == "" {
		log.SetLevel("DEBUG")
	} else {
		logFile, err := os.Create(logConfig.Filename)
		if err != nil {
			log.Fatalf("fail to create file %s: %v", logConfig.Filename, err)
		}
		logger := io.MultiWriter(os.Stderr, logFile)
		log.SetLevel(logConfig.Level)
		log.SetOutput(logger)
		defer logFile.Close()
	}
	log.Info("Starting d-way on port %s", cfg.GetPort())
	server := http.Server{Addr: ":" + cfg.GetPort()}
	matcher := cfg.GetMatcher()
	logMatcher(matcher)
	http.HandleFunc("/", docHandler(matcher))
	http.HandleFunc("/config/", configHandler(cfg.GetConfigData()))
	server.ListenAndServe()
}

func docHandler(matcher *rule.Matcher) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			path := cleanPath(request.URL.Path)
			if path == "/favicon.ico" {
				writer.WriteHeader(404)
				return
			}
			logRequest(request)
			document := doc.NewDocument(path)
			rule := matcher.Match(document)
			job := prc.NewJob(document, rule, request.Cookies())
			jobResult, err := prc.Exec(job)
			if err != nil {
				writer.Header().Set("Content-Type", "text/plain")
				writer.WriteHeader(500)
				fmt.Fprintf(writer, "Bad request: %v\n", err)
				log.Errorf("error %v", err)
			} else {
				writer.Header().Set("Content-Type", jobResult.ContentType())
				writer.WriteHeader(200)
				err = output(writer, jobResult.Reader())
				if err != nil {
					log.Errorf("error %v", err)
				}
				for _, result := range jobResult.Steps() {
					log.Debugf("Step %s took %s\n", result.Step, result.Duration)
				}
				logResponse(jobResult.Response())
			}
		default:
		}
	}
}

func output(dst io.Writer, src io.ReadCloser) error {
	defer src.Close()
	_, err := io.Copy(dst, src)
	return err
}

func configHandler(configData []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(configData)
		default:
		}
	}
}

func cleanPath(path string) string {
	return strings.TrimPrefix(path, pathPrefix)
}

func notFound(w http.ResponseWriter, url string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(404)
	w.Write([]byte(url))
}

func logMatcher(matcher *rule.Matcher) {
	log.Debugf("%s", matcher.Comment)
	log.Debug("rules:")
	for i, rule := range matcher.Rules {
		log.Debugf("\trule[%d]:", i)
		log.Debugf("\t\tregex: %s", rule.Regex)
		log.Debugf("\t\tsteps: %s", strings.Join(rule.Steps, ", "))
	}
}

func logRequest(r *http.Request) {
	log.Debug("===========================================================")
	for k, v := range r.Header {
		log.Debugf("key[%s] = %v\n", k, v)
	}
	for i, cookie := range r.Cookies() {
		log.Debugf("cookie[%d] = %v\n", i, cookie)
	}
}

func logResponse(r *http.Response) {
	if r != nil {
		for k, v := range r.Header {
			log.Debugf("key[%s] = %v\n", k, v)
		}
		for i, cookie := range r.Cookies() {
			log.Debugf("cookie[%d] = %v\n", i, cookie)
		}
	}
}
