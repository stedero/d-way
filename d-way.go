package main

import (
	"html"
	"log"
	"net/http"

	"ibfd.org/d-way/config"
)

func main() {
	server := http.Server{Addr: ":" + config.GetPort()}
	log.Printf("d-way %s started on %s", version, server.Addr)
	http.HandleFunc("/", handle)
	server.ListenAndServe()
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %s: %s", r.Method, r.RequestURI)
	switch r.Method {
	case "GET":
		log.Printf("We got: %q", html.EscapeString(r.URL.Path))
		showForm(w)
	case "POST":
		process(w, r)
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

func logHeaders(r *http.Request) {
	for k, v := range r.Header {
		log.Printf("key[%s] = %v\n", k, v)
	}
}
