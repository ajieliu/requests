package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	urlutil "net/url"
)

type requestOptions struct {
	headers H
	params  P
	bodyfn  func() (io.Reader, error)
}

var defaultRequestOptions = requestOptions{
	headers: H{},
	params:  P{},
	bodyfn: func() (io.Reader, error) {
		return nil, nil
	},
}

func (o *requestOptions) newRequest(method, url string) (req *http.Request, err error) {
	// params
	u, err := urlutil.Parse(url)
	if err != nil {
		return
	}
	for k, vs := range o.params {
		for _, v := range vs {
			u.Query().Add(k, v)
		}
	}

	// body
	br, err := o.bodyfn()
	if err != nil {
		return
	}

	req, err = http.NewRequest(method, u.RequestURI(), br)
	if err != nil {
		return
	}

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
		o.headers = headers
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
