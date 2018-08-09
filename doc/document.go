package doc

import (
	"mime"
	"strings"
)

// Document describes a document
type Document struct {
	path     string
	mimeType string
}

// NewDocument creates a document
func NewDocument(path string) *Document {
	return &Document{path, mime.TypeByExtension(extension(path))}
}

func (document *Document) String() string {
	return document.path
}

// Path get the document path.
func (document *Document) Path() string {
	return document.path
}

// MimeType get the document mime type.
func (document *Document) MimeType() string {
	return document.mimeType
}

func extension(path string) string {
	dotpos := strings.LastIndex(path, ".")
	if dotpos < 0 {
		return ""
	}
	return path[dotpos:]
}
