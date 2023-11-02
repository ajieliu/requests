package requests

import (
	"io"
	"net/http"
	"net/url"
)

type H http.Header

func (h H) Add(k, v string) H {
	http.Header(h).Add(k, v)
	return h
}

func (h H) Set(k, v string) H {
	http.Header(h).Set(k, v)
	return h
}

func (h H) Del(k string) H {
	http.Header(h).Del(k)
	return h
}

func (h H) override(hs ...H) H {
	for _, header := range hs {
		for k, vs := range header {
			h[k] = vs
		}
	}
	return h
}

type P url.Values

func (p P) Get(key string) string {
	return url.Values(p).Get(key)
}

func (p P) Set(key, value string) P {
	url.Values(p).Set(key, value)
	return p
}

func (p P) Add(key, value string) P {
	url.Values(p).Add(key, value)
	return p
}

func (p P) Del(key string) P {
	delete(p, key)
	return p
}

func (p P) String() string {
	return url.Values(p).Encode()
}

type File interface {
	Name() string
	io.ReadCloser
}

type mFile struct {
	name string
	io.ReadCloser
}

func NewRequestFile(name string, body io.ReadCloser) File {
	return &mFile{
		name:       name,
		ReadCloser: body,
	}
}

func (f *mFile) Name() string {
	return f.name
}
