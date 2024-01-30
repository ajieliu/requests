package requests

type Client interface {
	Get(url string, opts ...Option) (resp *Response, err error)
	Post(url string, opts ...Option) (resp *Response, err error)
	Delete(url string, opts ...Option) (resp *Response, err error)
	Put(url string, opts ...Option) (resp *Response, err error)
	Patch(url string, opts ...Option) (resp *Response, err error)
	Options(url string, opts ...Option) (resp *Response, err error)

	Request(method, url string, opts ...Option) (resp *Response, err error)
}
