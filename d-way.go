package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"ibfd.org/d-way/doc"
	"ibfd.org/d-way/prc"
	"ibfd.org/d-way/rule"

	"ibfd.org/d-way/config"
)

const pathPrefix = "/d-way"

func main() {
	server := http.Server{Addr: ":" + cfg.GetPort()}
	log.Printf("d-way %s started on %s", version, server.Addr)
	matcher := cfg.GetMatcher()
	logMatcher(matcher)
	http.HandleFunc("/", docHandler(matcher))
	http.HandleFunc("/config/", configHandler(cfg.GetConfigData()))
	server.ListenAndServe()
}

func docHandler(matcher *rule.Matcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		switch r.Method {
		case "GET":
			if !strings.HasPrefix(r.URL.Path, pathPrefix) {
				notFound(w, r.URL.Path)
			} else {
				path := strings.TrimPrefix(r.URL.Path, pathPrefix)
				document := doc.NewDocument(path)
				rule := matcher.Match(document)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(200)
				fmt.Fprintf(w, "Got    : %s\n", document)
				fmt.Fprintf(w, "Matched: %s\n", rule.Regex)
				job := prc.NewJob(document, rule, r.Cookies())
				jobResult, err := prc.Exec(job)
				if err != nil {
					log.Printf("error %v", err)
				} else {
					err = output(w, jobResult.Reader)
					for _, result := range jobResult.Steps {
						log.Printf("Step %s took %s\n", result.Step, result.Duration)
					}
				}
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

func notFound(w http.ResponseWriter, url string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(404)
	w.Write([]byte(url))
}

func logMatcher(matcher *rule.Matcher) {
	log.Printf("%s", matcher.Comment)
	log.Printf("rules: ")
	for i, rule := range matcher.Rules {
		log.Printf("\trule[%d]:", i)
		log.Printf("\t\tregex: %s", rule.Regex)
		log.Printf("\t\tsteps: %s", strings.Join(rule.Steps, ", "))
	}
}

func logRequest(r *http.Request) {
	for k, v := range r.Header {
		log.Printf("key[%s] = %v\n", k, v)
	}
	for i, cookie := range r.Cookies() {
		log.Printf("cookie[%d] = %v\n", i, cookie)
	}
}
