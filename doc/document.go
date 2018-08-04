package doc

// Document describes a document
type Document struct {
	Path string
}

// NewDocument creates a document
func NewDocument(path string) *Document {
	return &Document{path}
}

func (document *Document) String() string {
	return document.Path
}
