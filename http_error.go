package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPError is returned whenever rest gets a non-2xx response from
// the server.
type HTTPError struct {
	// URL is the url that the request was sent to
	URL string
	// Body is the body of the response
	Body []byte
	// StatusCode is the http status code of the response
	StatusCode int
}

// Error satisfies the error interface
func (e HTTPError) Error() string {
	return fmt.Sprintf("rest: http request to %s returned status code %d", e.URL, e.StatusCode)
}

// newHTTPError returns an HTTPError based on the given response. It
// may return a different error if there was a problem reading the response
// body.
func newHTTPError(res *http.Response) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("rest: Unexpected error reading response body: %s", err.Error())
	}
	return HTTPError{
		URL:        res.Request.URL.String(),
		Body:       body,
		StatusCode: res.StatusCode,
	}
}
