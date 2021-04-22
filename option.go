package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const (
	headerContentTypeKey = "Content-Type"
)

type requestOptions struct {
	headers H
	params  P
	bodyfn  func() (io.Reader, error)
}

var defaultRequestOptions = requestOptions{
	headers: nil,
	params:  nil,
	bodyfn: func() (io.Reader, error) {
		return nil, nil
	},
}

func (o *requestOptions) newRequest(method, url string) (req *http.Request, err error) {
	// body
	br, err := o.bodyfn()
	if err != nil {
		return
	}

	req, err = http.NewRequest(method, url, br)
	if err != nil {
		return
	}

	for k, vs := range o.headers {
		req.Header[k] = append(req.Header[k], vs...)
	}

	if req.URL.RawQuery != "" && !strings.HasSuffix(req.URL.RawQuery, "&") {
		req.URL.RawQuery += "&"
	}

	req.URL.RawQuery += o.params.String()

	return req, nil
}

type Option func(*requestOptions)

func WithParams(p P) Option {
	return func(o *requestOptions) {
		o.params = p
	}
}

func WithHeaders(headers H) Option {
	return func(o *requestOptions) {
		if o.headers == nil {
			o.headers = headers
			return
		}
		o.headers.override(headers)
	}
}

func WithBodyJson(v interface{}) Option {
	return func(o *requestOptions) {
		o.bodyfn = func() (io.Reader, error) {
			bs, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			return bytes.NewBuffer(bs), nil
		}

		if o.headers == nil {
			o.headers = H{}
		}

		if _, ok := o.headers[headerContentTypeKey]; !ok {
			o.headers.Set(headerContentTypeKey, "application/json")
		}

	}
}

func WithBodyBytes(b []byte) Option {
	return func(o *requestOptions) {
		o.bodyfn = func() (io.Reader, error) {
			return bytes.NewBuffer(b), nil
		}
	}
}

func WithBodyReader(r io.Reader) Option {
	return func(o *requestOptions) {
		o.bodyfn = func() (io.Reader, error) {
			return r, nil
		}
	}
}
