package requests

import "net/http"

type Request struct {
	c *http.Client
}

func (r *Request) Request(method, url string, opts ...Option) (response *Response, err error) {
	options := defaultRequestOptions
	for _, opt := range opts {
		opt(&options)
	}

	req, err := options.newRequest(method, url)
	if err != nil {
		return
	}
	resp, err := r.c.Do(req)
	if err != nil {
		return
	}

	return newResponse(resp), nil
}

func (r *Request) Get(url string, opts ...Option) (resp *Response, err error) {
	return r.Request(http.MethodGet, url, opts...)
}

func (r *Request) Post(url string, opts ...Option) (*Response, error) {
	return r.Request(http.MethodPost, url, opts...)
}

func (r *Request) Delete(url string, opts ...Option) (*Response, error) {
	return r.Request(http.MethodDelete, url, opts...)
}

func (r *Request) Put(url string, opts ...Option) (*Response, error) {
	return r.Request(http.MethodPut, url, opts...)
}

func (r *Request) Patch(url string, opts ...Option) (*Response, error) {
	return r.Request(http.MethodPatch, url, opts...)
}

func (r *Request) Head(url string, opts ...Option) (*Response, error) {
	return r.Request(http.MethodHead, url, opts...)
}
