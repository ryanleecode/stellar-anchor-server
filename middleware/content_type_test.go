package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

func (h *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}

func TestContentTypeMiddleware_ServeHTTP_CallsNextServeHttp(t *testing.T) {
	h := &MockHandler{}

	w := httptest.NewRecorder()
	r := &http.Request{}
	h.On("ServeHTTP", w, r).Return()

	m := NewContentTypeMiddleware(h, "")

	m.ServeHTTP(w, r)

	h.AssertCalled(t, "ServeHTTP", w, r)
}

func TestNewApplicationJSONMiddleware(t *testing.T) {
	h := &MockHandler{}

	w := httptest.NewRecorder()
	r := &http.Request{}
	h.On("ServeHTTP", w, r).Return()

	m := NewApplicationJSONMiddleware(h)

	m.ServeHTTP(w, r)

	assert.Equal(t, "application/json", w.Header().Get("content-type"))
}

func TestNewTextXMLMiddleware(t *testing.T) {
	h := &MockHandler{}

	w := httptest.NewRecorder()
	r := &http.Request{}
	h.On("ServeHTTP", w, r).Return()

	m := NewTextXMLMiddleware(h)

	m.ServeHTTP(w, r)

	assert.Equal(t, "text/xml", w.Header().Get("content-type"))
}

func TestNewContentTypeMiddleware(t *testing.T) {
	h := &MockHandler{}

	w := httptest.NewRecorder()
	r := &http.Request{}
	h.On("ServeHTTP", w, r).Return()

	cType := "application/test"
	m := NewContentTypeMiddleware(h, cType)

	m.ServeHTTP(w, r)

	assert.Equal(t, cType, w.Header().Get("content-type"))
}
