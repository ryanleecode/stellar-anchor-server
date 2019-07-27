package middleware

import "net/http"

const applicationJSON = "application/json"
const textXML = "text/xml"

type ContentTypeMiddleware struct {
	next        http.Handler
	contentType string
}

func NewApplicationJSONMiddleware(next http.Handler) http.Handler {
	return &ContentTypeMiddleware{
		next:        next,
		contentType: applicationJSON,
	}
}

func NewTextXMLMiddleware(next http.Handler) http.Handler {
	return &ContentTypeMiddleware{
		next:        next,
		contentType: textXML,
	}
}

func NewContentTypeMiddleware(next http.Handler, contentType string) http.Handler {
	return &ContentTypeMiddleware{
		next:        next,
		contentType: contentType,
	}
}

func (m *ContentTypeMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", string(m.contentType))
	m.next.ServeHTTP(w, r)
}
