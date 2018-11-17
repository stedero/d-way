package doc

import (
	"io"
	"io/ioutil"
	"mime"
	"strings"
)

// Source describes an input source.
type Source struct {
	content  string
	reader   io.ReadCloser
	mimeType string
}

// StringSource creates an input source from a string.
func StringSource(path string) *Source {
	return &Source{content: path, mimeType: mime.TypeByExtension(extension(path))}
}

// ReaderSource creates an input source froma reader.
func ReaderSource(r io.ReadCloser) *Source {
	return &Source{reader: r}
}

// Path get the source path.
func (src *Source) Path() string {
	return src.String()
}

func (src *Source) String() string {
	if src.content == "" && src.reader != nil {
		defer src.reader.Close()
		b, _ := ioutil.ReadAll(src.reader)
		src.content = string(b)
	}
	return src.content
}

// Reader returns a reader.
func (src *Source) Reader() io.Reader {
	if src.reader != nil {
		return src.reader
	}
	return strings.NewReader(src.content)
}

// MimeType returns the source mime type.
// Returns an empty string if the source content
// does not refer to a file or if the extension is unknown.
func (src *Source) MimeType() string {
	if src.mimeType != "" {
		return src.mimeType
	}
	return mime.TypeByExtension(extension(src.String()))
}

// Extension gives the extension if the source refers to a file.
func (src *Source) Extension() string {
	return extension(src.Path())
}

func extension(path string) string {
	dotpos := strings.LastIndex(path, ".")
	if dotpos < 0 {
		return ""
	}
	return path[dotpos:]
}
