package requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
}

func newResponse(resp *http.Response) *Response {
	return &Response{Response: resp}
}

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
