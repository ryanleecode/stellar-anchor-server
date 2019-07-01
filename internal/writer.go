package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// ResponseWriter intercepts its http.ResponseWriter and injects Definition metadata
// along with anything that is passed in the Write function
type ResponseWriter struct {
	ResponseProperties Properties
	Writer             http.ResponseWriter
	Request            *http.Request
}

// Header returns the original response writer headers
func (interceptor ResponseWriter) Header() http.Header {
	return interceptor.Writer.Header()
}

// Write writes additional response metadata to the original http.ResponseWriter
func (interceptor ResponseWriter) Write(b []byte) (int, error) {
	responseProperties := interceptor.ResponseProperties
	bytesWritten := 0

	ctx := interceptor.Request.Context()

	requestID, ok := ctx.Value("id").(string)
	if !ok {
		requestID = ""
	}
	method, ok := ctx.Value("method").(string)
	if !ok {
		requestID = ""
	}

	err := func() error {
		response := &Definition{
			Base: Base{
				APIVersion: responseProperties.APIVersion,
				ID:         requestID,
				Method:     method,
			},
		}
		err := json.Unmarshal(b, &response)
		if err != nil {
			return err
		}

		responseBytes := new(bytes.Buffer)
		err = json.NewEncoder(responseBytes).Encode(response)
		if err != nil {
			return err
		}

		bytesWritten, err = interceptor.Writer.Write(responseBytes.Bytes())

		return err
	}()

	return bytesWritten, err
}

// WriteHeader writes a header to the original http.ResponseWriter
func (interceptor ResponseWriter) WriteHeader(statusCode int) {
	interceptor.Writer.WriteHeader(statusCode)
}
