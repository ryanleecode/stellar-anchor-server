package internal_test

import (
	"github.com/drdgvhbh/stellar-fi-anchor/internal"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var uuidRegex = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func TestContentTypeMiddleware(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
	}

	assert := assert.New(t)

	handler := internal.ContentType(http.HandlerFunc(mockHandler))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)

	assert.EqualValues(res.Header.Get("Content-Type"), "application/json")
}
