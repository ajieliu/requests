package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
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
		if err != nil {
			t.Error(err)
		}

		data, err := io.ReadAll(r)
		if err != nil {
			t.Errorf("[%d] unexpect error with read data: %v", i, err)
		}

		if string(data) != string(tc.data) {
			t.Errorf("[%d] response body: %s, expect: %s", i, data, tc.data)
		}
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
		if err != nil {
			if tc.err == nil {
				t.Errorf("[%d] unexpected error %v", i, err)
			}

			if tc.err.Error() != err.Error() {
				t.Errorf("[%d] unexpeceted error %v != %v", i, tc.err, err)
			}

			return
		}

		if tc.err != nil {
			t.Errorf("[%d] expected error %v does not occurred", i, tc.err)
		}

		data, err := io.ReadAll(r)
		if err != nil {
			t.Errorf("[%d] unexpect error: %v", i, err)
		}
		if string(tc.expect) != string(data) {
			t.Errorf("[%d] unexpect value. %s != %s", i, tc.expect, data)
		}
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

	for _, tc := range testcases {
		opts := defaultRequestOptions
		WithBodyReader(tc.r)(&opts)

		r, err := opts.bodyfn()
		if err != nil {
			t.Fatal(err)
		}

		data, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}

		if string(tc.expect) != string(data) {
			t.Errorf("with body reader data: wanted %s, got %s", tc.expect, data)
		}
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

	for _, tc := range testcases {
		opts := defaultRequestOptions
		WithHeaders(tc.h)(&opts)
		if !reflect.DeepEqual(tc.h, opts.headers) {
			t.Errorf("with headers: wanted %v, got %v", tc.h, opts.headers)
		}
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

	for _, tc := range testcases {
		opts := defaultRequestOptions
		WithParams(tc.p)(&opts)

		if tc.p.String() != opts.params.String() {
			t.Errorf("params: wanted %s, got %s", tc.p, opts.params)
		}
	}
}

func TestNewRequest(t *testing.T) {
	ioErr := errors.New("io error")
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
				nil,
				nil,
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
				nil,
				nil,
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
					return nil, ioErr
				},
				nil,
				nil,
			},
			P{}.Add("name", "test"),
			[]byte{1},
			H{}.Add("key", "value"),
			ioErr,
		},
	}

	for _, tc := range testcases {
		req, err := tc.opts.newRequest(tc.method, tc.url)
		if tc.expectErr != err {
			t.Errorf("new request error: wanted %v, got %v", tc.expectErr, err)
		}
		if err != nil {
			return
		}

		if req.Method != tc.method {
			t.Errorf("request method: wanted %s, got %s", tc.method, req.Method)
		}

		if !reflect.DeepEqual(req.Header, http.Header(tc.expectHeaders)) {
			t.Errorf("request headers: wanted %v, got %v", tc.expectHeaders, req.Header)
		}

		if req.URL.Query().Encode() != tc.expectParams.String() {
			t.Errorf("request query: wanted %s, got %s", tc.expectParams.String(), req.URL.Query().Encode())
		}
	}
}
