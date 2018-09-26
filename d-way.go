package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/cors"
	"ibfd.org/d-way/cfg"
	"ibfd.org/d-way/doc"
	log "ibfd.org/d-way/log4u"
	"ibfd.org/d-way/prc"
	"ibfd.org/d-way/rule"
)

const appName = "d-way"
const pathPrefix = "/d-way"

func main() {
	defer cfg.CloseLog()
	cfg.SetUserAgent(appName + "/" + version)
	log.Infof("Starting %s %s on port %s", appName, version, cfg.GetPort())
	server := http.Server{Addr: ":" + cfg.GetPort()}
	matcher := cfg.GetMatcher()
	logMatcher(matcher)
	http.HandleFunc("/", withCORS(docHandler(matcher)))
	http.HandleFunc("/config/", withCORS(configHandler(cfg.GetConfigData())))
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
			rule := matcher.Match(path)
			src := getSource(rule, request)
			job := prc.NewJob(rule, request.Cookies())
			jobResult, err := prc.Exec(job, src)
			if err != nil {
				writer.Header().Set("Content-Type", "text/plain")
				writer.WriteHeader(500)
				fmt.Fprintf(writer, "Bad request: %v\n", err)
				log.Errorf("error %v", err)
			} else {
				writer.Header().Set("Allow", "GET, DELETE, OPTIONS, POST, PUT")
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
		case "OPTIONS":
			writer.Header().Set("Allow", "GET, DELETE, OPTIONS, POST, PUT")
			writer.WriteHeader(http.StatusOK)
		default:
		}
	}
}

func getSource(rule *rule.Rule, request *http.Request) *doc.Source {
	if rule.Steps[0] == "RESOLVE" {
		parts := strings.Split(request.URL.Path, "/")
		return doc.StringSource(parts[len(parts)-1])
	}
	return doc.StringSource(cleanPath(request.URL.Path))
}

func output(dst io.Writer, src io.ReadCloser) error {
	defer src.Close()
	_, err := io.Copy(dst, src)
	return err
}

func configHandler(configData []byte) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			writer.Header().Set("Allow", "GET, DELETE, OPTIONS, POST, PUT")
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(200)
			writer.Write(configData)
		default:
		}
	}
}

func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	options := cors.Options{
		Debug:            true,
		AllowedHeaders:   []string{"authorization"},
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true}
	cors := cors.New(options)
	return cors.Handler(handler).ServeHTTP
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
