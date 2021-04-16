package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithBodyBytes(t *testing.T) {
	testcases := []struct {
		data []byte
	}{
		{[]byte("name")},
		{[]byte{}},
	}

	for i, tc := range testcases {
		opts := defaultRequestOptions
		o := WithBodyBytes(tc.data)
		o(&opts)
		r, err := opts.bodyfn()
		assert.Nil(t, err, i)
		data, err := ioutil.ReadAll(r)
		assert.Nil(t, err, i)
		assert.Equal(t, tc.data, data, i)
	}
}

func TestWithBodyJson(t *testing.T) {
	testcases := []struct {
		data   interface{}
		expect []byte
		err    error
	}{
		{map[string]interface{}{"name": "zhangz", "age": 12}, []byte("{\"age\":12,\"name\":\"zhangz\"}"), nil},
		{struct {
			Name    string `json:"name"`
			Address string `json:"address"`
			Age     int
		}{"zhangz", "abc", 21}, []byte("{\"name\":\"zhangz\",\"address\":\"abc\",\"Age\":21}"), nil},
		{nil, []byte("null"), nil},
		{1, []byte("1"), nil},
		{true, []byte("true"), nil},
		{func() {}, []byte("true"), &json.UnsupportedTypeError{Type: reflect.TypeOf(func() {})}},
	}

	for i, tc := range testcases {
		opts := defaultRequestOptions
		WithBodyJson(tc.data)(&opts)

		r, err := opts.bodyfn()
		assert.Equal(t, tc.err, err, i)
		if err != nil {
			return
		}

		data, err := ioutil.ReadAll(r)
		assert.Nil(t, err, i)
		assert.Equal(t, string(tc.expect), string(data), i)
	}
}

func TestWithBodyReader(t *testing.T) {
	testcases := []struct {
		r      io.Reader
		expect []byte
	}{
		{bytes.NewReader([]byte{1, 2, 12}), []byte{1, 2, 12}},
		{strings.NewReader("this is a test case"), []byte("this is a test case")},
		{strings.NewReader(""), []byte{}},
	}

	for i, tc := range testcases {
		opts := defaultRequestOptions
		WithBodyReader(tc.r)(&opts)

		r, err := opts.bodyfn()
		assert.Nil(t, err, i)
		data, err := ioutil.ReadAll(r)
		assert.Nil(t, err, i)
		assert.Equal(t, tc.expect, data, i)
	}
}

func TestWithHeaders(t *testing.T) {
	testcases := []struct {
		h H
	}{
		{nil},
		{H{}},
		{H{}.Add("key", "value")},
	}

	for i, tc := range testcases {
		opts := defaultRequestOptions
		WithHeaders(tc.h)(&opts)
		assert.Equal(t, tc.h, opts.headers, i)
	}
}

func TestWithParams(t *testing.T) {
	testcases := []struct {
		p P
	}{
		{nil},
		{P{}},
		{P{}.Add("name", "liu")},
	}

	for i, tc := range testcases {
		opts := defaultRequestOptions
		WithParams(tc.p)(&opts)

		assert.Equal(t, tc.p, opts.params, i)
	}
}

func TestNewRequest(t *testing.T) {
	testcases := []struct {
		method        string
		url           string
		opts          requestOptions
		expectParams  P
		expectBody    []byte
		expectHeaders H
		expectErr     error
	}{
		{
			"GET",
			"http://localhost/a",
			requestOptions{
				H{}.Add("key", "value"),
				P{}.Add("name", "l"),
				func() (io.Reader, error) {
					return bytes.NewReader([]byte{1}), nil
				},
			},
			P{}.Add("name", "l"),
			[]byte{1},
			H{}.Add("key", "value"),
			nil,
		},
		{
			"GET",
			"http://localhost/a?name=test&age=12",
			requestOptions{
				H{}.Add("key", "value"),
				P{}.Add("name", "l"),
				func() (io.Reader, error) {
					return bytes.NewReader([]byte{1}), nil
				},
			},
			P{}.Add("name", "test").Add("name", "l").Add("age", "12"),
			[]byte{1},
			H{}.Add("key", "value"),
			nil,
		},
		{
			"GET",
			"http://localhost/a",
			requestOptions{
				H{}.Add("key", "value"),
				P{}.Add("name", "l"),
				func() (io.Reader, error) {
					return nil, errors.New("test error")
				},
			},
			P{}.Add("name", "test"),
			[]byte{1},
			H{}.Add("key", "value"),
			errors.New("test error"),
		},
	}

	for i, tc := range testcases {
		req, err := tc.opts.newRequest(tc.method, tc.url)
		assert.Equal(t, tc.expectErr, err, i)
		if err != nil {
			return
		}

		assert.Equal(t, req.Method, tc.method, i)
		assert.EqualValues(t, req.Header, tc.expectHeaders, i)
		assert.EqualValues(t, req.URL.Query(), tc.expectParams, i)
	}
}
