package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/cors"
	"ibfd.org/d-way/act"
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
			rule := matcher.Match(path)
			src := getSource(rule, path)
			job := prc.NewJob(rule, request.Cookies())
			jobResult, err := prc.Exec(job, src)
			if err != nil {
				statusCode, msg := statusAndMessage(err)
				writer.Header().Set("Content-Type", "text/plain")
				writer.WriteHeader(statusCode)
				fmt.Fprintf(writer, msg)
				log.Errorf("%d - %s", statusCode, msg)
			} else {
				writer.Header().Set("Allow", "GET, OPTIONS")
				writer.Header().Set("Content-Type", jobResult.ContentType())
				writer.WriteHeader(200)
				err = output(writer, jobResult.Reader())
				if err != nil {
					log.Errorf("error %v", err)
				}
				for _, result := range jobResult.Steps() {
					log.Debugf("%s took %s\n", result.Step, result.Duration)
				}
				total := jobResult.Total()
				log.Debugf("%s took %s\n", total.Step, total.Duration)
			}
		case "OPTIONS":
			writer.Header().Set("Allow", "GET, OPTIONS")
			writer.WriteHeader(http.StatusOK)
		default:
		}
	}
}

func statusAndMessage(err error) (int, string) {
	savedErr := err
	if err, ok := err.(*act.ActionError); ok {
		return err.StatusCode, err.Msg
	}
	return http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", savedErr)
}

func getSource(rule *rule.Rule, path string) *doc.Source {
	if rule.Steps[0] == "RESOLVE" {
		parts := strings.Split(path, "/")
		return doc.StringSource(parts[len(parts)-1])
	}
	return doc.StringSource(path)
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
			writer.Header().Set("Allow", "GET, OPTIONS")
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
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowCredentials: true}
	cors := cors.New(options)
	return cors.Handler(handler).ServeHTTP
}

func cleanPath(path string) string {
	return strings.TrimSuffix(strings.TrimPrefix(path, pathPrefix), "/")
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
