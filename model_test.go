package requests

import (
	"reflect"
	"testing"
)

func TestH_Add(t *testing.T) {
	testcases := []struct {
		key       string
		values    []string
		expectKey string
	}{
		{"key", []string{"value", "value2"}, "Key"},
		{"Key", []string{"value"}, "Key"},
		{"Key", nil, "Key"},
	}

	for _, tc := range testcases {
		h := H{}
		for _, v := range tc.values {
			h.Add(tc.key, v)
		}

		if !reflect.DeepEqual(tc.values, h[tc.expectKey]) {
			t.Errorf("header values: wanted %v, got %v", tc.values, h[tc.expectKey])
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

		if len(vs) != 1 {
			t.Fatalf("[%d] unexpected length of values %v", i, vs)
		}
		if tc.value != vs[0] {
			t.Errorf("[%d] unmatch value %s != %s", i, tc.value, vs[0])
		}
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
		if ok {
			t.Errorf("[%d] unexpect key %s", i, tc.expectKey)
		}
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

		if tc.expect != p.String() {
			t.Errorf("[%d] unexpect string %s != %s", i, tc.expect, p.String())
		}
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

		if len(vs) != 1 {
			t.Fatalf("[%d] unexpect length of %v", i, vs)
		}

		if tc.value != vs[0] {
			t.Errorf("[%d] unexpect value. %s != %s", i, tc.value, vs[0])
		}
	}
}
