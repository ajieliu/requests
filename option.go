package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	headerContentTypeKey = "Content-Type"
)

type requestOptions struct {
	headers H
	params  P
	bodyfn  func() (io.Reader, error)
	ctx     context.Context

	onBeforeRequestFns []func(req *http.Request) error
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

	ctx := o.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	req, err = http.NewRequestWithContext(ctx, method, url, br)
	if err != nil {
		return
	}

	for k, vs := range o.headers {
		if strings.ToLower(k) == "host" {
			req.Host = vs[len(vs)-1]
		}
		req.Header[k] = append(req.Header[k], vs...)
	}

	if req.URL.RawQuery != "" && !strings.HasSuffix(req.URL.RawQuery, "&") {
		req.URL.RawQuery += "&"
	}

	req.URL.RawQuery += o.params.String()

	for _, fn := range o.onBeforeRequestFns {
		err = fn(req)
		if err != nil {
			return
		}
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
		if o.headers == nil {
			o.headers = headers
			return
		}
		o.headers.override(headers)
	}
}

func WithHeader(key string, v ...string) Option {
	return func(o *requestOptions) {
		if o.headers == nil {
			o.headers = H{}
		}
		o.headers[key] = v
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

func WithFormUrlEncoded(fields map[string][]string) Option {
	return func(o *requestOptions) {
		o.bodyfn = func() (io.Reader, error) {
			WithHeaders(H{}.Add("Content-Type", "application/x-www-form-urlencoded"))(o)
			return strings.NewReader(url.Values(fields).Encode()), nil
		}
	}
}

func WithForm(fields map[string][]string, files map[string]File) Option {
	return func(o *requestOptions) {
		o.bodyfn = func() (io.Reader, error) {
			body := new(bytes.Buffer)
			mw := multipart.NewWriter(body)
			defer mw.Close()

			// write fields
			for k, vs := range fields {
				for _, v := range vs {
					if err := mw.WriteField(k, v); err != nil {
						return nil, err
					}
				}
			}

			// write files
			for k, fh := range files {
				w, err := mw.CreateFormFile(k, fh.Name())
				if err != nil {
					return nil, err
				}
				io.Copy(w, fh)
				fh.Close()
			}

			// set Content-Type
			if o.headers == nil {
				o.headers = H{}
			}

			o.headers.Set(headerContentTypeKey, mw.FormDataContentType())
			return body, nil
		}
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *requestOptions) {
		o.ctx = ctx
	}
}

// WithOnBeforeRequest add a function to be called before request.
func WithOnBeforeRequest(fn func(req *http.Request) (err error)) Option {
	return func(options *requestOptions) {
		if fn == nil {
			return
		}
		options.onBeforeRequestFns = append(options.onBeforeRequestFns, fn)
	}
}
