// tests/testutil/request.go
package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// Request is a fluent builder for test HTTP requests.
type Request struct {
	router  *gin.Engine
	method  string
	path    string
	body    interface{}
	headers map[string]string
}

// Do executes the request and returns a ResponseRecorder.
func (r *Request) Do() *httptest.ResponseRecorder {
	var bodyBytes []byte
	if r.body != nil {
		bodyBytes, _ = json.Marshal(r.body)
	}

	req, _ := http.NewRequest(r.method, r.path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	r.router.ServeHTTP(w, req)
	return w
}

// WithAuth adds a Bearer token header.
func (r *Request) WithAuth(token string) *Request {
	r.headers["Authorization"] = "Bearer " + token
	return r
}

// WithBody sets the request JSON body.
func (r *Request) WithBody(body interface{}) *Request {
	r.body = body
	return r
}

// GET builds a GET request.
func GET(router *gin.Engine, path string) *Request {
	return &Request{router: router, method: "GET", path: path, headers: map[string]string{}}
}

// POST builds a POST request.
func POST(router *gin.Engine, path string) *Request {
	return &Request{router: router, method: "POST", path: path, headers: map[string]string{}}
}

// PUT builds a PUT request.
func PUT(router *gin.Engine, path string) *Request {
	return &Request{router: router, method: "PUT", path: path, headers: map[string]string{}}
}

// DELETE builds a DELETE request.
func DELETE(router *gin.Engine, path string) *Request {
	return &Request{router: router, method: "DELETE", path: path, headers: map[string]string{}}
}

// ── JSON body helpers ─────────────────────────────────

// ParseBody decodes the response JSON into v.
func ParseBody(w *httptest.ResponseRecorder, v interface{}) error {
	return json.NewDecoder(w.Body).Decode(v)
}

// BodyString returns a map of the response JSON body for quick assertions.
func BodyString(w *httptest.ResponseRecorder) map[string]interface{} {
	var m map[string]interface{}
	json.NewDecoder(w.Body).Decode(&m)
	return m
}

// DataField returns the "data" field from a standard API response.
func DataField(w *httptest.ResponseRecorder) map[string]interface{} {
	body := BodyString(w)
	if data, ok := body["data"].(map[string]interface{}); ok {
		return data
	}
	return nil
}

// DataString returns a string field from response.data.<key>.
func DataString(w *httptest.ResponseRecorder, key string) string {
	data := DataField(w)
	if data == nil {
		return ""
	}
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}

// ErrorField returns the "error" field from a standard error response.
func ErrorField(w *httptest.ResponseRecorder) string {
	body := BodyString(w)
	if err, ok := body["error"].(string); ok {
		return err
	}
	return ""
}

// J is a shorthand for building JSON bodies inline.
func J(pairs ...interface{}) map[string]interface{} {
	if len(pairs)%2 != 0 {
		panic(fmt.Sprintf("testutil.J: odd number of arguments: %d", len(pairs)))
	}
	m := make(map[string]interface{}, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		m[fmt.Sprint(pairs[i])] = pairs[i+1]
	}
	return m
}
