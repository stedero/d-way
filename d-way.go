package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
	http.HandleFunc("/", withCORS(docHandler(matcher, fmt.Sprintf("max-age=%d", matcher.CacheMaxAgeSeconds))))
	http.HandleFunc("/config/", withCORS(configHandler(cfg.GetConfigData())))
	server.ListenAndServe()
}

func docHandler(matcher *rule.Matcher, maxAgeHeader string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			total := timer()
			path := cleanPath(request.URL.Path)
			mimeType := request.Header.Get("Accept")
			rule := matcher.Match(path, mimeType)
			src := doc.StringSource(path)
			job := prc.NewJob(rule, request.Cookies(), requestModSinceDate(request))
			jobResult, err := prc.Exec(job, src)
			if err == nil {
				sendSuccess(jobResult, maxAgeHeader, writer)
			} else {
				sendError(err, writer)
			}
			log.Debugf("rule %s took: %v", rule.Name, total())
		case "OPTIONS":
			writer.Header().Set("Server", cfg.GetUserAgent())
			writer.Header().Set("Allow", "GET, OPTIONS")
			writer.WriteHeader(http.StatusOK)
		default:
		}
	}
}

func sendError(err error, writer http.ResponseWriter) {
	statusCode, msg := statusAndMessage(err)
	writer.Header().Set("Server", cfg.GetUserAgent())
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(statusCode)
	fmt.Fprintf(writer, msg)
	log.Errorf("%d - %s", statusCode, msg)
}

func sendSuccess(jobResult *prc.JobResult, maxAgeHeader string, writer http.ResponseWriter) {
	writer.Header().Set("Allow", "GET, OPTIONS")
	writer.Header().Set("Server", cfg.GetUserAgent())
	switch jobResult.ResultType() {
	case act.Action:
		sendAction(jobResult, writer)
	case act.Content:
		sendContent(jobResult, maxAgeHeader, writer)
	}
	logJobResult(jobResult)
}

func sendAction(jobResult *prc.JobResult, writer http.ResponseWriter) {
	statusCode := jobResult.StatusCode()
	if statusCode == http.StatusFound {
		writer.Header().Set("Location", jobResult.Content())
	}
	writer.WriteHeader(statusCode)
}

func sendContent(jobResult *prc.JobResult, maxAgeHeader string, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", jobResult.ContentType())
	lastModifiedDate := jobResult.LastModifiedDate()
	if lastModifiedDate != "" {
		writer.Header().Set("Cache-Control", maxAgeHeader)
		writer.Header().Set("Last-Modified", lastModifiedDate)
	}
	statusCode := jobResult.StatusCode()
	writer.WriteHeader(statusCode)
	if statusCode != http.StatusNotModified {
		err := output(writer, jobResult.Reader())
		if err != nil {
			log.Errorf("error %v", err)
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
			writer.Header().Set("Server", cfg.GetUserAgent())
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
		AllowedOrigins:   []string{"http://localhost:4200", "http://steef.ibfd.org:4200"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowCredentials: true}
	cors := cors.New(options)
	return cors.Handler(handler).ServeHTTP
}

func requestModSinceDate(request *http.Request) *time.Time {
	modSinceDateStr := request.Header.Get("If-Modified-Since")
	msd, err := time.Parse(cfg.LastModifiedDateFormat, modSinceDateStr)
	if err == nil {
		return &msd
	}
	return nil
}

func cleanPath(path string) string {
	return strings.TrimSuffix(strings.TrimPrefix(path, pathPrefix), "/")
}

func timer() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Now().Sub(start)
	}
}

func logMatcher(matcher *rule.Matcher) {
	log.Infof("%s", matcher.Comment)
	log.Infof("rules:")
	for i, rule := range matcher.Rules {
		log.Infof("\trule[%d]:", i)
		log.Infof("\t\tname : %s", rule.Name)
		log.Infof("\t\tregex: %s", rule.Regex)
		if rule.MimeType != "" {
			log.Infof("\t\tmimeType: %s", rule.MimeType)
		}
		log.Infof("\t\tsteps: %s", strings.Join(rule.Steps, ", "))
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

func logJobResult(jobResult *prc.JobResult) {
	for _, result := range jobResult.Steps() {
		log.Debugf("step %s took %s\n", result.Step, result.Duration)
	}
}
