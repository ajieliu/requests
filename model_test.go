package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestH_Add(t *testing.T) {
	testcases := []struct {
		key       string
		values    []string
		expectKey string
	}{
		{"key", []string{"value", "value2"}, "Key"},
		{"Key", []string{"value"}, "Key"},
		{"Key", []string{}, "Key"},
	}

	for i, tc := range testcases {
		h := H{}
		for _, v := range tc.values {
			h.Add(tc.key, v)
		}

		assert.Equal(t, len(tc.values), len(h[tc.expectKey]))
		if len(tc.values) > 0 {
			assert.EqualValues(t, tc.values, h[tc.expectKey], i)
		}
	}
}

func TestH_Set(t *testing.T) {
	testcases := []struct {
		h         H
		key       string
		value     string
		expectKey string
	}{
		{H{"Key": []string{"1", "2"}}, "key", "value", "Key"},
		{H{"Key": []string{"1", "2"}}, "key", "", "Key"},
		{H{}, "1", "v", "1"},
	}

	for i, tc := range testcases {
		h := tc.h
		h.Set(tc.key, tc.value)
		vs := h[tc.expectKey]

		assert.Equal(t, 1, len(vs), i)
		assert.Equal(t, tc.value, vs[0], i)
	}
}

func TestH_Del(t *testing.T) {
	testcases := []struct {
		h         H
		key       string
		expectKey string
	}{
		{H{"Key": []string{"value"}}, "Key", "Key"},
		{H{"Key": []string{"value", "v2"}}, "key", "Key"},
		{H{}, "key", "Key"},
	}

	for i, tc := range testcases {
		h := tc.h
		h.Del(tc.key)

		_, ok := h[tc.expectKey]
		assert.False(t, ok, i)
	}
}

func TestP(t *testing.T) {
	testcases := []struct {
		p          P
		parameters [][2]string
		deletes    []string
		expect     string
	}{
		{P{}, [][2]string{{"a", "b"}, {"a", "c"}, {"b", "c"}}, []string{}, "a=b&a=c&b=c"},
		{P{}, [][2]string{{"a", "b"}, {"b", "c"}, {"d", "d"}, {"a", "c"}}, []string{"d"}, "a=b&a=c&b=c"},
		{P{"x": []string{"y", "z"}, "m": []string{"m"}}, [][2]string{{"d", "d"}, {"a", "c"}}, []string{"m"}, "a=c&d=d&x=y&x=z"},
	}

	for i, tc := range testcases {
		p := tc.p
		for _, param := range tc.parameters {
			p.Add(param[0], param[1])
		}

		for _, k := range tc.deletes {
			p.Del(k)
		}

		assert.Equal(t, tc.expect, p.String(), i)
	}
}

func TestP_Set(t *testing.T) {
	testcases := []struct {
		p     P
		key   string
		value string
	}{
		{P{}, "key", "value"},
		{P{"key": []string{"v1", "v2"}}, "key", "value"},
	}

	for i, tc := range testcases {
		p := tc.p
		p.Set(tc.key, tc.value)
		vs := p[tc.key]
		assert.Equal(t, 1, len(vs), i)
		assert.Equal(t, tc.value, vs[0], i)
	}
}
