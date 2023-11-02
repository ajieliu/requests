package requests

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

type mockReadCloser struct {
	reader  io.Reader
	IsClose bool
}

func newMockReadCloser(r io.Reader) *mockReadCloser {
	return &mockReadCloser{
		reader: r,
	}
}

func (rc *mockReadCloser) Close() error {
	rc.IsClose = true
	return nil
}

func (rc *mockReadCloser) Read(p []byte) (n int, err error) {
	return rc.reader.Read(p)
}

func TestResponse_Json(t *testing.T) {
	ioErr := errors.New("mock io error")
	type mockResponseData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	testcases := []struct {
		r           io.Reader
		expectValue *mockResponseData
		expectErr   error
	}{
		{strings.NewReader("{\"id\":\"123\",\"name\":\"abc\",\"age\":21}"), &mockResponseData{"123", "abc", 21}, nil},
		{strings.NewReader("{\"id\":\"123\",\"name\":\"abc\",\"others\":\"21\"}"), &mockResponseData{"123", "abc", 0}, nil},
		{strings.NewReader("{\"id\":\"123\",\"name\":\"abc\"}"), &mockResponseData{"123", "abc", 0}, nil},
		{strings.NewReader("{}"), &mockResponseData{"", "", 0}, nil},
		{iotest.ErrReader(ioErr), &mockResponseData{}, ioErr},
	}

	for i, tc := range testcases {
		rc := newMockReadCloser(tc.r)

		resp := newResponse(&http.Response{Body: rc})
		v := &mockResponseData{}
		err := resp.Json(v)

		if tc.expectErr != err {
			t.Fatalf("response json error: wanted %v, got %v", tc.expectErr, err)
		}

		if !reflect.DeepEqual(tc.expectValue, v) {
			t.Errorf("response json data: wanted %v, got %v", tc.expectValue, v)
		}

		if rc.IsClose {
			t.Errorf("[%d] unexpected close status", i)
		}
	}
}
