package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
}

func newResponse(resp *http.Response) *Response {
	return &Response{Response: resp}
}

// Json unmarshal body with json
func (r *Response) Json(v interface{}) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(body, v)
}

func (r *Response) String() string {
	return fmt.Sprintf("%v", r.Response)
}

// CloseBodySilently close body silently
func (r *Response) CloseBodySilently() {
	_ = r.Body.Close()
}

// WriteTo write body to the given writer
func (r *Response) WriteTo(w io.Writer) (int64, error) {
	return io.Copy(w, r.Body)
}
