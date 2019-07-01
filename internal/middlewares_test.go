package internal_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
)

var uuidRegex = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func TestContentTypeMiddleware(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
	}

	assert := assert.New(t)

	handler := ContentType(http.HandlerFunc(mockHandler))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)

	assert.EqualValues(res.Header.Get("Content-Type"), "application/json")
}

func TestRequestIDContextMiddleware(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("id")
		fmt.Fprintf(w, "%s", id)
	}

	assert := assert.New(t)

	handler := IDContext(http.HandlerFunc(mockHandler))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)

	id := string(body)

	assert.Regexp(uuidRegex, id)
}

func TestRequestMethodContextMiddleware(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		method := r.Context().Value("method")
		fmt.Fprintf(w, "%s", method)
	}

	assert := assert.New(t)

	const apiVersion = "v1"
	const route = "testing"

	handler := MethodContext(
		http.HandlerFunc(mockHandler))
	router := mux.NewRouter().PathPrefix(
		fmt.Sprintf("/%s", apiVersion)).Subrouter()
	router.HandleFunc(fmt.Sprintf("/%s", route), handler.ServeHTTP)
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/%s/%s", ts.URL, apiVersion, route))
	assert.NoError(err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)

	method := string(body)

	assert.EqualValues(fmt.Sprintf("%s.get", route), method)
}

func TestResponseWriterMiddleware(t *testing.T) {
	testMessage := "testing"
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testMessage)
	}

	writer := httptest.NewRecorder()
	writerFactory := func(http.ResponseWriter, *http.Request) http.ResponseWriter {
		return writer
	}

	assert := assert.New(t)

	handler := NewResponseWriter(
		writerFactory)(http.HandlerFunc(mockHandler))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	_, err := http.Get(ts.URL)
	assert.NoError(err)

	resp := writer.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err)

	message := string(body)

	assert.EqualValues(testMessage, message)
}

