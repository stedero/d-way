package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"ibfd.org/d-way/config"
)

var config *cfg.Config

func main() {
	server := http.Server{Addr: ":" + cfg.GetPort()}
	log.Printf("d-way %s started on %s", version, server.Addr)
	config = cfg.GetConfig()
	logConfig(config)
	http.HandleFunc("/d-way/", handle)
	server.ListenAndServe()
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := strings.TrimPrefix(r.URL.Path, "/d-way")
		log.Printf("We got : %s", path)
		log.Printf("Matched: %s", config.Match(path).Regex)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		fmt.Fprintf(w, "We got : %s\n", path)
		fmt.Fprintf(w, "Matched: %s", config.Match(path).Regex)
	default:
	}
}

func showForm(w http.ResponseWriter) {
	form := `<html>
				<head>
					<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
					<title>D-way - Document Gateway</title>
				</head>
				<body>
					<h1>D-way - Document Gateway</h1>
					<form action="/" method="post" enctype="multipart/form-data">
						<p>Provide a document path</p>
						<input type="text" name="url" size="100"><br><br>
						<input type="submit">
					</form>
				</body>
			</html>`
	w.Write([]byte(form))
}

func process(w http.ResponseWriter, r *http.Request) {
	logHeaders(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Thanks!"))
}

func logConfig(config *cfg.Config) {
	log.Printf("%s", config.Comment)
	log.Printf("rules: ")
	for i, rule := range config.Rules {
		log.Printf("\trule[%d]:", i)
		log.Printf("\t\tregex: %s", rule.Regex)
		log.Printf("\t\tprocesses: %s", strings.Join(rule.Processes, ", "))
	}
}

func logHeaders(r *http.Request) {
	for k, v := range r.Header {
		log.Printf("key[%s] = %v\n", k, v)
	}
}
